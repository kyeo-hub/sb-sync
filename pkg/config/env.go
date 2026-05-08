package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	envPrefix = "SB_SYNC_"
)

func LoadEnvOverrides() {
	if val := os.Getenv(envPrefix + "GITHUB_PROXY"); val != "" {
		AppConfig.GithubProxy = val
	}

	if val := os.Getenv(envPrefix + "WEBDAV_URL"); val != "" {
		AppConfig.WebDAV.URL = val
	}

	if val := os.Getenv(envPrefix + "WEBDAV_USER"); val != "" {
		AppConfig.WebDAV.Username = val
	}

	if val := os.Getenv(envPrefix + "WEBDAV_PASS"); val != "" {
		AppConfig.WebDAV.Password = val
	}

	if val := os.Getenv(envPrefix + "WEBDAV_PATH"); val != "" {
		AppConfig.WebDAV.FilePath = val
	}

	if val := os.Getenv(envPrefix + "INSTALL_DIR"); val != "" {
		AppConfig.InstallDir = val
	}

	if val := os.Getenv(envPrefix + "SYNC_INTERVAL"); val != "" {
		if interval, err := strconv.Atoi(val); err == nil && interval > 0 {
			AppConfig.SyncInterval = interval
		}
	}
}

func GetEnvList() []string {
	return []string{
		envPrefix + "GITHUB_PROXY",
		envPrefix + "WEBDAV_URL",
		envPrefix + "WEBDAV_USER",
		envPrefix + "WEBDAV_PASS",
		envPrefix + "WEBDAV_PATH",
		envPrefix + "INSTALL_DIR",
		envPrefix + "SYNC_INTERVAL",
	}
}

func ValidateConfig() error {
	var errors []string

	if AppConfig.InstallDir == "" {
		errors = append(errors, "install_dir cannot be empty")
	}

	if AppConfig.SyncInterval < 1 {
		errors = append(errors, "sync_interval must be at least 1 minute")
	}

	if AppConfig.GithubProxy != "" && !strings.HasSuffix(AppConfig.GithubProxy, "/") {
		errors = append(errors, "github_proxy must end with a trailing slash")
	}

	if errors != nil {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}
