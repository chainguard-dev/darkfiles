package main

import (
	"chainguard.dev/fagin/internal/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
