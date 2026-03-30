package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type ChatService struct {
	apiKey string
}

type OllamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"` // This is the actual command
	Done      bool      `json:"done"`
}

type OllamaModelsResponse struct {
	Models []struct {
		Name  string `json:"name"`
		Model string `json:"model"`
	} `json:"models"`
}

func NewChatService(apiKey string) *ChatService {
	return &ChatService{apiKey: apiKey}
}

func (s *ChatService) GetGeminiCommand(prompt string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	models := s.ListAvailableModels()
	if len(models) == 0 {
		log.Info().Msg("There are no models available")
		return "", nil
	}

	model := client.GenerativeModel(models[0])

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text("You are a CLI proxy. Return ONLY the terminal command for Ubuntu Linux based on the user request. No markdown, no explanations.")},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}

func (s *ChatService) ListAvailableModels() []string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.apiKey))
	if err != nil {
		log.Err(err)
	}
	defer client.Close()

	req := client.ListModels(ctx)
	info, err := req.Next()
	if err != nil {
		log.Err(err)
	}

	var models []string
	models = append(models, strings.Split(info.Name, "/")[1])
	fmt.Println(models)

	fmt.Printf("========GEMINI=============")
	fmt.Printf("Model Name: %s\n", info.Name)
	fmt.Printf("Display Name: %s\n", info.DisplayName)
	fmt.Printf("Description: %s\n", info.Description)
	fmt.Printf("Input Token Limit: %d\n", info.InputTokenLimit)
	fmt.Printf("Output Token Limit: %d\n", info.OutputTokenLimit)
	fmt.Printf("Supported Actions: %v\n", info.SupportedGenerationMethods)
	fmt.Println("-------------------------------------------")

	return models
}

func (s *ChatService) ListLocalModels() (map[string]string, error) {
	url := "http://localhost:11434/api/tags"

	resp, err := http.Get(url)
	if err != nil {
		log.Err(err).Msg("Failed to connect to ollama: %w")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Err(err).Msg("Ollama returned an error")
		return nil, err
	}

	var ollamaResp OllamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		log.Err(err).Msg("Failed to decode response:")
	}

	if len(ollamaResp.Models) == 0 {
		log.Info().Msg("There no available ollama models")
		return nil, err
	}

	var models = make(map[string]string)
	for _, model := range ollamaResp.Models {
		fmt.Printf("- %s\n", model.Name)
		models[model.Name] = model.Model
	}

	return models, nil
}

func (s *ChatService) RunLocalPrompt(prompt string, model string) (string, error) {
	models, err := s.ListLocalModels()
	if err != nil {
		return "", err
	}

	url := "http://localhost:11434/api/generate"
	payload := map[string]interface{}{
		"model":  models[model],
		"prompt": prompt,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "applicaton/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Failed to connect to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama returned an error: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}

func (s *ChatService) ollamaRequest(url string, payload map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "applicaton/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Failed to connect to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama returned an error: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}
