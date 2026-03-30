package services

import (
	"service/internal/domain"

	"gorm.io/gorm"
)

type DataService struct {
	db *gorm.DB
}

func NewDataService(db *gorm.DB) *DataService {
	return &DataService{db: db}
}

func (s *DataService) CreateData(data *domain.Data) error {
	return s.db.Model(&domain.Data{}).Create(data).Error
}

func (s *DataService) GetData(device int) ([]domain.Data, error) {
	var data []domain.Data
	if err := s.db.Model(&domain.Data{}).Where("device = ?", device).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
