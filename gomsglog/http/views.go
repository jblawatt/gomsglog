package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jblawatt/gomsglog/gomsglog"
	"github.com/jblawatt/gomsglog/gomsglog/parsers"
)

type InputMessage struct {
	Message string `json:"message"`
}

func WriteNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&map[string]interface{}{
		"error":  http.StatusText(http.StatusNotFound),
		"status": http.StatusNotFound,
	})
}

func WriteOK(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	if message == "" {
		message = http.StatusText(http.StatusOK)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"status":  http.StatusOK,
	})
}

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	input := new(InputMessage)

	json.NewDecoder(r.Body).Decode(&input)
	msg := parsers.NewMessage(input.Message)
	mm := gomsglog.Persist(msg)
	json.NewEncoder(w).Encode(&mm)

}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	tags := []string{}
	if tagsStr := r.URL.Query().Get("tags"); tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
	}

	users := []string{}
	if usersStr := r.URL.Query().Get("users"); usersStr != "" {
		users = strings.Split(usersStr, ",")
	}

	attrs := []string{}
	if attrsStr := r.URL.Query().Get("attrs"); attrsStr != "" {
		attrs = strings.Split(attrsStr, ",")
	}

	messages := gomsglog.LoadMessages(
		limit,
		offset,
		tags,
		users,
		attrs,
	)
	json.NewEncoder(w).Encode(&messages)
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	messageID, _ := strconv.Atoi(vars["messageID"])
	message, found := gomsglog.LoadMessage(messageID)
	if found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&message)
	} else {
		WriteNotFound(w)
	}
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	messageID, _ := strconv.Atoi(vars["messageID"])
	found := gomsglog.DeleteMessage(messageID)
	if found {
		WriteOK(w, "Message deleted")
	} else {
		WriteNotFound(w)
	}
}

func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	db := gomsglog.GetDB()
	var slugs []string
	db.Model(&gomsglog.TagModel{}).Pluck(`DISTINCT "slug"`, &slugs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&slugs)
}

func GetAttrsHandler(w http.ResponseWriter, r *http.Request) {

}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := ioutil.ReadFile("templates/index.html")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(t))
}

var updater = websocket.Upgrader{}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := updater.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	for {
		_, msg, _ := ws.ReadMessage()
		message := parsers.NewMessage(string(msg))
		mm := gomsglog.Persist(message)
		ws.WriteJSON(mm)
	}

}
