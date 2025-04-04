package postgres

import (
	"beetle/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

func (s *ProductService) ReadLicenseByID(id uuid.UUID) (*domain.ProductLicense, error) {
	var productLicense domain.ProductLicense
	if err := s.ReadDB.Where("id = ?", id).First(&productLicense).Error; err != nil {
		return nil, err
	}
	return &productLicense, nil
}
