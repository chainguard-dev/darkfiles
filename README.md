# darkfiles: Measure container image dark matter

<img src="https://user-images.githubusercontent.com/3935899/187793479-eb09b826-2e8a-486d-a9d5-fbdfdfa80691.png" />

Container images are often assembled by copying files straight
to their layers: examples include (but are not limited to) copying applications from
build images or copying entrypoint scripts. These files are
not tracked by the underlying operating system (OS) package managers like apt or apk,
which makes them invisible to security scanners since they are
not reflected in the package database.

To gain insight into this "dark matter," `darkfiles` inspects
container images to detect files not tracked by the OS package manager.

## Usage 

Just run `darkfiles stats imageref` to get some statistics about
files not installed via the OS package manager. Here is an example
scanning both the 
[official](https://github.com/docker-library/golang) and 
[distroless golang](https://github.com/distroless/go) images: 

```
darkfiles stats --distro=debian golang:latest
INFO flattening image index.docker.io/library/golang 
INFO flattened image to /tmp/image-dump-582865974.tar (962 MB) 
Total files in image:       21033
Files in packages:          8807
Files not in packages:      12226
Tracked by package manager: 41.872295%

darkfiles stats --distro=alpine distroless.dev/go
INFO flattening image distroless.dev/go           
INFO flattened image to /tmp/image-dump-7982759.tar (540 MB) 
Total files in image:       5734
Files in packages:          5734
Files not in packages:      0
Tracked by package manager: 100.000000%
```

There is also `darkfiles list --set=all imageref` which returns all files included in an image (`--set=all`), files tracked
by the package manager (`--set=tracked`) and all files found in
the image which were added through other means (`--set=untracked`):

```
darkfiles --distro=debian --set=untracked golang
INFO[0000] flattening image index.docker.io/library/golang 
INFO[0021] flattenned image to /tmp/image-dump-2645404725.tar (962 MB) 

/usr/local/go/CONTRIBUTING.md
/usr/local/go/LICENSE
/usr/local/go/PATENTS
/usr/local/go/README.md
/usr/local/go/SECURITY.md
/usr/local/go/VERSION
/usr/local/go/api/README
/usr/local/go/api/except.txt
/usr/local/go/api/go1.1.txt
/usr/local/go/api/go1.10.txt
/usr/local/go/api/go1.11.txt
/usr/local/go/api/go1.12.txt
/usr/local/go/api/go1.13.txt
/usr/local/go/api/go1.14.txt
/usr/local/go/api/go1.15.txt
/usr/local/go/api/go1.16.txt
/usr/local/go/api/go1.17.txt
/usr/local/go/api/go1.18.txt
/usr/local/go/api/go1.19.txt
/usr/local/go/api/go1.2.txt
/usr/local/go/api/go1.3.txt
/usr/local/go/api/go1.4.txt
/usr/local/go/api/go1.5.txt
... 
```

## TODO

### Automatic distro detection

Right now, specifying --distro is mandatory, it would be nice to 
have logic to automatically detect what kind of OS the image is
based on.

### Filter disable flag

There should be a flag to disable file filtering to get absolutel
all files.

### Expand stats output with:

* Number of files before and after filtering
* Image size

### Enable `--format=json` in stats subcommand

### ~apk support~
