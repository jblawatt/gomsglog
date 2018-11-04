package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/jblawatt/gomsglog/gomsglog/parsers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Edit an existing message.",
	Aliases: []string{"e"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "'%s' is an invalid message id", arg)
				os.Exit(1)
			}
			message, found := gomsglog.LoadMessage(id)
			if !found {
				fmt.Fprintf(os.Stderr, "Invalid message id %d\n", id)
				os.Exit(1)
			}
			raw, _ := cmd.PersistentFlags().GetString("message")
			if raw == "" {
				file, _ := ioutil.TempFile(os.TempDir(), "gomsglog_")
				file.WriteString(message.Original)
				filename := file.Name()
				file.Close()
				ecmd := exec.Command(viper.GetString("editor"), filename)
				ecmd.Stdout = os.Stdout
				ecmd.Stderr = os.Stderr
				ecmd.Stdin = os.Stdin
				ecmd.Run()
				rawBytes, _ := ioutil.ReadFile(filename)
				raw = string(rawBytes)
			}
			msg := parsers.NewMessage(raw)
			gomsglog.Update(id, msg)
		}

	},
}

func init() {
	RootCmd.AddCommand(editCmd)
	editCmd.PersistentFlags().StringP("message", "m", "", "message to set.")
}
