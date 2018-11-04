package cmd

import (
	"fmt"
	"os"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/jblawatt/gomsglog/gomsglog/parsers"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a not message.",
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		raw, _ := cmd.PersistentFlags().GetString("message")
		if raw != "" {
			msg := parsers.NewMessage(raw)
			gomsglog.Persist(msg)
		} else {
			fmt.Fprint(os.Stderr, "Please provide a mesage (-m/--message)")
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringP("message", "m", "", "message to add")
}
