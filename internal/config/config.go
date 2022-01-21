package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/knipferrc/fm/dirfs"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SyntaxThemeConfig represents light and dark syntax themes.
type SyntaxThemeConfig struct {
	Light string `mapstructure:"light"`
	Dark  string `mapstructure:"dark"`
}

// SettingsConfig struct represents the config for the settings.
type SettingsConfig struct {
	StartDir         string `mapstructure:"start_dir"`
	ShowIcons        bool   `mapstructure:"show_icons"`
	EnableLogging    bool   `mapstructure:"enable_logging"`
	EnableMouseWheel bool   `mapstructure:"enable_mousewheel"`
	PrettyMarkdown   bool   `mapstructure:"pretty_markdown"`
	Borderless       bool   `mapstructure:"borderless"`
	SimpleMode       bool   `mapstructure:"simple_mode"`
}

// ThemeConfig represents the config for themes.
type ThemeConfig struct {
	AppTheme    string            `mapstructure:"app_theme"`
	SyntaxTheme SyntaxThemeConfig `mapstructure:"syntax_theme"`
}

// Config represents the main config for the application.
type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
	Theme    ThemeConfig    `mapstructure:"theme"`
}

// LoadConfig loads a users config and creates the config if it does not exist
// located at ~/.config/fm.yml.
func LoadConfig(startDir, selectionPath *pflag.Flag) {
	var err error

	if runtime.GOOS != "windows" {
		homeDir, err := dirfs.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		err = dirfs.CreateDirectory(filepath.Join(homeDir, ".config", "fm"))
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath("$HOME/.config/fm")
	} else {
		viper.AddConfigPath("$HOME")
	}

	viper.SetConfigName("fm")
	viper.SetConfigType("yml")

	// Setup config defaults.
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("settings.enable_logging", false)
	viper.SetDefault("settings.enable_mousewheel", true)
	viper.SetDefault("settings.pretty_markdown", true)
	viper.SetDefault("settings.borderless", false)
	viper.SetDefault("settings.syntax_theme", "default")
	viper.SetDefault("theme.app_theme", "default")
	viper.SetDefault("theme.syntax_theme.light", "pygments")
	viper.SetDefault("theme.syntax_theme.dark", "dracula")
	viper.SetDefault("settings.simple_mode", false)

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
		}
	}

	// Setup flags.
	err = viper.BindPFlag("start-dir", startDir)
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlag("selection-path", selectionPath)
	if err != nil {
		log.Fatal(err)
	}

	// Setup flag defaults.
	viper.SetDefault("start-dir", "")
	viper.SetDefault("selection-path", "")
}

// GetConfig returns the users config.
func GetConfig() (config Config) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}
