package repository

import (
	"service/internal/domain"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type NodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r NodeRepository) GetNodes() ([]domain.Node, error) {
	var nodes []domain.Node
	err := r.db.Model(&domain.Node{}).Where("isActive = ?", true).Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r NodeRepository) GetAllNodes() ([]domain.Node, error) {
	var nodes []domain.Node
	err := r.db.Model(&domain.Node{}).Where("isActive = ?", true).Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r NodeRepository) GetNode(nodeName string) (domain.Node, error) {
	var node domain.Node
	err := r.db.Model(&domain.Node{}).Where("name = ?", nodeName).Find(&node).Error
	if err != nil {
		log.Err(err).Msg("Could not find node")
		return node, err
	}
	return node, nil
}
