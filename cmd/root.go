package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mistakenelf/fm/internal/config"
	"github.com/mistakenelf/fm/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "fm",
	Short:   "FM is a simple, configurable, and fun to use file manager",
	Version: "0.15.10",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startDir, err := cmd.Flags().GetString("start-dir")
		if err != nil {
			log.Fatal(err)
		}

		selectionPath, err := cmd.Flags().GetString("selection-path")
		if err != nil {
			log.Fatal(err)
		}

		cfg, err := config.ParseConfig()
		if err != nil {
			log.Fatal(err)
		}

		// If logging is enabled, logs will be output to debug.log.
		if cfg.Settings.EnableLogging {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err = f.Close(); err != nil {
					log.Fatal(err)
				}
			}()
		}

		if startDir == "" {
			startDir = cfg.Settings.StartDir
		}

		m := tui.New(startDir, selectionPath)
		var opts []tea.ProgramOption

		// Always append alt screen program option.
		opts = append(opts, tea.WithAltScreen())

		// Initialize and start app.
		p := tea.NewProgram(m, opts...)
		if _, err := p.Run(); err != nil {
			log.Fatal("Failed to start fm", err)
			os.Exit(1)
		}
	},
}

// Execute runs the root command and starts the application.
func Execute() {
	rootCmd.AddCommand(updateCmd)
	rootCmd.PersistentFlags().String("selection-path", "", "Path to write to file on open.")
	rootCmd.PersistentFlags().String("start-dir", "", "Starting directory for FM")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
