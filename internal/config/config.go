package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// AppDir is the name of the directory where the config file is stored.
const AppDir = "fm"

// ConfigFileName is the name of the config file that gets created.
const ConfigFileName = "config.yml"

// SyntaxThemeConfig represents light and dark syntax themes.
type SyntaxThemeConfig struct {
	Light string `yaml:"light"`
	Dark  string `yaml:"dark"`
}

// SettingsConfig struct represents the config for the settings.
type SettingsConfig struct {
	StartDir       string `yaml:"start_dir"`
	EnableLogging  bool   `yaml:"enable_logging"`
	PrettyMarkdown bool   `yaml:"pretty_markdown"`
}

// ThemeConfig represents the config for themes.
type ThemeConfig struct {
	AppTheme    string            `yaml:"app_theme"`
	SyntaxTheme SyntaxThemeConfig `yaml:"syntax_theme"`
}

// Config represents the main config for the application.
type Config struct {
	Settings SettingsConfig `yaml:"settings"`
	Theme    ThemeConfig    `yaml:"theme"`
}

// configError represents an error that occurred while parsing the config file.
type configError struct {
	configDir string
	parser    ConfigParser
	err       error
}

// ConfigParser is the parser for the config file.
type ConfigParser struct{}

// getDefaultConfig returns the default config for the application.
func (parser ConfigParser) getDefaultConfig() Config {
	return Config{
		Settings: SettingsConfig{
			StartDir:       ".",
			EnableLogging:  false,
			PrettyMarkdown: true,
		},
		Theme: ThemeConfig{
			AppTheme: "default",
			SyntaxTheme: SyntaxThemeConfig{
				Dark:  "dracula",
				Light: "pygments",
			},
		},
	}
}

// getDefaultConfigYamlContents returns the default config file contents.
func (parser ConfigParser) getDefaultConfigYamlContents() string {
	defaultConfig := parser.getDefaultConfig()
	yaml, _ := yaml.Marshal(defaultConfig)

	return string(yaml)
}

// Error returns the error message for when a config file is not found.
func (e configError) Error() string {
	return fmt.Sprintf(
		`Couldn't find a config.yml configuration file.
Create one under: %s
Example of a config.yml file:
%s
For more info, go to https://github.com/mistakenelf/fm
press q to exit.
Original error: %v`,
		path.Join(e.configDir, AppDir, ConfigFileName),
		e.parser.getDefaultConfigYamlContents(),
		e.err,
	)
}

// writeDefaultConfigContents writes the default config file contents to the given file.
func (parser ConfigParser) writeDefaultConfigContents(newConfigFile *os.File) error {
	_, err := newConfigFile.WriteString(parser.getDefaultConfigYamlContents())

	if err != nil {
		return err
	}

	return nil
}

// createConfigFileIfMissing creates the config file if it doesn't exist.
func (parser ConfigParser) createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		newConfigFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}

		defer newConfigFile.Close()
		return parser.writeDefaultConfigContents(newConfigFile)
	}

	return nil
}

// getConfigFileOrCreateIfMissing returns the config file path or creates the config file if it doesn't exist.
func (parser ConfigParser) getConfigFileOrCreateIfMissing() (*string, error) {
	var err error
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return nil, configError{parser: parser, configDir: configDir, err: err}
		}
	}

	prsConfigDir := filepath.Join(configDir, AppDir)
	err = os.MkdirAll(prsConfigDir, os.ModePerm)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}

	configFilePath := filepath.Join(prsConfigDir, ConfigFileName)
	err = parser.createConfigFileIfMissing(configFilePath)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}

	return &configFilePath, nil
}

// parsingError represents an error that occurred while parsing the config file.
type parsingError struct {
	err error
}

// Error represents an error that occurred while parsing the config file.
func (e parsingError) Error() string {
	return fmt.Sprintf("failed parsing config.yml: %v", e.err)
}

// readConfigFile reads the config file and returns the config.
func (parser ConfigParser) readConfigFile(path string) (Config, error) {
	config := parser.getDefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return config, configError{parser: parser, configDir: path, err: err}
	}

	err = yaml.Unmarshal((data), &config)
	return config, err
}

// initParser initializes the parser.
func initParser() ConfigParser {
	return ConfigParser{}
}

// ParseConfig parses the config file and returns the config.
func ParseConfig() (Config, error) {
	var config Config
	var err error

	parser := initParser()

	configFilePath, err := parser.getConfigFileOrCreateIfMissing()
	if err != nil {
		return config, parsingError{err: err}
	}

	config, err = parser.readConfigFile(*configFilePath)
	if err != nil {
		return config, parsingError{err: err}
	}

	return config, nil
}
