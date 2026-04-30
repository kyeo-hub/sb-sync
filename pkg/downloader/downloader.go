package downloader

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-resty/resty/v2"
	"sb-sync/pkg/config"
)

func GetLatestVersion() (string, error) {
	client := resty.New()
	url := fmt.Sprintf("%shttps://api.github.com/repos/SagerNet/sing-box/releases/latest", config.AppConfig.GithubProxy)
	
	// If proxy is empty, use direct
	if config.AppConfig.GithubProxy == "" {
		url = "https://api.github.com/repos/SagerNet/sing-box/releases/latest"
	}

	var result struct {
		TagName string `json:"tag_name"`
	}

	resp, err := client.R().
		SetResult(&result).
		Get(url)

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to get latest version: %s", resp.Status())
	}

	return result.TagName, nil
}

func DownloadSingBox(version string) (string, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map GOARCH to sing-box naming
	switch arch {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	extension := "tar.gz"
	if osName == "windows" {
		extension = "zip"
	}

	fileName := fmt.Sprintf("sing-box-%s-%s-%s.%s", version[1:], osName, arch, extension)
	downloadURL := fmt.Sprintf("%shttps://github.com/SagerNet/sing-box/releases/download/%s/%s", 
		config.AppConfig.GithubProxy, version, fileName)

	if config.AppConfig.GithubProxy == "" {
		downloadURL = fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/%s/%s", version, fileName)
	}

	tmpFile := filepath.Join(os.TempDir(), fileName)
	
	client := resty.New()
	_, err := client.R().
		SetOutput(tmpFile).
		Get(downloadURL)

	if err != nil {
		return "", err
	}

	return tmpFile, nil
}

func ExtractBinary(archivePath, destDir string) error {
	defer os.Remove(archivePath) // Clean up temp file

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	} else if strings.HasSuffix(archivePath, ".tar.gz") {
		return extractTarGz(archivePath, destDir)
	}

	return fmt.Errorf("unknown archive format")
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "sing-box.exe") || strings.HasSuffix(f.Name, "sing-box") {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			target := filepath.Join(dest, filepath.Base(f.Name))
			outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			return err
		}
	}
	return fmt.Errorf("sing-box binary not found in zip")
}

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if strings.HasSuffix(header.Name, "/sing-box") || header.Name == "sing-box" {
			target := filepath.Join(dest, filepath.Base(header.Name))
			outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			return err
		}
	}
	return fmt.Errorf("sing-box binary not found in tar.gz")
}
