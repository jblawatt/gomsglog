package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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

func NewMessage(input string) Message {
	message := Message{
		ID:         newUUID(),
		Original:   input,
		HTML:       input,
		Attributes: make(map[string]interface{}),
		Tags:       make([]string, 0),
		URLs:       make([]string, 0),
	}
	replacements := make(map[string]string)

	HandleUsers(&message, &replacements)
	HandleTags(&message, &replacements)
	HandleLinks(&message, &replacements)
	HandleAttrs(&message, &replacements)

	for key, value := range replacements {
		message.HTML = strings.Replace(message.HTML, key, value, 1)
	}

	return message
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{messageID}", GetMessage).Methods("GET")

	// router.HandleFunc("/api/tags",).Methods("GET")
	// router.HandleFunc("/api/tags",).Methods("PATCH")

	log.Fatal(http.ListenAndServe("127.0.0.1:12345", router))

}
