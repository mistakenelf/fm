package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type SettingsConfig struct {
	StartDir  string `mapstructure:"start_dir"`
	ShowIcons bool   `mapstructure:"show_icons"`
}

type ColorsConfig struct {
	SelectedItem string `mapstructure:"selected_dir_item"`
	ActivePane   string `mapstructure:"active_pane"`
	InactivePane string `mapstructure:"inactive_pane"`
	Spinner      string `mapstructure:"spinner"`
}

type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
	Colors   ColorsConfig   `mapstructure:"colors"`
}

func LoadConfig() {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "fm")
	configFile := filepath.Join(home, ".config", "fm", "config.yml")

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configFile), 0770); err != nil {
			log.Fatal("Error creating config file")
		}

		os.Create(configFile)
		viper.WriteConfig()
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Error loading config:", err)
		}
	}
}

func GetConfig() (config Config) {
	err := viper.Unmarshal(&config)

	if err != nil {
		fmt.Println("Error parsing config", err)
	}

	return
}

func SetDefaults() {
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("colors.selected_dir_item", "#F25D94")
	viper.SetDefault("colors.active_pane", "#F25D94")
	viper.SetDefault("colors.inactive_pane", "#FFFFFF")
	viper.SetDefault("colors.spinner", "#F25D94")
}
