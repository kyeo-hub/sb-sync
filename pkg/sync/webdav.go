package sync

import (
	"io"
	"os"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
	"sb-sync/pkg/config"
)

func SyncFromWebDAV() error {
	c := gowebdav.NewClient(config.AppConfig.WebDAV.URL, config.AppConfig.WebDAV.Username, config.AppConfig.WebDAV.Password)

	reader, err := c.ReadStream(config.AppConfig.WebDAV.FilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	destPath := filepath.Join(home, ".sb-sync", "config.json")
	
	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	return err
}
