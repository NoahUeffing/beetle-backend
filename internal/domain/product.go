package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	HealthCanadaBaseProductURL = "https://health-products.canada.ca/api/natural-licences/"
	ProductLicenseURL          = HealthCanadaBaseProductURL + "productlicence/?lang=en&type=json"
	ProductURLID               = "&id="
	ProductURLPageNumber       = "&page="
)

type DosageForm struct {
	Entity
	Name string `json:"name_id"`
}

type SubmissionType struct {
	Entity
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type ProductLicense struct {
	Entity
	LNHPDID               int       `json:"lnhpd_id"`
	LicenseNumber         int       `json:"license_number"`
	LicenseDate           time.Time `json:"license_date"`
	RevisedDate           time.Time `json:"revised_date"`
	TimeReceipt           time.Time `json:"time_receipt"`
	DateStart             time.Time `json:"date_start"`
	ProductNameID         int       `json:"product_name_id"`
	ProductName           string    `json:"product_name"`
	DosageFormID          uuid.UUID `json:"dosage_form_id"`
	CompanyID             uuid.UUID `json:"company_id"`
	SubmissionTypeID      uuid.UUID `json:"submission_type_id"`
	FlagPrimaryName       bool      `json:"flag_primary_name"`
	FlagProductStatus     bool      `json:"flag_product_status"`
	FlagAttestedMonograph bool      `json:"flag_attested_monograph"`
}

type ProductNameFilter struct {
	ProductName string `json:"product_name"`
}

type IProductService interface {
	ReadLicenseByID(id uuid.UUID) (*ProductLicense, error)
	GetLicenses(pi *PaginationQuery, filters ...Filter) (*PaginatedResults, error) // TODO: Add filtering and sorting
	GetDosageForms(pi *PaginationQuery) (*PaginatedResults, error)
	ReadDosageFormByID(id uuid.UUID) (*DosageForm, error)
	GetSubmissionTypes(pi *PaginationQuery) (*PaginatedResults, error)
	ReadSubmissionTypeByID(id uuid.UUID) (*SubmissionType, error)
}
