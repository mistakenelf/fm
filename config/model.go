package config

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

type SelectedFileColors struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

type BarColors struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

type TotalFilesColors struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

type LogoColors struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

type StatusBarColors struct {
	SelectedFile SelectedFileColors `mapstructure:"selected_file"`
	Bar          BarColors          `mapstructure:"bar"`
	TotalFiles   TotalFilesColors   `mapstructure:"total_files"`
	Logo         LogoColors         `mapstructure:"logo"`
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
