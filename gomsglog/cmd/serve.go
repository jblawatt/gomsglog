package cmd

import (
	"github.com/jblawatt/gomsglog/gomsglog/http"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Serves a REST API for this",
	Aliases: []string{"server"},
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.PersistentFlags().GetString("bind")
		http.Serve(addr)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringP("bind", "b", ":1234", "Where to bind to.")
}
