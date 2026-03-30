package server

import (
	"encoding/json"
	"net/http"
	"service/internal/util/response"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetNodes(c *gin.Context) {
	nodes, err := s.node.GetNodes()
	if err != nil {
		response.InternalServerError(c, "Unable to get nodes", err)
		return
	}

	resp, err := json.Marshal(nodes)
	if err != nil {
		response.InternalServerError(c, "Could not create response", err)
		return
	}

	response.Success(c, "", string(resp))
}

func (s *Server) GetNodeData(c *gin.Context) {
	data, err := s.node.GetNetDataByName("localhost")
	if err != nil {
		response.InternalServerError(c, "Could not load net data", err)
		return
	}
	response.Success(c, "", data)
}

func (s *Server) GetNodeDataByName(c *gin.Context) {
	name := c.Param("name")

	data, err := s.node.GetNetDataByName(name)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Could not find node", err)
		return
	}

	response.Success(c, "", data)
}
