package oci

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sirupsen/logrus"
)

// FlattenImageToTmp writes an image to a temporary file and returns the path
func FlattenImageToTmp(reference string) (string, error) {
	tmp, err := os.CreateTemp("", "image-dump-*.tar")
	if err != nil {
		return "", fmt.Errorf("creating temporary file: %w", err)
	}
	defer tmp.Close()

	if err := FlattenImageToWriter(reference, tmp); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("flattening image: %w", err)
	}
	info, err := os.Stat(tmp.Name())
	if err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("inspecting tmp file")
	}
	logrus.Infof("flattenned image to %s (%d MB)", tmp.Name(), info.Size()/1024/1024)
	return tmp.Name(), nil
}

// FlattenImageToWriter flattens an image and writes it to the passed writer
func FlattenImageToWriter(referenceString string, w io.Writer) error {
	ref, err := name.ParseReference(referenceString)
	if err != nil {
		log.Fatalf("parsing %s: %v", referenceString, err)
	}

	logrus.Infof("flattening image %s", ref.Context().String())

	// desc, err := remote.Get(ref, o.Remote...)
	desc, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return fmt.Errorf("pulling %s: %w", ref, err)
	}
	img, err := desc.Image()
	if err != nil {
		return fmt.Errorf("getting imasge from descriptor: %w", err)
	}

	fs := mutate.Extract(img)
	if _, err = io.Copy(w, fs); err != nil {
		return fmt.Errorf("writing flattened image: %w", err)
	}
	return nil
}
