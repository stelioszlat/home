package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service/internal/domain"
	"service/internal/util/response"

	"github.com/gin-gonic/gin"
)

type GetDataRequest struct {
	device string `binding:"required"`
}

func (s *Server) CreateData(c *gin.Context) {
	var data domain.Data
	fmt.Println(c.Request)
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if err := s.data.CreateData(&data); err != nil {
		response.Error(c, http.StatusInternalServerError, "Could not create data", err)
		return
	}
	response.Success(c, "Data created successfully", nil)
}

func (s *Server) GetData(c *gin.Context) {
	data, err := s.data.GetData(c.GetInt("deviceId"))
	if err != nil {
		response.InternalServerError(c, "Failed to fetch data", err)
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		response.InternalServerError(c, "Could not create response", err)
		return
	}

	response.Success(c, "", string(resp))
}
