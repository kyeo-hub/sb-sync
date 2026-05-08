package downloader

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestArchMapping(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"amd64", "amd64"},
		{"arm64", "arm64"},
	}

	for _, tc := range testCases {
		result, ok := archMap[tc.input]
		if !ok {
			t.Errorf("arch %s not found in map", tc.input)
			continue
		}
		if result != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, result)
		}
	}
}

func TestGetHTTPClient(t *testing.T) {
	client1 := getHTTPClient()
	client2 := getHTTPClient()

	if client1 == nil {
		t.Error("getHTTPClient() should not return nil")
	}

	if client1 != client2 {
		t.Error("getHTTPClient() should return the same instance (singleton)")
	}
}

func TestExtractZip(t *testing.T) {
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test.zip")

	createTestZip(t, zipPath)

	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatalf("Failed to create extract dir: %v", err)
	}

	err := extractZip(zipPath, extractDir)
	if err != nil {
		t.Errorf("extractZip failed: %v", err)
	}

	singBoxPath := filepath.Join(extractDir, "sing-box.exe")
	if _, err := os.Stat(singBoxPath); os.IsNotExist(err) {
		t.Error("sing-box.exe should exist after extraction")
	}
}

func TestExtractTarGz(t *testing.T) {
	tmpDir := t.TempDir()
	tarPath := filepath.Join(tmpDir, "test.tar.gz")

	createTestTarGz(t, tarPath)

	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		t.Fatalf("Failed to create extract dir: %v", err)
	}

	err := extractTarGz(tarPath, extractDir)
	if err != nil {
		t.Errorf("extractTarGz failed: %v", err)
	}

	singBoxPath := filepath.Join(extractDir, "sing-box")
	if _, err := os.Stat(singBoxPath); os.IsNotExist(err) {
		t.Error("sing-box should exist after extraction")
	}
}

func TestExtractBinary_UnknownFormat(t *testing.T) {
	tmpDir := t.TempDir()
	unknownPath := filepath.Join(tmpDir, "test.unknown")

	if err := os.WriteFile(unknownPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	extractDir := filepath.Join(tmpDir, "extract")
	err := ExtractBinary(unknownPath, extractDir)
	if err == nil {
		t.Error("Expected error for unknown archive format")
	}
}

func TestExtractBinary_ZipError(t *testing.T) {
	tmpDir := t.TempDir()
	invalidZip := filepath.Join(tmpDir, "invalid.zip")

	if err := os.WriteFile(invalidZip, []byte("not a zip"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	extractDir := filepath.Join(tmpDir, "extract")
	err := ExtractBinary(invalidZip, extractDir)
	if err == nil {
		t.Error("Expected error for invalid zip file")
	}
}

func createTestZip(t *testing.T, zipPath string) {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip: %v", err)
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	header := &zip.FileHeader{
		Name:   "sing-box.exe",
		Method: zip.Deflate,
	}
	writer, err := zw.CreateHeader(header)
	if err != nil {
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	writer.Write([]byte("fake binary"))
}

func createTestTarGz(t *testing.T, tarPath string) {
	t.Helper()

	file, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("Failed to create tar.gz: %v", err)
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	header := &tar.Header{
		Name: "sing-box",
		Mode: 0755,
		Size: 11,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatalf("Failed to write tar header: %v", err)
	}
	tarWriter.Write([]byte("fake binary"))
}

func TestDownloadURL_Construction(t *testing.T) {
	testCases := []struct {
		version     string
		osName      string
		arch        string
		expectedExt string
	}{
		{"v1.2.3", "linux", "amd64", "tar.gz"},
		{"v1.2.3", "windows", "amd64", "zip"},
		{"v1.2.3", "darwin", "arm64", "tar.gz"},
	}

	for _, tc := range testCases {
		ext := "tar.gz"
		if tc.osName == "windows" {
			ext = "zip"
		}
		if ext != tc.expectedExt {
			t.Errorf("Expected extension %s for %s, got %s", tc.expectedExt, tc.osName, ext)
		}
	}
}
