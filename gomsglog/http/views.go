package http

import (
	"encoding/json"
	"fmt"
	"html/template"
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

	arch := false
	if archStr := r.URL.Query().Get("archived"); archStr != "" {
		arch, _ = strconv.ParseBool(archStr)
	}

	messages := gomsglog.LoadMessages(
		limit,
		offset,
		tags,
		users,
		attrs,
		arch,
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
	db := gomsglog.GetDB()
	var attrs []string
	db.Model(&gomsglog.AttributeSet{}).Pluck(`DISTINCT "slug"`, &attrs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&attrs)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

}

type IndexViewModel struct {
	Messages           []gomsglog.MessageModel
	EditMessageID      uint
	EditMessageContent string
	EditMessage        bool
}

func indexHandlerRaw(w http.ResponseWriter, r *http.Request, viewModel *IndexViewModel) {
	funcMap := template.FuncMap{
		"safeHTML": func(input string) template.HTML {
			return template.HTML(input)
		},
		"quote": func(id string) template.HTML {
			if intID, err := strconv.Atoi(id); err != nil {
				return template.HTML("not a number")
			} else {
				m, found := gomsglog.LoadMessage(int(intID))
				if !found {
					return template.HTML("not found")
				} else {
					return template.HTML(m.HTML)
				}
			}
		},
	}
	// var terr error
	// tmpl := template.New("templates/index")
	// tmpl = tmpl.Funcs(funcMap)
	// tmpl, terr = tmpl.ParseFiles("templates/index.html")
	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html"))
	if err := tmpl.Execute(w, viewModel); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}
	// if terr != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "Error loading template: %s", terr.Error())
	// } else {

	// }
}

func SplitQuery(r *http.Request, param string) []string {
	value := r.URL.Query().Get(param)
	var splitted []string
	for _, s := range strings.Split(strings.Trim(value, ""), ",") {
		if s != "" {
			splitted = append(splitted, s)
		}
	}
	return splitted
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	messages := gomsglog.LoadMessages(100, 0, SplitQuery(r, "tags"), SplitQuery(r, "users"), []string{}, false)
	viewModel := &IndexViewModel{messages, 0, "", false}
	editIDstr := r.URL.Query().Get("edit")
	if editIDstr != "" {
		editID, _ := strconv.Atoi(editIDstr)
		if message, found := gomsglog.LoadMessage(editID); found {
			viewModel.EditMessage = true
			viewModel.EditMessageContent = message.Original
			viewModel.EditMessageID = message.ID
		} else {
			println("no message found with id ", editIDstr)
		}

	}

	indexHandlerRaw(w, r, viewModel)

}

func SubmitMessageHandler(w http.ResponseWriter, r *http.Request) {
	newMessage := r.FormValue("message")
	messageIDstr := r.FormValue("message-id")
	if newMessage != "" {
		fmt.Printf("Wert: '%s'", messageIDstr)
		if messageIDstr != "" {
			if messageID, err := strconv.Atoi(messageIDstr); err != nil {
				// todo parse error
				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			} else {
				msg := parsers.NewMessage(newMessage)
				gomsglog.Update(messageID, msg)
			}
		} else {
			msg := parsers.NewMessage(newMessage)
			gomsglog.Persist(msg)
		}
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
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
