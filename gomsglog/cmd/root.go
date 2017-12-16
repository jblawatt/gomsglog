package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use: "ml",
}

func init() {
	f := RootCmd.PersistentFlags()
	f.Bool("debug", false, "DEBUG")
	viper.BindPFlag("debug", f.Lookup("debug"))
}
