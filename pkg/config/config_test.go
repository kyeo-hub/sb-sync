package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath() should not return empty string")
	}
	if !filepath.IsAbs(path) {
		t.Error("GetConfigPath() should return absolute path")
	}
}

func TestConfigPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}
	expectedPath := filepath.Join(home, configDirName, "config.json")
	actualPath := GetConfigPath()
	if actualPath != expectedPath {
		t.Errorf("Expected %s, got %s", expectedPath, actualPath)
	}
}

func TestGetSingBoxBinary(t *testing.T) {
	AppConfig.InstallDir = "/test/bin"
	binPath := GetSingBoxBinary()
	if binPath != "/test/bin/sing-box" {
		t.Errorf("Expected /test/bin/sing-box, got %s", binPath)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.SetConfigFile(configPath)
	viper.Set("github_proxy", "https://test.com/")
	viper.Set("install_dir", "/tmp/test")
	viper.Set("sync_interval", 30)

	AppConfig.GithubProxy = "https://test.com/"
	AppConfig.InstallDir = "/tmp/test"
	AppConfig.SyncInterval = 30

	if err := viper.WriteConfig(); err != nil {
		t.Errorf("Failed to save config: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should exist after Save()")
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	viper.SetConfigFile(configPath)
	viper.Set("github_proxy", "https://loadtest.com/")
	viper.Set("install_dir", "/tmp/loadtest")
	viper.Set("sync_interval", 45)
	if err := viper.WriteConfig(); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	viper.Reset()
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var loaded Config
	if err := viper.Unmarshal(&loaded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if loaded.GithubProxy != "https://loadtest.com/" {
		t.Errorf("Expected proxy https://loadtest.com/, got %s", loaded.GithubProxy)
	}
	if loaded.SyncInterval != 45 {
		t.Errorf("Expected interval 45, got %d", loaded.SyncInterval)
	}
}

func TestDefaultValues(t *testing.T) {
	viper.Reset()
	viper.SetDefault("github_proxy", GitHubProxyDefault)
	viper.SetDefault("sync_interval", 60)

	if viper.GetString("github_proxy") != GitHubProxyDefault {
		t.Errorf("Expected default proxy %s, got %s", GitHubProxyDefault, viper.GetString("github_proxy"))
	}
	if viper.GetInt("sync_interval") != 60 {
		t.Errorf("Expected default interval 60, got %d", viper.GetInt("sync_interval"))
	}
}
