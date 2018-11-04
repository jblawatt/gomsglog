package parsers

import (
	"plugin"
	"strings"

	"github.com/satori/go.uuid"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type Message struct {
	ID           string
	Original     string
	HTML         string
	RelatedUsers []string
	Tags         []string
	Attributes   map[string]Attr
	URLs         []string
}

type ReplacementManager map[string]string

func (rm ReplacementManager) Add(key string, value string) {
	rm[key] = value
}

func NewUUID() string {
	return uuid.NewV4().String()
}

type MessageParser interface {
	Parse(m *Message, repl *ReplacementManager) error
}

var parserRegistry []MessageParser

func GetRegisteredParsers() []MessageParser {
	return parserRegistry
}

func RegisterParser(parser MessageParser) {
	parserRegistry = append(parserRegistry, parser)
}

func ApplyToMessage(input string, message *Message) {

	rm := make(ReplacementManager)

	parserPlugins := viper.GetStringSlice("parsers")
	for _, p := range parserPlugins {
		plg, err := plugin.Open(p)
		if err != nil {
			jww.ERROR.Printf("Error loading plugin %s: %s\n", p, err)
			continue
		}
		parseFunc, fnErr := plg.Lookup("Parse")
		if fnErr != nil {
			jww.ERROR.Printf("Error loading parse func of %s: %s", p, fnErr.Error())
			continue
		}
		perr := parseFunc.(func(*Message, *ReplacementManager) error)(message, &rm)
		if perr != nil {
			jww.ERROR.Printf("Error parsing in %s: %s", p, perr.Error())
			continue
		}
	}

	for _, parser := range GetRegisteredParsers() {
		err := parser.Parse(message, &rm)
		if err != nil {
			jww.ERROR.Printf("Error parsing: %s\n", err.Error())
		}
	}

	for key, value := range rm {
		message.HTML = strings.Replace(message.HTML, key, value, 1)
	}

}

func UpdateMessage(input string, message *Message) {
	message.Attributes = make(map[string]Attr)
	message.Original = input
	message.HTML = input
	message.Tags = make([]string, 0)
	message.URLs = make([]string, 0)
	ApplyToMessage(input, message)
}

func NewMessage(input string) Message {
	message := Message{
		ID:         NewUUID(),
		Original:   input,
		HTML:       input,
		Attributes: make(map[string]Attr),
		Tags:       make([]string, 0),
		URLs:       make([]string, 0),
	}
	ApplyToMessage(input, &message)
	return message
}
