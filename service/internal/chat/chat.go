package chat

import (
	"fmt"
	"service/config"
	"service/internal/logger"
	"service/internal/services"

	"github.com/rs/zerolog/log"
)

type Chat struct {
	service *services.ChatService
}

func New() *Chat {
	log.Logger = logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err)
	}
	chatService := services.NewChatService(cfg.Chat.ApiKey)

	return &Chat{
		service: chatService,
	}
}

func (c *Chat) Chat(prompt, model string) {
	resp, err := c.service.RunLocalPrompt(prompt, model)
	if err != nil {
		log.Err(err).Msg("Could not load chat response")
	}

	fmt.Println(resp)
}

func (c *Chat) Models() {
	c.service.ListLocalModels()
}
