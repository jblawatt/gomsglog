package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jblawatt/gomsglog/gomsglog/cmd"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

func main() {
	setupViper()
	setupLog()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func setupViper() {
	viper.SetEnvPrefix("ML_")
	viper.AutomaticEnv()
	viper.BindEnv("debug", "ML_DEBUG")
	viper.SetConfigName(".mlrc")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	viper.SetDefault("database.dialect", "sqlite3")
	viper.SetDefault("database.connectionstring", "db.sqlite3")
	viper.SetDefault("database.debug", false)
	viper.SetDefault("loglevel", "WARN")
	viper.SetDefault("debug", false)

	fmt.Println(color.RedString("DEBUG ", viper.GetBool("debug")))
}

var thresholds = []jww.Threshold{
	jww.LevelTrace,
	jww.LevelDebug,
	jww.LevelInfo,
	jww.LevelWarn,
	jww.LevelError,
	jww.LevelFatal,
	jww.LevelCritical,
}

func setupLog() {
	logfile := viper.GetString("logfile")
	level := viper.GetString("loglevel")
	level = strings.ToUpper(level)
	for _, t := range thresholds {
		if t.String() == level {
			jww.SetLogThreshold(t)
		}
	}
	if logfile != "" {
		f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(f)
		jww.SetLogOutput(f)
	} else {
		log.SetOutput(os.Stdout)
		jww.SetLogOutput(os.Stdout)
	}
}
