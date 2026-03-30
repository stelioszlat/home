package server

import (
	"net/http"
	"service/internal/util/response"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetModules(c *gin.Context) {
	modules, err := s.modules.GetModules()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Could not fetch modules", err)
		return
	}
	response.Success(c, "", modules)
}
