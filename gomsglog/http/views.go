package http

import (
	"encoding/json"
	"fmt"
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
	messages := gomsglog.LoadMessages(
		limit,
		offset,
		tags,
		users,
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
		<head>
			<title>GoMSG-Log</title>
			
		</head>
		<body>
			<div id="app">
				<ul>
					<li v-for="msg in messages">
						{{msg.ID}} <br>
						<div v-html="msg.HTML"></div>
					</li>
				</ul>
				<form method="POST" action="/api/messages">
					<input placeholder="tell me more">
					</input>
				</form>
			</div>
			<script src="https://unpkg.com/vue"></script>
			<script src="/static/js/gomsglog.js"></script>
		</body>
	</html>`)
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
