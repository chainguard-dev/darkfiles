/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package distro

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	dpkgLib = "var/lib/dpkg"
	dbDir   = dpkgLib + "/info/" // + util-linux.list

	apkdbPath = "/lib/apk/db/installed"
)

type fsEventHandler func(*tar.Header, *tar.Reader) error

func ScanImageArchive(archivePath, format string) (filesInPackages, filesInImage []string, err error) {
	/*switch format {
	case "debian":
		return ScanDebian(archivePath)
	case "alpine":
		return ScanAlpine(archivePath)
	default:
		return nil, nil, fmt.Errorf("unkown distribution (use debian or alpine)")
	}*/
	filesInPackages, filesInImage, err = ScanDebian(archivePath)
	if err != nil || len(filesInPackages) != 0 {
		return filesInPackages, filesInImage, err
	}

	filesInPackages, filesInImage, err = ScanAlpine(archivePath)
	if err != nil || len(filesInPackages) != 0 {
		return filesInPackages, filesInImage, err
	}

	return nil, nil, fmt.Errorf("Cannot find valid package manager (debian or alpine)")
}

func ScanAlpine(archivePath string) (filesInPackages, filesInImage []string, err error) {
	filesInPackages = []string{}

	apkdbCopyFile, err := os.CreateTemp("", "apkdb-")
	if err != nil {
		return nil, nil, fmt.Errorf("opening temporary file for apkdb: %w", err)
	}

	defer func() {
		apkdbCopyFile.Close()
		os.Remove(apkdbCopyFile.Name())
	}()

	fsHandler := func(header *tar.Header, tr *tar.Reader) error {
		if "/"+header.Name == apkdbPath {
			if _, err := io.Copy(apkdbCopyFile, tr); err != nil {
				return fmt.Errorf("writing to files list: %w", err)
			}
		}
		return nil
	}

	others := map[string]struct{}{}
	otherHandler := func(header *tar.Header, tr *tar.Reader) error {
		others["/"+header.Name] = struct{}{}
		return nil
	}

	filesInImage, err = readTarFiles(archivePath, fsHandler, otherHandler)
	if err != nil {
		return nil, nil, fmt.Errorf("reading tar files: %w", err)
	}

	// Parse the apk db
	if _, err := apkdbCopyFile.Seek(0, 0); err != nil {
		return nil, nil, fmt.Errorf("rewinding apkdb: %w", err)
	}

	scanner := bufio.NewScanner(apkdbCopyFile)
	dupecheck := map[string]struct{}{}
	inDir := ""
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "F:") {
			inDir = "/" + strings.TrimPrefix(scanner.Text(), "F:")
			continue
		}

		if inDir == "" || !strings.HasPrefix(scanner.Text(), "R:") {
			continue
		}

		line := filepath.Join(inDir, strings.TrimPrefix(scanner.Text(), "R:"))
		if _, ok := others[line]; ok {
			continue
		}
		if _, ok := dupecheck[line]; ok {
			continue
		}
		filesInPackages = append(filesInPackages, line)
		dupecheck[line] = struct{}{}
	}

	return filesInPackages, filesInImage, err
}

// readTarFiles returns a list of all file entries in the tarball
// it support receiving two functions that act as handlers for
// nodes in the file system. When defined, the functions will be
// called with the tar header and the tar stream reader.
func readTarFiles(archivePath string, fileHandler, nodeHandler fsEventHandler) (files []string, err error) {
	files = []string{}
	tarFile, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("opening tar file: %w", err)
	}

	tr := tar.NewReader(tarFile)
LOOP:
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			break LOOP

		// return any other error
		case err != nil:
			return nil, fmt.Errorf("scanning file: %w", err)

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		switch header.Typeflag {
		case tar.TypeReg:
			// If its a file, add it to the all files list
			files = append(files, "/"+header.Name)
			if fileHandler != nil {
				if err := fileHandler(header, tr); err != nil {
					return nil, fmt.Errorf("calling file handler: %w", err)
				}
			}
		default:
			if err := nodeHandler(header, tr); err != nil {
				return nil, fmt.Errorf("calling other node handler")
			}
		}
	}
	return files, nil
}

// ScanDebian reads all files in a debian based image, appending the
// package index data to a file
func ScanDebian(archivePath string) (filesInPackages, filesInImage []string, err error) {
	f, err := os.CreateTemp("", "file-list-*.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("opening temp file: %w", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	// Define the handlers to be called when the file and other
	// nodes are found
	fsHandler := func(header *tar.Header, tr *tar.Reader) error {
		if strings.HasPrefix(header.Name, dbDir) && strings.HasSuffix(header.Name, ".list") {
			if _, err := io.Copy(f, tr); err != nil {
				return fmt.Errorf("writing to files list: %w", err)
			}
		}
		return nil
	}

	others := map[string]struct{}{}
	otherHandler := func(header *tar.Header, tr *tar.Reader) error {
		others["/"+header.Name] = struct{}{}
		return nil
	}

	// Read all files in the tarball
	filesInImage, err = readTarFiles(archivePath, fsHandler, otherHandler)
	if err != nil {
		return nil, nil, fmt.Errorf("reading archive files: %w", err)
	}

	// Build the pavckaged files list
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	dupecheck := map[string]struct{}{}
	for scanner.Scan() {
		if _, ok := others[scanner.Text()]; ok {
			continue
		}
		if _, ok := dupecheck[scanner.Text()]; ok {
			continue
		}
		filesInPackages = append(filesInPackages, scanner.Text())
		dupecheck[scanner.Text()] = struct{}{}
	}
	return filesInPackages, filesInImage, nil
}
