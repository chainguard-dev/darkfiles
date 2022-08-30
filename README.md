# fagin

<img src=https://i.dailymail.co.uk/i/pix/scaled/2011/12/03/article-0-034532240000044D-201_308x185.jpg>

When you need to learn how to pick a pocket or two, call Fagin.

Fagin (named after [the Oliver Twist Villain](https://en.wikipedia.org/wiki/Fagin)) is a tool finds orphaned files in container images and makes them do bad deeds, 
like [making them drink gin](https://youtu.be/-BtRMxBYaqs?t=28).

Run `fagin` against a container image reference and it will return a list of files
not tracker or installed by the OS package manager. 

## Usage 

Just run `fagin stats imageref` to get some statistics about files not installed via the os package 
manager:

```
go run ./main.go stats golang
INFO[0000] flattening image index.docker.io/library/golang 
INFO[0021] flattenned image to /tmp/image-dump-2645404725.tar (962 MB) 

Total files in image:       24906
Files in packages:          12881
Files not in packages:      12025
Tracked by package manager: 51.718460%

```

There is also `fagin list --set=all imageref` which can give you all files in 
an image (`--set=all`), files tracked by the package manager (`--set=tracked`)
and all files found in the image which were add via other means (`--set=untracked`):

```
go run ./main.go list --set=untracked golang
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

### apk support

Right now, fagin is limited to analyzing debian images (and derivatives). 