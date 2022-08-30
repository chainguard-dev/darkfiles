package distro

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	dpkgLib = "var/lib/dpkg"
	dbDir   = dpkgLib + "/info/" // + util-linux.list
)

func ScanImageArchive(archivePath string) (filesInPackages, filesInImage []string, err error) {
	filesInPackages, filesInImage, err = ScanDebian(archivePath)
	if err != nil {
		return nil, nil, fmt.Errorf("scanning debian archive: %w", err)
	}
	return filesInPackages, filesInImage, nil
}

/*
	func ScanAlpine(archivePath string) (filesInPackages, filesInImage []string, err error) {
		filesInImage = []string{}
		filesInPackages = []string{}
		tarFile, err := os.Open(archivePath)
		if err != nil {
			return nil, nil, fmt.Errorf("opening tar file: %w", err)
		}

		tr := tar.NewReader(tarFile)

}
*/

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

type fsEventHandler func(*tar.Header, *tar.Reader) error

// ScanDebian reads all files in a debian based image, appending the
// package index data to a file
func ScanDebian(archivePath string) (filesInPackages, filesInImage []string, err error) {
	f, err := os.CreateTemp("", "file-list-*.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("opening temp file: %w", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	others := map[string]struct{}{}

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
