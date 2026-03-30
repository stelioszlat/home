package server

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"service/config"
	"service/internal/database"
	"service/internal/logger"
	"service/lights"

	"service/internal/repository"
	"service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const broker = "localhost:9092"

type APILogEvent struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Path      string `json:"path,omitempty"`
	Status    int    `json:"status,omitempty"`
	LatencyMs int64  `json:"latency_ms,omitempty"`
}

type Server struct {
	config  *config.Config
	db      *gorm.DB
	logger  zerolog.Logger
	modules *services.ModuleService
	node    *services.NodeService
	chat    *services.ChatService
	data    *services.DataService
	// authService    *services.AuthService
	// productService *services.ProductService
	// userService    *services.UserService
	// uploadService  *services.UploadService
}

func New() *Server {
	log.Logger = logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msg(fmt.Sprintf("Server will run on %s:%s", cfg.Server.Host, cfg.Server.Port))

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Database initialization failed")
	}
	if db == nil {
		log.Fatal().Msg("Database connection is nil")
	}

	moduleRepository := repository.NewModuleRepository(db)
	moduleService := services.NewModuleService(moduleRepository)

	nodeRepository := repository.NewNodeRepository(db)
	nodeService := services.NewNodeService(nodeRepository)

	chatService := services.NewChatService(cfg.Chat.ApiKey)

	dataService := services.NewDataService(db)

	return &Server{
		config:  cfg,
		db:      db,
		logger:  log.Logger,
		modules: moduleService,
		node:    nodeService,
		chat:    chatService,
		data:    dataService,
		// authService:    authService,
		// productService: productService,
		// userService:    userService,
		// uploadService:  uploadService,
	}
}

type MessageResponse struct {
	Message string `json:"message"`
}

type BorisRequest struct {
	Message string `json:"message"`
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
}

//	func sendNotification(w Socket, message string) {
//		notification := MessageResponse{Message: message}
//		w.WriteJSON(notification)
//	}
func heartBeatHandler(w http.ResponseWriter, r *http.Request) {
	msg := MessageResponse{Message: "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

}

func sendToDiscord(message string) {
	webhookURL := "https://discordapp.com/api/webhooks/1440459383757602928/kUzZRyZw0XdxB_C67MkW9NgnpINwsKWqUJfRxUbJ_bcpmhdjeLpnrsFSz7K0Pnw23-uR"

	payload := map[string]string{"content": message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending message to Discord:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Unexpected response from Discord:", resp.Status)
	}
}

func borisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintln(w, "Hello from Boris")
	case http.MethodGet:
		requestBodyMessage := r.Header.Get("Message")
		fmt.Println("Request Body Message: ", requestBodyMessage)

		if requestBodyMessage == "" {
			fmt.Fprintln(w, "Please provide a message in the 'Message' header.")
			return
		}

		requestBody := ChatRequest{
			Model: "deepseek-r1:1.5b",
			Messages: []Message{
				{
					Role:    "user",
					Content: requestBodyMessage,
				},
			},
		}

		data, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Fprintln(w, "Sorry, I could not understand your request.")
			return

		}

		resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Fprintln(w, "The chat service is currently unavailable.")
			return
		}
		defer resp.Body.Close()

		var out strings.Builder
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Bytes()
			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}

			var chunk ChatResponse
			err := json.Unmarshal(line, &chunk)
			if err != nil {
				fmt.Fprintln(w, "I could not process the response from the model.")
			}

			if strings.HasPrefix(chunk.Message.Content, "<") {
				continue
			}

			out.WriteString(chunk.Message.Content)

			fmt.Print(chunk.Message.Content)

			if chunk.Done {
				break
			}
		}

		fmt.Fprintln(w, out.String())

		go sendToDiscord(out.String())
	}
}

func lightsHandler(w http.ResponseWriter, r *http.Request) {
	lights.ToggleLights()
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
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
		response := AskChat(string(message))
		if response == "" {
			response = "Sorry, I could not process your request."
		}
		broadcast <- string(response)
	}
}

func AskChat(message string) string {
	token := "s2_d44d659a88b6445e8726e4960c0213ac"
	abacusUrl := ""

	body := map[string]string{
		"entity_id": "light.living_room",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", abacusUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to Abacus AI:", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	return message
}

func GetAbacusModels() []string {
	token := "s2_d44d659a88b6445e8726e4960c0213ac"
	abacusUrl := "https://api.abacus.ai/api/v0/listRouteLLMModels"

	req, err := http.NewRequest("GET", abacusUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	var models []string
	json.NewDecoder(resp.Body).Decode(&models)

	return models
}

func GetAbacusPipelines() []string {
	token := "s2_d44d659a88b6445e8726e4960c0213ac"
	abacusUrl := "https://api.abacus.ai/api/v0/listPipelines"

	req, err := http.NewRequest("GET", abacusUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	var pipelines []string
	json.NewDecoder(resp.Body).Decode(&pipelines)

	return pipelines
}

func GetAbacusSessions() []string {
	token := "s2_d44d659a88b6445e8726e4960c0213ac"
	abacusUrl := "https://api.abacus.ai/api/v0/listChatSessions"

	req, err := http.NewRequest("GET", abacusUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	var sessions []string
	json.NewDecoder(resp.Body).Decode(&sessions)

	return sessions
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

func (s *Server) Serve(port int) {
	// producer, err := kafka.NewProducer(brokers)
	// if err != nil {
	// 	logger.Log.Fatal("failed to initialize producer", zap.Error(err))
	// }
	// defer producer.Close()

	// consumer, err := kafka.NewConsumer(brokers, "boris-log-consumer", []string{kafka.TopicLogs})
	// if err != nil {
	// 	logger.Log.Fatal("failed to initialize consumer", zap.Error(err))
	// }
	// defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Trap OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Info().Msg(sig.String())
	}()

	// consumer.Start(ctx, func(topic, key string, payload []byte) error {
	// 	event, err := kafka.UnmarshalPayload[APILogEvent](payload)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	logger.Log.Info("api log consumed",
	// 		zap.String("service", event.Service),
	// 		zap.String("path", event.Path),
	// 		zap.Int("status", event.Status),
	// 		zap.Int64("latency_ms", event.LatencyMs),
	// 		zap.String("request_id", event.RequestID),
	// 	)
	// 	return nil
	// })

	// InitNotificationHandlers()

	router := s.setupRouter()

	go handleMessages()

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   s.config.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Wrap router with CORS handler
	handler := corsHandler.Handler(router)

	// Create HTTP server
	var address string = fmt.Sprint(":", port)
	server := &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server is running on http://%s", address)
		log.Printf("WebSocket endpoint: ws://%s/ws", address)
		log.Printf("API endpoints available at http://%s/api/", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msg(fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Msg(fmt.Sprintf("Forced shutdown: %v", err))
	}

	log.Print("Server exited cleanly")
}

func (s *Server) setupRouter() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Health check endpoint
	router.GET("/health", s.healthCheck)

	api := router.Group("/api")
	api.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello from Home")
	})
	api.GET("/boris", func(c *gin.Context) {
		// boris
	})
	lights := api.Group("lights")
	lights.GET("/lights", func(c *gin.Context) {
		// lights
	})
	chat := api.Group("/chat")
	chat.GET("/boris", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello from Home")
	})

	modules := api.Group("modules")
	modules.GET("/", s.GetModules)

	node := api.Group("node")
	node.GET("/data", s.GetNodeData)
	nodes := api.Group("nodes")
	nodes.GET("/", s.GetNodes)
	nodes.GET("/:name", s.GetNodeDataByName)

	data := api.Group("data")
	data.POST("/", s.CreateData)
	data.GET("/:deviceId")
	// api.HandleFunc("/modules/{id}", handlers.UpdateModule).Methods("PUT")

	// api.HandleFunc("/notifications/register", handlers.RegisterDeviceToken).Methods("POST")
	// api.HandleFunc("/notifications/send", handlers.SendNotification).Methods("POST")

	// router.HandleFunc("/ws", wsHandler)

	// router.Use(loggingMiddleware)

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
