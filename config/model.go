package config

type SettingsConfig struct {
	StartDir     string `mapstructure:"start_dir"`
	ShowIcons    bool   `mapstructure:"show_icons"`
	ShowHidden   bool   `mapstructure:"show_hidden"`
	RoundedPanes bool   `mapstructure:"rounded_panes"`
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
