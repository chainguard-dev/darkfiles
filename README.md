# fagin

When you need to learn how to pick a pocket or two, call Fagin.

Fagin (named after [the Oliver Twist Villain](https://en.wikipedia.org/wiki/Fagin)) is a tool finds orphaned files in container images and makes them do bad deeds, 
like [making them drink gin](https://youtu.be/-BtRMxBYaqs?t=28).

Run `fagin` against a container image reference and it will return a list of files
not tracker or installed by the OS package manager. 

## Usage 

Just run `fagin imageref` to get a list of files not installed via the os package 
manager:

```
fagin golang

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
...
(and 12'000+ files)
```

## TODO

### Finish Image Flattening Code

For now, fagin can only run on flattened image archives only, to get some samples:

```bash
wget https://storage.googleapis.com/puerco-chainguard-public/flat-image-samples/debian.tar.gz
wget https://storage.googleapis.com/puerco-chainguard-public/flat-image-samples/golang.tar.gz
wget https://storage.googleapis.com/puerco-chainguard-public/flat-image-samples/wordpress.tar.gz
```

Then run fagin with the archive as first arg

```
fagin wordpress.tar.gz
```




