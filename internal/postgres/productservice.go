package postgres

import (
	"beetle/internal/domain"

	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
	ReadDB            *gorm.DB
	WriteDB           *gorm.DB
	PaginationService domain.IPaginationService
}

// paginateQuery is a helper function that handles common pagination logic
func (s *ProductService) paginateQuery(model interface{}, results *domain.PaginatedResults, offset int) error {
	var total int64
	if err := s.ReadDB.Model(model).Count(&total).Error; err != nil {
		return err
	}
	results.Total = int(total)

	if err := s.ReadDB.Offset(offset).Limit(results.Limit).Find(model).Error; err != nil {
		return err
	}

	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
	length := modelValue.Len()
	interfaceSlice := make([]interface{}, length)
	for i := 0; i < length; i++ {
		interfaceSlice[i] = modelValue.Index(i).Interface()
	}
	results.Data = &interfaceSlice
	return nil
}

func (s *ProductService) ReadLicenseByID(id uuid.UUID) (*domain.ProductLicense, error) {
	var productLicense domain.ProductLicense
	if err := s.ReadDB.Where("id = ?", id).First(&productLicense).Error; err != nil {
		return nil, err
	}
	return &productLicense, nil
}

func (s *ProductService) GetLicenses(pi *domain.PaginationQuery) (*domain.PaginatedResults, error) {
	var licenses []domain.ProductLicense
	results, offset := pi.CreateResults()

	if err := s.PaginationService.Paginate(&licenses, &results, offset); err != nil {
		return nil, err
	}
	return &results, nil
}

func (s *ProductService) ReadDosageFormByID(id uuid.UUID) (*domain.DosageForm, error) {
	var dosageForm domain.DosageForm
	if err := s.ReadDB.Where("id = ?", id).First(&dosageForm).Error; err != nil {
		return nil, err
	}
	return &dosageForm, nil
}

func (s *ProductService) GetDosageForms(pi *domain.PaginationQuery) (*domain.PaginatedResults, error) {
	var dosageForms []domain.DosageForm
	results, offset := pi.CreateResults()

	if err := s.PaginationService.Paginate(&dosageForms, &results, offset); err != nil {
		return nil, err
	}
	return &results, nil
}

func (s *ProductService) ReadSubmissionTypeByID(id uuid.UUID) (*domain.SubmissionType, error) {
	var submissionType domain.SubmissionType
	if err := s.ReadDB.Where("id = ?", id).First(&submissionType).Error; err != nil {
		return nil, err
	}
	return &submissionType, nil
}

func (s *ProductService) GetSubmissionTypes(pi *domain.PaginationQuery) (*domain.PaginatedResults, error) {
	var submissionTypes []domain.SubmissionType
	results, offset := pi.CreateResults()

	if err := s.PaginationService.Paginate(&submissionTypes, &results, offset); err != nil {
		return nil, err
	}
	return &results, nil
}
