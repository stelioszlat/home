package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func connect() {
	url := "ws://localhost:8600/ws"

	dialer := websocket.DefaultDialer

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error: ", err)
	}
	defer conn.Close()

	var message string
	for {
		fmt.Printf("Say: ")
		_, err := fmt.Scan(&message)
		if err != nil {
			log.Fatal("Error reading: ", err)
		}
		err = conn.WriteMessage(websocket.TextMessage, bytes.NewBufferString(message).Bytes())
		if err != nil {
			log.Println("Write:", err)
			return
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read: ", err)
			return
		}
		fmt.Printf("Receiveed: %s\n", msg)
	}
}
