package cmd

import (
	"strconv"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Removes a specific message",
	Aliases: []string{"remove"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, rawID := range args {
			id, err := strconv.Atoi(rawID)
			if err != nil {
				panic(err)
			}
			gomsglog.DeleteMessage(id)
		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
