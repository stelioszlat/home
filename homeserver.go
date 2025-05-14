package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type BorisRequest struct {
	Prompt string `json:"prompt"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Model              string  `json:"model"`
		CreatedAt          string  `json:"created_at"`
		Message            Message `json:"message"`
		Done               bool    `json:"done"`
		DoneReason         string  `json:"done_reason,omitempty"`
		TotalDuration      int     `json:"total_duration,omitempty"`
		LoadDuration       int     `json:"load_duration,omitempty"`
		PromptEvalCount    int     `json:"prompt_eval_count,omitempty"`
		PromptEvalDuration int     `json:"prompt_eval_duation,omitempty"`
		EvalCount          int     `json:"eval_count,omitempty"`
		EvalDuration       int     `json:"eval_duration,omitempty"`
	} `json:"choices"`
}

func heartBeatHandler(w http.ResponseWriter, r *http.Request) {
	msg := MessageResponse{Message: "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func borisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintln(w, "This is a GET request")
	case http.MethodGet:
		requestBody := ChatRequest{
			Model: "llama3.2",
			Messages: []Message{
				{
					Role:    "user",
					Content: "Hello! What can you do?",
				},
			},
		}

		data, err := json.Marshal(requestBody)
		if err != nil {
			panic(err)
		}

		resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(body))

		var chatResp ChatResponse
		err = json.Unmarshal(body, &chatResp)
		if err != nil {
			panic(err)
		}

		if len(chatResp.Choices) > 0 {
			fmt.Println("Assistant:", chatResp.Choices[0].Message.Content)
		} else {
			fmt.Println("No response from Ollama")
		}
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	token := os.Getenv("TOKEN")

	haURL := "http://localhost:8123/api/services/light/turn_on"

	body := map[string]string{
		"entity_id": "light.living_room",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", haURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
var mutex = &sync.Mutex{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading: ", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		fmt.Println(string(message))
		broadcast <- string(message)
	}
}

func handleMessages() {
	for {
		message := <-broadcast

		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, bytes.NewBufferString(message).Bytes())
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func server(port int) {
	fmt.Println("Welcome home!")

	http.HandleFunc("/hello", heartBeatHandler)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/boris", borisHandler)
	http.HandleFunc("/ws", wsHandler)
	go handleMessages()

	var address string = fmt.Sprint(":", port)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
