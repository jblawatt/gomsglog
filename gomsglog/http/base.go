package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Run starts the program.
func Serve(addr string) {

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessageHandler).Methods("POST")
	router.HandleFunc("/api/messages", GetMessagesHandler).Methods("GET")
	router.HandleFunc("/api/messages/{messageID}", GetMessageHandler).Methods("GET")
	router.HandleFunc("/api/messages/{messageID}", DeleteMessageHandler).Methods("DELETE")
	router.HandleFunc("/api/tags", GetTagsHandler).Methods("GET")
	router.PathPrefix("/static").Handler(http.FileServer(http.Dir("./")))
	router.HandleFunc("/", IndexHandler).Methods("GET")
	router.HandleFunc("/ws", WebSocketHandler)

	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	fmt.Printf("Listening on %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, loggingRouter))

}
