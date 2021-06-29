package app

import (
	"fmt"
	"log"
	"os"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "fm",
	Short:   "FM is a simple, configurable and fun to use file manager",
	Version: "v0.0.5",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup default config and load the users config
		config.SetDefaults()
		config.LoadConfig()

		cfg := config.GetConfig()

		// If logging is enabled, logs will be output to debug.log
		if cfg.Settings.EnableLogging {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			defer f.Close()
		}

		m := ui.NewModel()
		var opts []tea.ProgramOption

		// Always append alt screen program option
		opts = append(opts, tea.WithAltScreen())

		// If mousewheel is enabled, append it to the program options
		if cfg.Settings.EnableMouseWheel {
			opts = append(opts, tea.WithMouseAllMotion())
		}

		// Initialize new app
		p := tea.NewProgram(m, opts...)

		// If the program fails to start, exit the program
		if err := p.Start(); err != nil {
			log.Fatal("Failed to start fm", err)
			os.Exit(1)
		}
	},
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
