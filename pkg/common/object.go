package common

import (
	"context"
	"os"

	"github.com/mholt/archiver/v3"
)

func ExtractObjectFile(ctx context.Context, objectPath, destPath string) error {
	if _, err := os.Stat(destPath); !os.IsNotExist(err) {
		// Folder already exists, so skip extraction
		return nil
	}

	os.MkdirAll(destPath, 0755)

	// Check if the object file exists
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		return err
	}

	zip := archiver.NewZip()
	if err := zip.Unarchive(objectPath, destPath); err != nil {
		return err
	}

	return nil
}
