package cmd

import (
	"os"
	"text/template"

	"github.com/fatih/color"
	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var COMPLETE_TEMPLATE = `
ID:      {{.ID}}
Date:    {{.Created.Format "2006-01-02 15:04"}}
{{ "Message:" | hired }} {{.Original | hired }}
Tags:    {{range .Tags}}{{.ScreenName}} {{end}}
Users:   {{range .RelatedUsers}}{{.ScreenName}} {{end}}
Attrs:   {{range .Attributes}}{{.ScreenName}}={{.String}} {{end}}
URLs: {{range .URLs}}
  - {{.URL}}{{end}}
`

var DEFAULT_TEMPLATE = `
ID:      {{.ID}}
Date:    {{.Created.Format "2006-01-02 15:04"}}
{{ "Message:" | hired }} {{.Original | hired }}
`

var SHORT_TEMPLATE = `
{{.ID}}: {{.Original}}
`

var templates = map[string]*template.Template{
	"default":  template.New("default"),
	"short":    template.New("short"),
	"json":     template.New("json"),
	"complete": template.New("complete"),
}

var logCmd = &cobra.Command{
	Use:     "log",
	Short:   "Lists all messages.",
	Aliases: []string{"l", "ls"},
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.PersistentFlags()
		limit, _ := flags.GetInt("limit")
		templ := viper.GetString("log.template")
		tags, _ := flags.GetStringArray("tag")
		users, _ := flags.GetStringArray("user")
		messages := gomsglog.LoadMessages(limit, 0, tags, users)
		for i := len(messages) - 1; i >= 0; i-- {
			msg := messages[i]
			templates[templ].Execute(os.Stdout, msg)
		}
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
	flags := logCmd.PersistentFlags()
	flags.IntP("limit", "l", 100, "Number of Entries")
	flags.StringP("template", "T", "default", "Template")
	flags.StringArrayP("user", "u", []string{}, "Users to filter")
	flags.StringArrayP("tag", "t", []string{}, "Tags to filter")

	viper.BindPFlag("log.template", flags.Lookup("template"))

	funcMap := template.FuncMap{
		"green":     color.GreenString,
		"red":       color.RedString,
		"yellow":    color.YellowString,
		"blue":      color.BlueString,
		"magenta":   color.MagentaString,
		"cyan":      color.CyanString,
		"white":     color.WhiteString,
		"black":     color.BlackString,
		"higreen":   color.HiGreenString,
		"hired":     color.HiRedString,
		"hiyellow":  color.HiYellowString,
		"hiblue":    color.HiBlueString,
		"himagenta": color.HiMagentaString,
		"hicyan":    color.HiCyanString,
		"hiwhite":   color.HiWhiteString,
		"hiblack":   color.HiBlackString,
	}

	templates["default"].Funcs(funcMap)
	templates["short"].Funcs(funcMap)
	templates["complete"].Funcs(funcMap)

	if _, err := templates["default"].Parse(DEFAULT_TEMPLATE); err != nil {
		panic(err)
	}

	if _, err := templates["complete"].Parse(COMPLETE_TEMPLATE); err != nil {
		panic(err)
	}

	if _, err := templates["short"].Parse(SHORT_TEMPLATE); err != nil {
		panic(err)
	}

}
