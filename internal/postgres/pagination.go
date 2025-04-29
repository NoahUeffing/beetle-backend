package postgres

import (
	"beetle/internal/domain"
	"reflect"

	"gorm.io/gorm"
)

// PaginateQuery is a helper function that handles common pagination logic
func PaginateQuery(db *gorm.DB, model interface{}, results *domain.PaginatedResults, offset int) error {
	var total int64
	if err := db.Model(model).Count(&total).Error; err != nil {
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
