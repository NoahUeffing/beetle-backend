package domain

import "github.com/google/uuid"

type Company struct {
	Entity
	CompanyID     int    `json:"company_id"`
	CompanyName   string `json:"company_name"`
	CompanyNameID int    `json:"company_name_id"`
}

type CompanyFilter struct {
	CompanyName string `json:"company_name"`
}

type ICompanyService interface {
	ReadByID(id uuid.UUID) (*Company, error)
	GetCompanies(pi *PaginationQuery) (*PaginatedResults, error) // TODO: Add filtering and sorting
}
