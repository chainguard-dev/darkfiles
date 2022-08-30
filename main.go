package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// "github.com/google/go-containerregistry/pkg/crane"
// import ( k8s.io )

const (
	dpkgLib = "var/lib/dpkg"
	dbDir   = dpkgLib + "/info/" // + util-linux.list
)

func main() {
	filesInPackages, filesInImage, err := buildFileLists(os.Args[1])
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	inversePackageDict := map[string]struct{}{}
	for _, filename := range filesInPackages {
		inversePackageDict[filename] = struct{}{}
	}

	inverseBaseDict := map[string]struct{}{}
	for _, filename := range baseLinuxFiles {
		inverseBaseDict[filename] = struct{}{}
	}

	// fmt.Printf("%+v\n", inversePackageDict)

	for _, filename := range filesInImage {
		if _, ok := inverseBaseDict[filename]; ok {
			continue
		}
		if _, ok := inversePackageDict[filename]; !ok {
			// Exceptions:

			// Files created by dpkg itsef:
			if strings.HasPrefix(filename, "/"+dpkgLib) {
				continue
			}

			// debconf database and helpers
			if strings.HasPrefix(filename, "/var/cache/debconf/") ||
				strings.HasPrefix(filename, "/var/lib/systemd/deb-systemd-helper-enabled") ||
				strings.HasPrefix(filename, "/var/lib/systemd/deb-systemd-user-helper-enabled") {
				continue
			}

			// Python cache
			if strings.Contains(filename, "/__pycache__/") {
				continue
			}

			// Python packages databases
			if strings.HasPrefix(filename, "/var/lib/python/") &&
				strings.HasSuffix(filename, "_installed") {
				continue
			}

			// For now, lets ignore all config, logs and runtime directories
			if strings.HasPrefix(filename, "/var/log/") ||
				strings.HasPrefix(filename, "/etc/") ||
				strings.HasPrefix(filename, "/var/cache") ||
				strings.HasPrefix(filename, "/run/") {
				continue
			}

			// User confg databases
			if strings.HasPrefix(filename, "/var/lib/ucf/") {
				continue
			}

			// pam configs created on the fly
			if strings.HasPrefix(filename, "/var/lib/pam/") {
				continue
			}

			//

			fmt.Println(filename)
		}
	}

}

func buildFileLists(tarFilePath string) (filesInPackages, filesInImage []string, err error) {
	filesInImage = []string{}
	filesInPackages = []string{}
	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("opening tar file: %w", err)
	}
	f, err := os.CreateTemp("", "file-list-*.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("opening temp file: %w", err)
	}
	defer f.Close()
	gzr, err := gzip.NewReader(tarFile)
	if err != nil {
		return nil, nil, fmt.Errorf("creating tar reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
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

		// check the file type
		switch header.Typeflag {

		// if it's a file create it
		case tar.TypeReg:
			// If its a file, add it to the all files list
			filesInImage = append(filesInImage, "/"+header.Name)
			if strings.HasPrefix(header.Name, dbDir) && strings.HasSuffix(header.Name, ".list") {
				// And if the file is a package file list, read it
				if _, err := io.Copy(f, tr); err != nil {
					return nil, nil, fmt.Errorf("writing to files list: %w", err)
				}
			}
		}
	}

	// Build the pavckaged files list
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// the db has a hash and the path
		/*
			_, path, ok := strings.Cut(scanner.Text(), " ")
			if ok {
				filesInPackages = append(filesInPackages, path)
			} else {
				return nil, nil, fmt.Errorf("unknown entry: %s", scanner.Text())
			}
		*/
		filesInPackages = append(filesInPackages, scanner.Text())
	}
	return filesInPackages, filesInImage, nil
}

var baseLinuxFiles = []string{
	"/etc/environment",
	"/etc/fstab",
	"/etc/group",
	"/etc/group-",
	"/etc/gshadow",
	"/etc/hostname",
	"/etc/ld.so.cache",
	"/etc/motd",
	"/etc/nsswitch.conf",
	"/etc/pam.d/common-account",
	"/etc/pam.d/common-auth",
	"/etc/pam.d/common-password",
	"/etc/pam.d/common-session",
	"/etc/pam.d/common-session-noninteractive",
	"/etc/passwd",
	"/etc/passwd-",
	"/etc/profile",
	"/etc/resolv.conf",
	"/etc/security/opasswd",
	"/etc/shadow",
	"/etc/shadow-",
	"/etc/shells",
	"/etc/subgid",
	"/etc/subuid",
	"/etc/timezone",
	"/root/.bashrc",
	"/root/.profile",
	// Standard base logs
	"/var/log/btmp",
	"/var/log/faillog",
	"/var/log/lastlog",
	"/var/log/wtmp",
	"/usr/share/info/dir",
	"/usr/sbin/policy-rc.d",
	// Deb db states
	"/var/lib/apt/extended_states",
}
