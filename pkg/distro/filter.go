/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package distro

import "strings"

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

	"/dev/pts", // Shared memory mount
	"/dev/pts", // terms
}

// ignoreDirectories is a list of directories to ignore which shift
// dynamically as the OS runs or that are created at install time
var ignoreDirectories = []string{
	"/etc",          // All configuration files under etc
	"/var/log/",     // Log files
	"/var/cache",    // Cache files
	"/run/",         // Runtime files
	"/tmp",          // Temporary files
	"/var/lib/ucf/", // User config data
	"/var/lib/pam/", // pam session files and such
	//  unit files
	"/var/lib/systemd/deb-systemd-helper-enabled",
	"/var/lib/systemd/deb-systemd-user-helper-enabled",
	"/" + dpkgLib, // apt  system data

	"/usr/share/man/",    // man pages
	"/usr/share/doc/",    // documentation
	"/usr/share/info/",   // info pages
	"/usr/share/locale/", // localization files

	// Alpine exceptions
	"/lib/apk/db/", // apkdb files mutated by the package manager
}

func FilterPaths(fileList []string) (newFileList []string) {
	inverseBaseDict := map[string]struct{}{}
	for _, filename := range baseLinuxFiles {
		inverseBaseDict[filename] = struct{}{}
	}

	newFileList = []string{}

LOOP:
	for _, filename := range fileList {
		if _, ok := inverseBaseDict[filename]; ok {
			continue
		}

		// Exceptions:

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
		for _, prefix := range ignoreDirectories {
			if strings.HasPrefix(filename, prefix) {
				continue LOOP
			}
		}

		newFileList = append(newFileList, filename)
	}
	return newFileList
}
