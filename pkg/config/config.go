package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configDirName  = ".sb-sync"
	configFileName = "config"
	configFileType = "yaml"
)

type Config struct {
	WebDAV struct {
		URL      string `mapstructure:"url"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		FilePath string `mapstructure:"file_path"`
	} `mapstructure:"webdav"`
	GithubProxy  string `mapstructure:"github_proxy"`
	InstallDir   string `mapstructure:"install_dir"`
	SyncInterval int    `mapstructure:"sync_interval"`
}

var AppConfig Config

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
		if home == "" {
			home = "/tmp"
		}
	}
	return filepath.Join(home, configDirName, "config.json")
}

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(home, configDirName)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configDir)

	viper.SetDefault("github_proxy", GitHubProxyDefault)
	viper.SetDefault("install_dir", filepath.Join(configDir, "bin"))
	viper.SetDefault("sync_interval", 60)
	viper.SetDefault("webdav.url", "")
	viper.SetDefault("webdav.username", "")
	viper.SetDefault("webdav.password", "")
	viper.SetDefault("webdav.file_path", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		if err := viper.SafeWriteConfig(); err != nil {
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if !filepath.IsAbs(AppConfig.InstallDir) {
		AppConfig.InstallDir = filepath.Join(configDir, AppConfig.InstallDir)
	}

	return nil
}

func Save() error {
	viper.Set("webdav.url", AppConfig.WebDAV.URL)
	viper.Set("webdav.username", AppConfig.WebDAV.Username)
	viper.Set("webdav.password", AppConfig.WebDAV.Password)
	viper.Set("webdav.file_path", AppConfig.WebDAV.FilePath)
	viper.Set("github_proxy", AppConfig.GithubProxy)
	viper.Set("install_dir", AppConfig.InstallDir)
	viper.Set("sync_interval", AppConfig.SyncInterval)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}

func GetInstallDir() string {
	if filepath.IsAbs(AppConfig.InstallDir) {
		return AppConfig.InstallDir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, configDirName, AppConfig.InstallDir)
}

func GetSingBoxBinary() string {
	binPath := filepath.Join(GetInstallDir(), "sing-box")
	if filepath.Separator == '\\' {
		binPath += ".exe"
	}
	return binPath
}
