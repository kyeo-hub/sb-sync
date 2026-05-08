package sync

import (
	"encoding/json"
	"fmt"
	"os"

	"sb-sync/pkg/config"
)

func ValidateConfigJSON() error {
	configPath := config.GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist at: %s", configPath)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var jsonConfig map[string]interface{}
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return fmt.Errorf("config file is not valid JSON: %w", err)
	}

	return nil
}

func CheckWebDAVConnection() error {
	if config.AppConfig.WebDAV.URL == "" {
		return fmt.Errorf("WebDAV URL is not configured")
	}

	if config.AppConfig.WebDAV.FilePath == "" {
		return fmt.Errorf("WebDAV file path is not configured")
	}

	client := newWebDAVClient()
	_, err := client.Stat(config.AppConfig.WebDAV.FilePath)
	if err != nil {
		return fmt.Errorf("failed to connect to WebDAV server: %w", err)
	}

	return nil
}

func IsInstalled() bool {
	binPath := config.GetSingBoxBinary()
	_, err := os.Stat(binPath)
	return err == nil
}

type HealthStatus struct {
	SingBoxInstalled bool `json:"singbox_installed"`
	ConfigExists     bool `json:"config_exists"`
	WebDAVConfigured bool `json:"webdav_configured"`
	WebDAVReachable  bool `json:"webdav_reachable"`
}

func HealthCheck() HealthStatus {
	status := HealthStatus{
		SingBoxInstalled: IsInstalled(),
		ConfigExists:     false,
		WebDAVConfigured: false,
		WebDAVReachable:  false,
	}

	configPath := config.GetConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		status.ConfigExists = true
	}

	if config.AppConfig.WebDAV.URL != "" && config.AppConfig.WebDAV.FilePath != "" {
		status.WebDAVConfigured = true
		if err := CheckWebDAVConnection(); err == nil {
			status.WebDAVReachable = true
		}
	}

	return status
}
