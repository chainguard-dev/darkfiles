# fagin

<img src=https://i.dailymail.co.uk/i/pix/scaled/2011/12/03/article-0-034532240000044D-201_308x185.jpg>

When you need to learn how to pick a pocket or two, call Fagin.

Fagin (named after [the Oliver Twist Villain](https://en.wikipedia.org/wiki/Fagin)) is a tool finds orphaned files in container images and makes them do bad deeds, 
like [making them drink gin](https://youtu.be/-BtRMxBYaqs?t=28).

Run `fagin` against a container image reference and it will return a list of files
not tracker or installed by the OS package manager. 

## Usage 

Just run `fagin imageref` to get a list of files not installed via the os package 
manager:

```
go run ./main.go golang
INFO[0000] flattening image index.docker.io/library/golang 
INFO[0021] flattenned image to /tmp/image-dump-2645404725.tar (962 MB) 

Total files in image:       24906
Files in packages:          12881
Files not in packages:      12025
Tracked by package manager: 51.718460%

```

## TODO

### apk support

Right now, fagin is limited to analyzing debian images (and derivatives). 