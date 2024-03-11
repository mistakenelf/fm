package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/fm/filesystem"
	"github.com/mistakenelf/fm/internal/theme"
	"github.com/mistakenelf/fm/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "fm",
	Short:   "FM is a simple, configurable, and fun to use file manager",
	Version: "1.0.0",
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

		enableLogging, err := cmd.Flags().GetBool("enable-logging")
		if err != nil {
			log.Fatal(err)
		}

		prettyMarkdown, err := cmd.Flags().GetBool("pretty-markdown")
		if err != nil {
			log.Fatal(err)
		}

		applicationTheme, err := cmd.Flags().GetString("theme")
		if err != nil {
			log.Fatal(err)
		}

		// If logging is enabled, logs will be output to debug.log.
		if enableLogging {
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

		appTheme := theme.GetTheme(applicationTheme)

		cfg := tui.Config{
			StartDir:       startDir,
			SelectionPath:  selectionPath,
			EnableLogging:  enableLogging,
			PrettyMarkdown: prettyMarkdown,
			Theme:          appTheme,
		}

		m := tui.New(cfg)

		p := tea.NewProgram(m, tea.WithAltScreen())
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
	rootCmd.PersistentFlags().String("start-dir", filesystem.CurrentDirectory, "Starting directory for FM")
	rootCmd.PersistentFlags().Bool("enable-logging", false, "Enable logging for FM")
	rootCmd.PersistentFlags().Bool("pretty-markdown", true, "Render markdown to look nice")
	rootCmd.PersistentFlags().String("theme", "default", "Application theme")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
