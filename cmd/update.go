package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/muesli/coral"
)

var updateCmd = &coral.Command{
	Use:   "update",
	Short: "Update FM to the latest version",
	Long:  `Update FM to the latest version.`,
	Run: func(cmd *coral.Command, args []string) {
		updateCommand := exec.Command("bash", "-c", "curl -sfL https://raw.githubusercontent.com/knipferrc/fm/main/install.sh | sh")
		updateCommand.Stdin = os.Stdin
		updateCommand.Stdout = os.Stdout
		updateCommand.Stderr = os.Stderr

		err := updateCommand.Run()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	},
}
