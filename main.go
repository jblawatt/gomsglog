package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/jblawatt/gomsglog/gomsglog"

	"github.com/jblawatt/gomsglog/gomsglog/cmd"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	binary   = "ml"
	version  = "dev"
	builDate = time.Now().Format("2006-01-02")
)

func main() {
	setupViper()
	setVersionInfo(binary, version, builDate)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func setVersionInfo(binary string, version string, build string) {
	gomsglog.SetBinaryInfo(binary)
	gomsglog.SetVersionInfo(version)
	gomsglog.SetBuildDateInfo(builDate)
}

func setupViper() {
	viper.SetEnvPrefix("ML_")
	viper.AutomaticEnv()
	viper.BindEnv("debug", "ML_DEBUG")
	viper.SetConfigName(".mlrc")
	viper.AddConfigPath(".")
	viper.AddConfigPath(os.Getenv("HOME"))
	if runtime.GOOS == "windows" {
		viper.AddConfigPath(path.Join(os.Getenv("APPDATA"), "ml"))
	}
	viper.ReadInConfig()

	viper.SetDefault("database.dialect", "sqlite3")
	viper.SetDefault("database.connectionstring", "db.sqlite3")
	viper.SetDefault("database.debug", false)
	viper.SetDefault("loglevel", "WARN")
	viper.SetDefault("debug", false)
	if runtime.GOOS == "windows" {
		viper.SetDefault("editor", "notepad")
	}
	if runtime.GOOS == "linux" {
		viper.SetDefault("editor", "vim")
	}

	setupLog()

	jww.DEBUG.Printf("Using config file %s.", viper.ConfigFileUsed())

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
	if logfile != "" || logfile == "stdout" {
		f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(f)
		jww.SetLogOutput(f)
	} else {
		log.SetOutput(os.Stdout)
		jww.SetLogOutput(os.Stdout)
	}

}
