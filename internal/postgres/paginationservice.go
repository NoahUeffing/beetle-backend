package postgres

import (
	"beetle/internal/domain"
	"reflect"

	"gorm.io/gorm"
)

// PaginationService handles pagination operations for database queries
type PaginationService struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

// Paginate handles common pagination logic for database queries
func (s *PaginationService) Paginate(model any, results *domain.PaginatedResults, offset int, filters ...domain.Filter) error {
	db := s.ReadDB.Model(model)
	for _, filter := range filters {
		db = filter.Apply(db)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return err
	}
	results.Total = int(total)

	if err := db.Offset(offset).Limit(results.Limit).Find(model).Error; err != nil {
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
