package unpack

import (
	"fmt"
	"log"
	"os"

	"chainguard.dev/fagin/pkg/distro"
	"chainguard.dev/fagin/pkg/oci"
)

func New() *Unpack {
	return &Unpack{}
}

type Unpack struct {
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

	filesInPackages, filesInImage, err := distro.ScanImageArchive(tmpFile)
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
