package domain

type Company struct {
	Entity
	CompanyID     int    `json:"company_id"`
	CompanyName   string `json:"company_name"`
	CompanyNameID int    `json:"company_name_id"`
}
