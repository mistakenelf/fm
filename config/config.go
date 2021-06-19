package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/utils"

	"github.com/spf13/viper"
)

// Main app settings
type SettingsConfig struct {
	StartDir         string `mapstructure:"start_dir"`
	ShowIcons        bool   `mapstructure:"show_icons"`
	RoundedPanes     bool   `mapstructure:"rounded_panes"`
	EnableLogging    bool   `mapstructure:"enable_logging"`
	EnableMouseWheel bool   `mapstructure:"enable_mousewheel"`
	PrettyMarkdown   bool   `mapstructure:"pretty_markdown"`
}

// Directory tree has a selected item and unselected item color
type DirTreeColors struct {
	SelectedItem   string `mapstructure:"selected_item"`
	UnselectedItem string `mapstructure:"unselected_item"`
}

// A pane has a active and inactive border color
type PaneColors struct {
	ActiveBorderColor   string `mapstructure:"active_border_color"`
	InactiveBorderColor string `mapstructure:"inactive_border_color"`
}

// Color consists of both a background and foreground
type ColorVariant struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

// Colors for the 4 different parts of the status bar
type StatusBarColors struct {
	SelectedFile ColorVariant `mapstructure:"selected_file"`
	Bar          ColorVariant `mapstructure:"bar"`
	TotalFiles   ColorVariant `mapstructure:"total_files"`
	Logo         ColorVariant `mapstructure:"logo"`
}

// Color config for the different components
type ColorsConfig struct {
	DirTree   DirTreeColors   `mapstructure:"dir_tree"`
	Pane      PaneColors      `mapstructure:"pane"`
	Spinner   string          `mapstructure:"spinner"`
	StatusBar StatusBarColors `mapstructure:"status_bar"`
}

// Main app config
type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
	Colors   ColorsConfig   `mapstructure:"colors"`
}

// Load users config and create the config if it does not exist
func LoadConfig() {
	// Get users home directory and get path to create a config at
	// if it does not exist
	homeDir, _ := utils.GetHomeDirectory()
	configPath := filepath.Join(homeDir, ".config", "fm")
	configFile := filepath.Join(homeDir, ".config", "fm", "config.yml")

	// Set viper config name, type and path
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	// if a config file already exists, do nothing, else create a new config
	// file in the .config/fm directory
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configFile), 0770); err != nil {
			log.Fatal("Error creating config file")
		}

		// Create the config and write it to viper
		os.Create(configFile)
		viper.WriteConfig()
	}
}

// Get the users config and return it for use
func GetConfig() (config Config) {
	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}

// Setup viper defaults
func SetDefaults() {
	// App Settings
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("settings.rounded_panes", false)
	viper.SetDefault("settings.enable_logging", false)
	viper.SetDefault("settings.enable_mousewheel", true)
	viper.SetDefault("settings.pretty_markdown", true)

	// DirTree colors
	viper.SetDefault("colors.dir_tree.selected_item", constants.Pink)
	viper.SetDefault("colors.dir_tree.unselected_item", constants.White)

	// Pane colors
	viper.SetDefault("colors.pane.active_border_color", constants.Pink)
	viper.SetDefault("colors.pane.inactive_border_color", constants.White)

	// Spinner colors
	viper.SetDefault("colors.spinner", constants.Pink)

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
