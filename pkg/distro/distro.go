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

func ScanDebian(archivePath string) (filesInPackages, filesInImage []string, err error) {
	filesInImage = []string{}
	filesInPackages = []string{}
	tarFile, err := os.Open(archivePath)
	if err != nil {
		return nil, nil, fmt.Errorf("opening tar file: %w", err)
	}

	tr := tar.NewReader(tarFile)

	f, err := os.CreateTemp("", "file-list-*.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("opening temp file: %w", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	others := map[string]struct{}{}

LOOP:
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			break LOOP

		// return any other error
		case err != nil:
			return nil, nil, fmt.Errorf("scanning file: %w", err)

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		switch header.Typeflag {
		case tar.TypeReg:
			// If its a file, add it to the all files list
			filesInImage = append(filesInImage, "/"+header.Name)
			if strings.HasPrefix(header.Name, dbDir) && strings.HasSuffix(header.Name, ".list") {
				// And if the file is a package file list, read it
				if _, err := io.Copy(f, tr); err != nil {
					return nil, nil, fmt.Errorf("writing to files list: %w", err)
				}
			}
		default:
			others["/"+header.Name] = struct{}{}
		}
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
