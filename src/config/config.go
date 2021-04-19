package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	StartDir          string `mapstructure:"start_dir"`
	ShowIcons         bool   `mapstructure:"show_icons"`
	SelectedItemColor string `mapstructure:"selected_item_color"`
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("$HOME/.config/fm")

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
	viper.SetDefault("start_dir", ".")
	viper.SetDefault("show_icons", true)
	viper.SetDefault("selected_item_color", "#F25D94")
}
