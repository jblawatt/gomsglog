package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jblawatt/gomsglog/gomsglog/cmd"
	"github.com/spf13/viper"
)

func main() {

	viper.SetEnvPrefix("ML_")
	viper.SetConfigName(".mlrc")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	viper.SetDefault("database.dialect", "sqlite3")
	viper.SetDefault("database.connectionstring", "db.sqlite3")
	viper.SetDefault("database.debug", false)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func setupLog() {
	logfile := viper.GetString("logfile")

	if logfile != "" {
		f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}
}
