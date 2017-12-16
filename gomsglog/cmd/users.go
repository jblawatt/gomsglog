package cmd

import (
	"fmt"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"u", "user"},
	Short:   "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		for _, u := range gomsglog.LoadUsers() {
			fmt.Println(u.Slug)
		}
	},
}

func init() {
	RootCmd.AddCommand(usersCmd)
}
