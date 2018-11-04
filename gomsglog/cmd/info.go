package cmd

import (
	"fmt"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"version"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(
			gomsglog.GetBinary(),
			"- Version:",
			gomsglog.GetVersion(),
			"- Date:",
			gomsglog.GetBuildDate(),
		)
		fmt.Print("Options :: \n")
		fmt.Printf("DEBUG = %t\n", viper.GetBool("debug"))
		fmt.Printf("BIND = %s\n", viper.GetString("bind"))
		fmt.Printf("LOG.TEMPLATE = %d\n", viper.GetInt("log.limit"))
		fmt.Printf("LOG.OFFSET = %d\n", viper.GetInt("log.offset"))
		fmt.Printf("LOGLEVEL = %s\n", viper.GetString("LOGLEVEL"))
		fmt.Printf("LOGFILE = %s\n", viper.GetString("logfile"))
		fmt.Printf("DATABASE.DIALECT = %s\n", viper.GetString("database.dialect"))
		fmt.Printf("DATABASE.CONNECTIONSTRING = %s\n", viper.GetString("database.connectionstring"))
		fmt.Printf("DATABASE.DEBUG = %t\n", viper.GetBool("database.debug"))
		fmt.Printf("MLRC = %s\n", viper.ConfigFileUsed())
		fmt.Printf("EDITOR = %s\n", viper.GetString("editor"))

	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
