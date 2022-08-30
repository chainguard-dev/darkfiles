package main

import (
	"os"

	"chainguard.dev/fagin/pkg/unpack"
)

func main() {
	u := unpack.Unpack{}
	u.PrintStats(os.Args[1])
}
