package repository

import (
	"service/internal/domain"

	"gorm.io/gorm"
)

type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	if db == nil {
		panic("Module Repository initialized with NIL database pointer")
	}
	return &ModuleRepository{db: db}
}

func (m *ModuleRepository) GetModules() ([]domain.Module, error) {
	var modules []domain.Module
	var total int64
	err := m.db.Model(&domain.Module{}).Where("isEnabled = ?", true).Count(&total).Find(&modules).Error

	return modules, err
}

func (m *ModuleRepository) GetAllModules() ([]domain.Module, error) {
	var modules []domain.Module
	var total int64
	err := m.db.Model(&domain.Module{}).Count(&total).Find(&modules).Error

	return modules, err
}
