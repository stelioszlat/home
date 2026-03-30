package repository

import (
	"service/internal/domain"

	"gorm.io/gorm"
)

type DataRepository struct {
	DB *gorm.DB
}

func NewDataRepository(db *gorm.DB) *DataRepository {
	if db == nil {
		panic("Data Repository initialized with NIL database pointer")
	}
	return &DataRepository{DB: db}
}

func (r *DataRepository) GetData(device string) ([]domain.Data, error) {
	var data []domain.Data
	err := r.DB.Model(&domain.Data{}).Where("device = ?", device).Find(&data).Error

	return data, err
}
