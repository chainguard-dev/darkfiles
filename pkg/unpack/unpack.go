package unpack

import (
	"errors"
	"fmt"
	"log"
	"os"

	"chainguard.dev/darkfiles/pkg/distro"
	"chainguard.dev/darkfiles/pkg/oci"
)

func New() *Unpack {
	return &Unpack{
		Options: Options{},
	}
}

type Options struct {
	Distro string
}

type Unpack struct {
	Options Options
}

type ImageStats struct {
	FilesInImage       int     `json:"files_in_image"`
	FilesInPackages    int     `json:"files_in_package"`
	FilesNotInPackages int     `json:"files_not_in_package"`
	PercentTracked     float32 `json:"percent_tracked"`
}

// ImageStats returns an image struct
func (u *Unpack) ImageStats(imageRef string) (stats ImageStats, err error) {
	tmpFile, err := oci.FlattenImageToTmp(imageRef)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)

	var filesInPackages []string
	var filesInImage []string

	switch u.Options.Distro {
	case "debian":
		filesInPackages, filesInImage, err = distro.ScanDebian(tmpFile)
	case "alpine":
		filesInPackages, filesInImage, err = distro.ScanAlpine(tmpFile)
	default:
		return stats, errors.New("unkown distribution")
	}

	if err != nil {
		return stats, fmt.Errorf("scanning image: %w", err)
	}

	filesInImage = distro.FilterPaths(filesInImage)
	filesInPackages = distro.FilterPaths(filesInPackages)

	return ImageStats{
		FilesInImage:       len(filesInImage),
		FilesInPackages:    len(filesInPackages),
		FilesNotInPackages: len(filesInImage) - len(filesInPackages),
		PercentTracked:     (float32(len(filesInPackages)) / float32(len(filesInImage))) * 100,
	}, nil
}

func (u *Unpack) PrintStats(imageRef string) error {
	stats, err := u.ImageStats(imageRef)
	if err != nil {
		return fmt.Errorf("getting image stats: %w", err)
	}

	fmt.Printf("Total files in image:       %d\n", stats.FilesInImage)
	fmt.Printf("Files in packages:          %d\n", stats.FilesInPackages)
	fmt.Printf("Files not in packages:      %d\n", stats.FilesNotInPackages)
	fmt.Printf("Tracked by package manager: %f%%\n", stats.PercentTracked)
	return nil
}

func (u *Unpack) List(imageRef, set string) error {
	tmpFile, err := oci.FlattenImageToTmp(imageRef)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)

	filesInPackages, filesInImage, err := distro.ScanImageArchive(tmpFile, u.Options.Distro)
	if err != nil {
		return fmt.Errorf("scanning image: %w", err)
	}

	// TODO: Make filtering optional?
	filesInPackages = distro.FilterPaths(filesInPackages)
	filesInImage = distro.FilterPaths(filesInImage)

	switch set {
	case "all":
		for _, p := range filesInImage {
			fmt.Println(p)
		}
	case "tracked":
		for _, p := range filesInPackages {
			fmt.Println(p)
		}
	case "untracked":
		reverseDict := map[string]struct{}{}
		for _, p := range filesInPackages {
			reverseDict[p] = struct{}{}
		}

		for _, p := range filesInImage {
			if _, ok := reverseDict[p]; !ok {
				fmt.Println(p)
			}
		}
	default:
		return errors.New("nnknown set of files")
	}
	return nil
}
