package sync

import (
	"io"
	"os"

	"github.com/studio-b12/gowebdav"
	"sb-sync/pkg/config"
)

func SyncFromWebDAV() error {
	_, _, err := SyncFromWebDAVWithStatus()
	return err
}

func SyncFromWebDAVWithStatus() (bool, bool, error) {
	c := gowebdav.NewClient(config.AppConfig.WebDAV.URL, config.AppConfig.WebDAV.Username, config.AppConfig.WebDAV.Password)

	info, err := c.Stat(config.AppConfig.WebDAV.FilePath)
	if err != nil {
		return false, false, err
	}

	destPath := config.GetConfigPath()
	localInfo, err := os.Stat(destPath)
	
	// If local file exists, check if remote is different
	if err == nil {
		if localInfo.Size() == info.Size() && localInfo.ModTime().After(info.ModTime()) {
			// Probably same or local is newer (unlikely but possible)
			// This is a simple heuristic. Better would be ETag or Hash.
			return true, false, nil
		}
	}

	reader, err := c.ReadStream(config.AppConfig.WebDAV.FilePath)
	if err != nil {
		return false, false, err
	}
	defer reader.Close()

	tmpPath := destPath + ".tmp"
	outFile, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return false, false, err
	}

	if _, err := io.Copy(outFile, reader); err != nil {
		outFile.Close()
		os.Remove(tmpPath)
		return false, false, err
	}
	outFile.Close()

	if err := os.Rename(tmpPath, destPath); err != nil {
		return false, false, err
	}

	// Update local mod time to match remote
	os.Chtimes(destPath, info.ModTime(), info.ModTime())

	return true, true, nil
}
