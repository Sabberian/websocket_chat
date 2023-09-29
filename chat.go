package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	clients[conn] = true
	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			delete(clients, conn)
			return
		}
		broadcast <- message
	}
}

func handleMessage() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error :%v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.HandleFunc("/ws", handleConnection)
	go handleMessage()

	port := ":8080"
	fmt.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
