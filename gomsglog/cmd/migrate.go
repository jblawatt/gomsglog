package cmd

import (
	"fmt"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/spf13/cobra"
)

var migrCmd = &cobra.Command{
	Use:     "_migrate",
	Aliases: []string{"_mg"},
	Run: func(cmd *cobra.Command, args []string) {
		gomsglog.AutoMigrate()
		fmt.Println("Ok")
	},
}

func init() {
	RootCmd.AddCommand(migrCmd)
}
