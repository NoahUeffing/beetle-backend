package domain

import (
	"gorm.io/gorm"
)

const (
	MaxFilterIDs = 5
)

var fieldMap = map[string]string{
	"product search": "product_name",
	"company search": "company_name",
	"company":        "company_id",
	"form":           "dosage_form_id",
}

type GormFilter interface {
	Apply(db *gorm.DB) *gorm.DB
}

type Filter struct {
	Field    string
	Operator string // e.g., "eq", "in", "like" TODO: make and enum
	Value    interface{}
}

func (f Filter) Apply(db *gorm.DB) *gorm.DB {
	dbField, ok := fieldMap[f.Field]
	if !ok {
		// TODO: Error handling
		return db
	}
	switch f.Operator {
	case "in":
		return db.Where(dbField+" IN ?", f.Value)
	case "like":
		if name, ok := f.Value.(string); ok {
			return db.Where(dbField+" LIKE ?", "%"+name+"%")
		}
	default:
		// TODO: add error handling
		return db
	}
	return db
}
