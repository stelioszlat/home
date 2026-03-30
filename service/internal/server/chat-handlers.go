package server

import (
	"service/internal/util/response"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetChatHealth(c *gin.Context) {
	response.Success(c, "Successfull chat request", nil)
}
