package config

import (
	"log"
	"os"

	"github.com/knipferrc/fm/internal/constants"

	"github.com/spf13/viper"
)

// SettingsConfig struct represents the config for the settings.
type SettingsConfig struct {
	StartDir         string `mapstructure:"start_dir"`
	ShowIcons        bool   `mapstructure:"show_icons"`
	RoundedPanes     bool   `mapstructure:"rounded_panes"`
	EnableLogging    bool   `mapstructure:"enable_logging"`
	EnableMouseWheel bool   `mapstructure:"enable_mousewheel"`
	PrettyMarkdown   bool   `mapstructure:"pretty_markdown"`
}

//DirTreeColors struct represents the colors for the dirtree.
type DirTreeColors struct {
	SelectedItem   string `mapstructure:"selected_item"`
	UnselectedItem string `mapstructure:"unselected_item"`
}

// PaneColors represents the colors for a pane.
type PaneColors struct {
	ActiveBorderColor   string `mapstructure:"active_border_color"`
	InactiveBorderColor string `mapstructure:"inactive_border_color"`
}

// ColorVariant struct represents a color.
type ColorVariant struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

// StatusBarColors represents the colors for the status bar.
type StatusBarColors struct {
	SelectedFile ColorVariant `mapstructure:"selected_file"`
	Bar          ColorVariant `mapstructure:"bar"`
	TotalFiles   ColorVariant `mapstructure:"total_files"`
	Logo         ColorVariant `mapstructure:"logo"`
}

// ColorsConfig struct represets the colors of the UI.
type ColorsConfig struct {
	DirTree   DirTreeColors   `mapstructure:"dir_tree"`
	Pane      PaneColors      `mapstructure:"pane"`
	Spinner   string          `mapstructure:"spinner"`
	StatusBar StatusBarColors `mapstructure:"status_bar"`
}

// Config represents the main config for the application.
type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
	Colors   ColorsConfig   `mapstructure:"colors"`
}

// LoadConfig loads a users config and creates the config if it does not exist
// located at ~/.fm.yml.
func LoadConfig() {
	viper.AddConfigPath("$HOME")
	viper.SetConfigName(".fm")
	viper.SetConfigType("yml")

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfig()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		} else {
			log.Fatal(err)
		}
	}
}

// GetConfig returns the users config.
func GetConfig() (config Config) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}

// SetDefaults sets default values for the config.
func SetDefaults() {
	// App Settings.
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("settings.rounded_panes", false)
	viper.SetDefault("settings.enable_logging", false)
	viper.SetDefault("settings.enable_mousewheel", true)
	viper.SetDefault("settings.pretty_markdown", true)

	// DirTree colors.
	viper.SetDefault("colors.dir_tree.selected_item", constants.Pink)
	viper.SetDefault("colors.dir_tree.unselected_item", constants.White)

	// Pane colors.
	viper.SetDefault("colors.pane.active_border_color", constants.Pink)
	viper.SetDefault("colors.pane.inactive_border_color", constants.White)

	// Spinner colors.
	viper.SetDefault("colors.spinner", constants.Pink)

	// StatusBar colors.
	viper.SetDefault("colors.status_bar.selected_file.foreground", constants.White)
	viper.SetDefault("colors.status_bar.selected_file.background", constants.Pink)
	viper.SetDefault("colors.status_bar.bar.foreground", constants.White)
	viper.SetDefault("colors.status_bar.bar.background", constants.DarkGray)
	viper.SetDefault("colors.status_bar.total_files.foreground", constants.White)
	viper.SetDefault("colors.status_bar.total_files.background", constants.LightPurple)
	viper.SetDefault("colors.status_bar.logo.foreground", constants.White)
	viper.SetDefault("colors.status_bar.logo.background", constants.DarkPurple)
}
