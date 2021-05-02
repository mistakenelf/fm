package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/constants"

	"github.com/spf13/viper"
)

type SettingsConfig struct {
	StartDir   string `mapstructure:"start_dir"`
	ShowIcons  bool   `mapstructure:"show_icons"`
	ShowHidden bool   `mapstructure:"show_hidden"`
}

type DirTreeColors struct {
	SelectedItem      string `mapstructure:"selected_dir_item"`
	UnselectedDirItem string `mapstructure:"unselected_dir_item"`
}
type PaneColors struct {
	ActivePane   string `mapstructure:"active_pane"`
	InactivePane string `mapstructure:"inactive_pane"`
}

type ComponentColors struct {
	Spinner string `mapstructure:"spinner"`
}

type ColorsConfig struct {
	DirTree    DirTreeColors   `mapstructure:"dir_tree"`
	Pane       PaneColors      `mapstructure:"pane"`
	Components ComponentColors `mapstructure:"components"`
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
	viper.SetDefault("settings.show_hidden", true)

	viper.SetDefault("colors.dir_tree.selected_dir_item", constants.Pink)
	viper.SetDefault("colors.dir_tree.unselected_dir_item", constants.White)

	viper.SetDefault("colors.pane.active_pane", constants.Pink)
	viper.SetDefault("colors.pane.inactive_pane", constants.White)
	viper.SetDefault("colors.components.spinner", constants.Pink)
}
