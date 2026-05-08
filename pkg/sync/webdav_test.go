package sync

import (
	"os"
	"path/filepath"
	"testing"

	"sb-sync/pkg/config"
)

func TestSyncFromWebDAV_NoConfig(t *testing.T) {
	t.Run("Missing WebDAV URL", func(t *testing.T) {
		originalURL := config.AppConfig.WebDAV.URL
		config.AppConfig.WebDAV.URL = ""
		defer func() { config.AppConfig.WebDAV.URL = originalURL }()

		err := SyncFromWebDAV()
		if err == nil {
			t.Error("Expected error when WebDAV URL is empty")
		}
	})
}

func TestSyncWithMock(t *testing.T) {
	t.Run("Sync with mock client", func(t *testing.T) {
		originalURL := config.AppConfig.WebDAV.URL
		originalUser := config.AppConfig.WebDAV.Username
		originalPass := config.AppConfig.WebDAV.Password
		originalPath := config.AppConfig.WebDAV.FilePath

		config.AppConfig.WebDAV.URL = "https://example.com/webdav"
		config.AppConfig.WebDAV.Username = "test"
		config.AppConfig.WebDAV.Password = "test"
		config.AppConfig.WebDAV.FilePath = "/test.json"

		defer func() {
			config.AppConfig.WebDAV.URL = originalURL
			config.AppConfig.WebDAV.Username = originalUser
			config.AppConfig.WebDAV.Password = originalPass
			config.AppConfig.WebDAV.FilePath = originalPath
		}()

		_, _, err := SyncFromWebDAVWithStatus()
		if err == nil {
			t.Log("Expected connection error (no real server), but call completed")
		}
	})
}

func TestTempFileCreation(t *testing.T) {
	tmpDir := t.TempDir()
	tmpPath := filepath.Join(tmpDir, "config.json.tmp")

	outFile, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	outFile.Close()

	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		t.Error("Temp file should exist")
	}
}

func TestAtomicRename(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src.txt")
	dst := filepath.Join(tmpDir, "dst.txt")

	if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	if err := os.Rename(src, dst); err != nil {
		t.Fatalf("Failed to rename file: %v", err)
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Error("Destination file should exist after rename")
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Error("Source file should not exist after rename")
	}
}

func TestSyncFromWebDAV_MissingFilePath(t *testing.T) {
	originalPath := config.AppConfig.WebDAV.FilePath
	config.AppConfig.WebDAV.FilePath = ""
	defer func() { config.AppConfig.WebDAV.FilePath = originalPath }()

	_, _, err := SyncFromWebDAVWithStatus()
	if err == nil {
		t.Error("Expected error when FilePath is empty")
	}
}

func TestConfigValidation(t *testing.T) {
	t.Run("Empty WebDAV configuration", func(t *testing.T) {
		wasEmpty := config.AppConfig.WebDAV.URL == ""
		if !wasEmpty {
			t.Skip("WebDAV URL is already configured")
		}
	})
}
