package services

import (
	"fmt"
	"io"
	"net/http"
	"service/internal/domain"
	"service/internal/repository"

	"github.com/rs/zerolog/log"
)

type NodeService struct {
	NodeRepository *repository.NodeRepository
}

func NewNodeService(nodeRepository *repository.NodeRepository) *NodeService {
	return &NodeService{NodeRepository: nodeRepository}
}

func (s *NodeService) GetNodes() ([]domain.Node, error) {
	return s.NodeRepository.GetNodes()
}

func (s *NodeService) GetNetDataByName(nodeName string) (string, error) {
	node, err := s.NodeRepository.GetNode(nodeName)
	if err != nil {
		return "", err
	}

	fmt.Println(node.Name, " - ")
	url := fmt.Sprintf("http://%s:19999/api/v3/data?filter=system.*&points=1&format=json", nodeName)

	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reach Netdata")
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return "", err
	}

	log.Info().
		Int("status", resp.StatusCode).
		RawJSON("payload", body).
		Msg("Netdata API Response")
	return string(body), nil
}
