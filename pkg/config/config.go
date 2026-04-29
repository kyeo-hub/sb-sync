package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	WebDAV struct {
		URL      string `mapstructure:"url"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		FilePath string `mapstructure:"file_path"` // Path to config.json on WebDAV
	} `mapstructure:"webdav"`
	GithubProxy string `mapstructure:"github_proxy"`
	InstallDir  string `mapstructure:"install_dir"`
}

var AppConfig Config

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".sb-sync")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.SetDefault("github_proxy", "https://ghproxy.com/")
	viper.SetDefault("install_dir", filepath.Join(configDir, "bin"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// Write default config if not exists
		if err := viper.SafeWriteConfig(); err != nil {
			// If file already exists, it's fine
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	return nil
}

func Save() error {
	viper.Set("webdav", AppConfig.WebDAV)
	viper.Set("github_proxy", AppConfig.GithubProxy)
	viper.Set("install_dir", AppConfig.InstallDir)
	return viper.WriteConfig()
}
