package services

import (
	"service/internal/domain"
	"service/internal/repository"
)

type ModuleService struct {
	moduleRepository *repository.ModuleRepository
}

func NewModuleService(ModuleRepository *repository.ModuleRepository) *ModuleService {
	return &ModuleService{moduleRepository: ModuleRepository}
}

func (s *ModuleService) GetModules() ([]domain.Module, error) {
	modules, err := s.moduleRepository.GetModules()
	if err != nil {
		return nil, err
	}

	return modules, nil
}
