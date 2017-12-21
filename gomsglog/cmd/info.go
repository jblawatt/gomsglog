package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use: "info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DEBUG =", viper.GetBool("debug"))
		fmt.Println("")
		fmt.Println("LOGLEVEL =", viper.GetString("LOGLEVEL"))
		fmt.Println("LOGFILE =", viper.GetString("logfile"))
		fmt.Println("")
		fmt.Println("DATABASE.DIALECT =", viper.GetString("database.dialect"))
		fmt.Println("DATABASE.CONNECTIONSTRING =", viper.GetString("database.connectionstring"))
		fmt.Println("DATABASE.DEBUG =", viper.GetBool("database.debug"))
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
