package postgres

import (
	"beetle/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyService struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func (s *CompanyService) ReadByID(id uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	if err := s.ReadDB.Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (s *CompanyService) GetCompanies(pi *domain.PaginationQuery) (*domain.PaginatedResults, error) {
	var companies []domain.Company
	results, offset := pi.CreateResults()

	if err := PaginateQuery(s.ReadDB, &companies, &results, offset); err != nil {
		return nil, err
	}
	return &results, nil
}
