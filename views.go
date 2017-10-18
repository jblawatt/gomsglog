package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type InputMessage struct {
	Message string `json:"message"`
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	input := new(InputMessage)

	json.NewDecoder(r.Body).Decode(&input)
	msg := NewMessage(input.Message)
	mm := Persist(msg)
	json.NewEncoder(w).Encode(&mm)

}

func GetMessages(w http.ResponseWriter, r *http.Request) {

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
	messages := LoadMessages(
		limit,
		offset,
		tags,
		users,
	)
	json.NewEncoder(w).Encode(&messages)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	messageID, _ := strconv.Atoi(vars["messageID"])
	message, found := LoadMessage(messageID)
	if found {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&message)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"error":  http.StatusText(http.StatusNotFound),
			"status": http.StatusNotFound,
		})
	}
}
