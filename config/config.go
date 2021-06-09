package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/utils"

	"github.com/spf13/viper"
)

type SettingsConfig struct {
	StartDir         string `mapstructure:"start_dir"`
	ShowIcons        bool   `mapstructure:"show_icons"`
	RoundedPanes     bool   `mapstructure:"rounded_panes"`
	EnableLogging    bool   `mapstructure:"enable_logging"`
	EnableMouseWheel bool   `mapstructure:"enable_mousewheel"`
}

type DirTreeColors struct {
	SelectedItem   string `mapstructure:"selected_item"`
	UnselectedItem string `mapstructure:"unselected_item"`
}
type PaneColors struct {
	ActiveBorderColor   string `mapstructure:"active_border_color"`
	InactiveBorderColor string `mapstructure:"inactive_border_color"`
}

type ComponentColors struct {
	Spinner string `mapstructure:"spinner"`
}
type ColorVariant struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

type StatusBarColors struct {
	SelectedFile ColorVariant `mapstructure:"selected_file"`
	Bar          ColorVariant `mapstructure:"bar"`
	TotalFiles   ColorVariant `mapstructure:"total_files"`
	Logo         ColorVariant `mapstructure:"logo"`
}

type ColorsConfig struct {
	DirTree    DirTreeColors   `mapstructure:"dir_tree"`
	Pane       PaneColors      `mapstructure:"pane"`
	Components ComponentColors `mapstructure:"components"`
	StatusBar  StatusBarColors `mapstructure:"status_bar"`
}

type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
	Colors   ColorsConfig   `mapstructure:"colors"`
}

func LoadConfig() {
	configPath := filepath.Join(utils.GetHomeDirectory(), ".config", "fm")
	configFile := filepath.Join(utils.GetHomeDirectory(), ".config", "fm", "config.yml")

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
			log.Fatal("Error loading config:", err)
		}
	}
}

func GetConfig() (config Config) {
	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}

func SetDefaults() {
	// App Settings
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("settings.rounded_panes", false)
	viper.SetDefault("settings.enable_logging", false)
	viper.SetDefault("settings.enable_mousewheel", true)

	// DirTree colors
	viper.SetDefault("colors.dir_tree.selected_item", constants.Pink)
	viper.SetDefault("colors.dir_tree.unselected_item", constants.White)

	// Pane colors
	viper.SetDefault("colors.pane.active_border_color", constants.Pink)
	viper.SetDefault("colors.pane.inactive_border_color", constants.White)

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
