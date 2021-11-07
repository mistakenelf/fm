package config

import (
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SettingsConfig struct represents the config for the settings.
type SettingsConfig struct {
	StartDir         string `mapstructure:"start_dir"`
	ShowIcons        bool   `mapstructure:"show_icons"`
	EnableLogging    bool   `mapstructure:"enable_logging"`
	EnableMouseWheel bool   `mapstructure:"enable_mousewheel"`
	PrettyMarkdown   bool   `mapstructure:"pretty_markdown"`
	Borderless       bool   `mapstructure:"borderless"`
	Theme            string `mapstructure:"theme"`
	SyntaxTheme      string `mapstructure:"syntax_theme"`
}

// Config represents the main config for the application.
type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
}

// LoadConfig loads a users config and creates the config if it does not exist
// located at ~/.fm.yml.
func LoadConfig(startDir, selectionPath *pflag.Flag) {
	viper.AddConfigPath("$HOME")
	viper.SetConfigName(".fm")
	viper.SetConfigType("yml")

	// Setup flags.
	err := viper.BindPFlag("start-dir", startDir)
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

	// Setup config defaults.
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("settings.enable_logging", false)
	viper.SetDefault("settings.enable_mousewheel", true)
	viper.SetDefault("settings.pretty_markdown", true)
	viper.SetDefault("settings.borderless", false)
	viper.SetDefault("settings.theme", "default")
	viper.SetDefault("settings.syntax_theme", "dracula")

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
}

// GetConfig returns the users config.
func GetConfig() (config Config) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}
