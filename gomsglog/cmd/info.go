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
		fmt.Print("Options :: ")
		fmt.Printf("DEBUG = %t; ", viper.GetBool("debug"))
		fmt.Printf("BIND = %s; ", viper.GetString("bind"))
		fmt.Printf("LOG.TEMPLATE = %d; ", viper.GetInt("log.limit"))
		fmt.Printf("LOG.OFFSET = %d; ", viper.GetInt("log.offset"))
		fmt.Printf("LOGLEVEL = %s; ", viper.GetString("LOGLEVEL"))
		fmt.Printf("LOGFILE = %s; ", viper.GetString("logfile"))
		fmt.Printf("DATABASE.DIALECT = %s; ", viper.GetString("database.dialect"))
		fmt.Printf("DATABASE.CONNECTIONSTRING = %s; ", viper.GetString("database.connectionstring"))
		fmt.Printf("DATABASE.DEBUG = %t; ", viper.GetBool("database.debug"))
		fmt.Printf("MLRC = %s", viper.ConfigFileUsed())

	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
