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

	"sb-sync/pkg/config"
)

const (
	archAmd64 = "amd64"
	archArm64 = "arm64"

	extTarGz = "tar.gz"
	extZip   = "zip"
)

var archMap = map[string]string{
	"amd64": archAmd64,
	"arm64": archArm64,
}

func GetLatestVersion() (string, error) {
	client := getHTTPClient()

	baseURL := config.AppConfig.GithubProxy
	if baseURL == "" {
		baseURL = ""
	}
	url := baseURL + "https://api.github.com/repos/" + config.SingBoxRepoOwner + "/" + config.SingBoxRepoName + "/releases/latest"

	var result struct {
		TagName string `json:"tag_name"`
	}

	resp, err := client.R().
		SetResult(&result).
		Get(url)

	if err != nil {
		return "", fmt.Errorf("failed to fetch latest version: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to get latest version: %s: %w", resp.Status(), err)
	}

	return result.TagName, nil
}

func DownloadSingBox(version string) (string, error) {
	client := getHTTPClient()

	osName := runtime.GOOS
	arch := runtime.GOARCH

	mappedArch, ok := archMap[arch]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	extension := extTarGz
	if osName == "windows" {
		extension = extZip
	}

	fileName := fmt.Sprintf("sing-box-%s-%s-%s.%s", strings.TrimPrefix(version, "v"), osName, mappedArch, extension)

	baseURL := config.AppConfig.GithubProxy
	var downloadURL string
	if baseURL == "" {
		downloadURL = fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
			config.SingBoxRepoOwner, config.SingBoxRepoName, version, fileName)
	} else {
		downloadURL = fmt.Sprintf("%shttps://github.com/%s/%s/releases/download/%s/%s",
			baseURL, config.SingBoxRepoOwner, config.SingBoxRepoName, version, fileName)
	}

	tmpFile := filepath.Join(os.TempDir(), fileName)

	_, err := client.R().
		SetOutput(tmpFile).
		Get(downloadURL)

	if err != nil {
		return "", fmt.Errorf("failed to download sing-box: %w", err)
	}

	return tmpFile, nil
}

func ExtractBinary(archivePath, destDir string) error {
	defer os.Remove(archivePath)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	} else if strings.HasSuffix(archivePath, ".tar.gz") {
		return extractTarGz(archivePath, destDir)
	}

	return fmt.Errorf("unknown archive format: %s", archivePath)
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "sing-box.exe") || strings.HasSuffix(f.Name, "/sing-box") || f.Name == "sing-box" {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer rc.Close()

			target := filepath.Join(dest, "sing-box.exe")
			if !strings.HasSuffix(f.Name, ".exe") {
				target = filepath.Join(dest, "sing-box")
			}
			outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("sing-box binary not found in zip")
}

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		if strings.HasSuffix(header.Name, "/sing-box") || header.Name == "sing-box" || strings.HasSuffix(header.Name, "/sing-box.exe") {
			target := filepath.Join(dest, "sing-box")
			outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			if err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("sing-box binary not found in tar.gz")
}
