package domain

import "time"

const (
	HealthCanadaBaseProductURL = "https://health-products.canada.ca/api/natural-licences/"
	ProductLicenseURL          = HealthCanadaBaseProductURL + "productlicence/?lang=en&type=json"
	ProductURLID               = "&id="
	ProductURLPageNumber       = "&page="
)

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
	DosageForm            string    `json:"dosage_form"`
	CompanyID             int       `json:"company_id"`
	CompanyName           string    `json:"company_name"`
	CompanyNameID         int       `json:"company_name_id"`
	SubSubmissionTypeCode int       `json:"sub_submission_type_code"`
	SubSubmissionTypeDesc string    `json:"sub_submission_type_desc"`
	FlagPrimaryName       bool      `json:"flag_primary_name"`
	FlagProductStatus     bool      `json:"flag_product_status"`
	FlagAttestedMonograph bool      `json:"flag_attested_monograph"`
}
