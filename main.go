package main

import (
	"chainguard.dev/imgfiles/internal/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
