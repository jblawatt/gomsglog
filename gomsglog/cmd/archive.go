package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:     "archive",
	Short:   "Archive message.",
	Aliases: []string{"arch"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Please provide a message id")
			os.Exit(1)
		}
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "'%s' is an invalid message id", arg)
				os.Exit(1)
			}
			_, found := gomsglog.LoadMessage(id)
			if !found {
				fmt.Fprintf(os.Stderr, "Invalid message id %d\n", id)
				os.Exit(1)
			}
			gomsglog.Archive(id)
			fmt.Fprintf(os.Stdout, "Message %d archived.\n", id)
		}
	},
}

func init() {
	RootCmd.AddCommand(archiveCmd)
}
