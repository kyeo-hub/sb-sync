package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/studio-b12/gowebdav"
	"sb-sync/pkg/config"
)

func newWebDAVClient() *gowebdav.Client {
	return gowebdav.NewClient(config.AppConfig.WebDAV.URL, config.AppConfig.WebDAV.Username, config.AppConfig.WebDAV.Password)
}

func SyncFromWebDAV() error {
	_, _, err := SyncFromWebDAVWithStatus()
	return err
}

func SyncFromWebDAVWithStatus() (bool, bool, error) {
	c := newWebDAVClient()

	info, err := c.Stat(config.AppConfig.WebDAV.FilePath)
	if err != nil {
		return false, false, fmt.Errorf("failed to stat remote file: %w", err)
	}

	destPath := config.GetConfigPath()
	localInfo, err := os.Stat(destPath)

	if err == nil {
		if localInfo.Size() == info.Size() && localInfo.ModTime().After(info.ModTime()) {
			return true, false, nil
		}
	}

	reader, err := c.ReadStream(config.AppConfig.WebDAV.FilePath)
	if err != nil {
		return false, false, fmt.Errorf("failed to read remote file: %w", err)
	}
	defer reader.Close()

	tmpPath := destPath + ".tmp"
	outFile, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return false, false, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := io.Copy(outFile, reader); err != nil {
		outFile.Close()
		os.Remove(tmpPath)
		return false, false, fmt.Errorf("failed to copy file: %w", err)
	}
	outFile.Close()

	if err := os.Rename(tmpPath, destPath); err != nil {
		return false, false, fmt.Errorf("failed to rename file: %w", err)
	}

	os.Chtimes(destPath, info.ModTime(), info.ModTime())

	return true, true, nil
}
