package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/constants"

	"github.com/spf13/viper"
)

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
	viper.SetDefault("settings.rounded_panes", false)

	// DirTree colors
	viper.SetDefault("colors.dir_tree.selected_dir_item", constants.Pink)
	viper.SetDefault("colors.dir_tree.unselected_dir_item", constants.White)

	// Pane colors
	viper.SetDefault("colors.pane.active_pane", constants.Pink)
	viper.SetDefault("colors.pane.inactive_pane", constants.White)

	// Component colors
	viper.SetDefault("colors.components.spinner", constants.Pink)

	// StatusBar colors
	viper.SetDefault("colors.status_bar.selected_file.foreground", constants.White)
	viper.SetDefault("colors.status_bar.selected_file.background", constants.Pink)
	viper.SetDefault("colors.status_bar.bar.foreground", constants.White)
	viper.SetDefault("colors.status_bar.bar.background", constants.DarkGray)
	viper.SetDefault("colors.status_bar.total_files.foreground", constants.White)
	viper.SetDefault("colors.status_bar.total_files.background", constants.LightPurple)
	viper.SetDefault("colors.status_bar.logo.foreground", constants.White)
	viper.SetDefault("colors.status_bar.logo.background", constants.DarkPurple)
}
