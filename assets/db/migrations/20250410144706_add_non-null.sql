-- +goose Up
-- +goose StatementBegin
ALTER TABLE companies
    ALTER COLUMN company_id SET NOT NULL,
    ALTER COLUMN company_name SET NOT NULL,
    ALTER COLUMN company_name_id SET NOT NULL;

ALTER TABLE dosage_forms
    ALTER COLUMN name SET NOT NULL;

ALTER TABLE product_licenses
    ALTER COLUMN time_receipt SET NOT NULL,
    ALTER COLUMN date_start SET NOT NULL,
    ALTER COLUMN flag_primary_name SET NOT NULL,
    ALTER COLUMN flag_product_status SET NOT NULL,
    ALTER COLUMN flag_attested_monograph SET NOT NULL,
    ALTER COLUMN product_name SET NOT NULL,
    ALTER COLUMN product_name_id SET NOT NULL;
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product_licenses
    ALTER COLUMN time_receipt DROP NOT NULL,
    ALTER COLUMN date_start DROP NOT NULL,
    ALTER COLUMN flag_primary_name DROP NOT NULL,
    ALTER COLUMN flag_product_status DROP NOT NULL,
    ALTER COLUMN flag_attested_monograph DROP NOT NULL,
    ALTER COLUMN product_name DROP NOT NULL,
    ALTER COLUMN product_name_id DROP NOT NULL;

ALTER TABLE dosage_forms
    ALTER COLUMN name DROP NOT NULL;

ALTER TABLE companies
    ALTER COLUMN company_id DROP NOT NULL,
    ALTER COLUMN company_name DROP NOT NULL,
    ALTER COLUMN company_name_id DROP NOT NULL;
-- +goose StatementEnd
