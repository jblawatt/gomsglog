package cmd

import (
	"fmt"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:     "tags",
	Aliases: []string{"t", "tag"},
	Short:   "List all tags",
	Run: func(cmd *cobra.Command, args []string) {
		for _, t := range gomsglog.LoadTags() {
			fmt.Println(t.Slug)
		}
	},
}

func init() {
	RootCmd.AddCommand(tagsCmd)
}
