package main

import (
	_ "embed"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

//go:embed index.html
var indexHTML string

type CommandRequest struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

var upgrader = websocket.Upgrader{}

type WSConnectionWrapper struct {
	Conn  *websocket.Conn
	Mutex *sync.Mutex
}

func NewWSConnectionWrapper(w http.ResponseWriter, r *http.Request) (*WSConnectionWrapper, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &WSConnectionWrapper{
		Conn:  conn,
		Mutex: &sync.Mutex{},
	}, nil
}

func (sc *WSConnectionWrapper) close() {
	sc.Mutex.Unlock()
	err := sc.Conn.Close()
	if err != nil {
		log.Println("Unable to close websocket")
	}
}

func (sc *WSConnectionWrapper) readRequest() (CommandRequest, error) {
	_, responseBody, err := sc.Conn.ReadMessage()
	if err != nil {
		return CommandRequest{}, err
	}
	var commandRequest CommandRequest
	if err = json.Unmarshal(responseBody, &commandRequest); err != nil {
		return CommandRequest{}, err
	}
	return commandRequest, nil
}

func (sc *WSConnectionWrapper) writeResponse(cmdID string, status string, output string) {
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()
	err := sc.Conn.WriteJSON(map[string]string{"id": cmdID, "status": status, "output": output})
	if err != nil {
		log.Println("Unable to write out json to websocket", cmdID)
	}
}

func (params *Config) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	connWrapped, err := NewWSConnectionWrapper(w, r)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer connWrapped.close()
	for {
		commandRequest, readError := connWrapped.readRequest()
		if readError != nil {
			log.Println("Error reading JSON payload:", readError)
			break
		}
		go handleCommand(&commandRequest, params, connWrapped)
	}
}

func handleCommand(command *CommandRequest, params *Config, conn *WSConnectionWrapper) {
	conn.writeResponse(command.ID, "Running", "")
	if err := ExecCmd(command.Command, params, func(line string) {
		log.Println("ytdlp:", line)
		conn.writeResponse(command.ID, "Running", line)
		if strings.HasPrefix(line, "DoneFile$") {
			downloadedFilePath := strings.TrimPrefix(line, "DoneFile$")
			conn.writeResponse(command.ID, "Meta", "Updating metadata from Deezer...")
			metaResult := UpdateFileMetadata(downloadedFilePath)
			conn.writeResponse(command.ID, "Meta", metaResult)
		}
	}); err != nil {
		conn.writeResponse(command.ID, "Error", err.Error())
	} else {
		conn.writeResponse(command.ID, "Done", "Task Done!")
	}
}

func ServeHTML(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(indexHTML)); err != nil {
		log.Println("Error serving index page:", err)
	}
}
