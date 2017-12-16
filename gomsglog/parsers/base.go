package parsers

import (
	"strings"

	"github.com/satori/go.uuid"
)

type Message struct {
	ID           string
	Original     string
	HTML         string
	RelatedUsers []string
	Tags         []string
	Attributes   map[string]interface{}
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
	Parse(m *Message, repl *ReplacementManager)
}

var parserRegistry []MessageParser

func GetRegisteredParsers() []MessageParser {
	return parserRegistry
}

func RegisterParser(parser MessageParser) {
	parserRegistry = append(parserRegistry, parser)
}

func NewMessage(input string) Message {
	message := Message{
		ID:         NewUUID(),
		Original:   input,
		HTML:       input,
		Attributes: make(map[string]interface{}),
		Tags:       make([]string, 0),
		URLs:       make([]string, 0),
	}
	// replacements := make(map[string]string)
	rm := make(ReplacementManager)

	for _, parser := range GetRegisteredParsers() {
		parser.Parse(&message, &rm)
	}

	for key, value := range rm {
		message.HTML = strings.Replace(message.HTML, key, value, 1)
	}

	return message
}
