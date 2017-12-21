package gomsglog

import (
	"time"
)

var (
	binary    = "ml.exe"
	version   = "dev"
	builddate = time.Now().Format("2006-01-02")
)

func SetVersionInfo(v string) {
	version = v
}

func SetBinaryInfo(b string) {
	binary = b
}

func SetBuildDateInfo(b string) {
	builddate = b
}

func GetVersion() string {
	return version
}

func GetBinary() string {
	return binary
}

func GetBuildDate() string {
	return builddate
}
