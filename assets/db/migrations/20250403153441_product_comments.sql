-- +goose Up
-- +goose StatementBegin
COMMENT ON COLUMN company.company_name_id IS 'The company name ID from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON TABLE company IS 'Companies as listed by the Canadian Food Inspection Agency (CFIA) https://health-products.canada.ca/api/documentation/lnhpd-documentation-en.html';
COMMENT ON COLUMN company.company_name IS 'The company name from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN dosage_form.name IS 'The dosage form from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON TABLE dosage_form IS 'The dosage form from the Canadian Food Inspection Agency (CFIA) https://health-products.canada.ca/api/documentation/lnhpd-documentation-en.html';
COMMENT ON TABLE product_license IS 'The product license from the Canadian Food Inspection Agency (CFIA)https://health-products.canada.ca/api/documentation/lnhpd-documentation-en.html';
COMMENT ON COLUMN product_license.lnhpd_id IS 'The license ID from the Natural Health Products Directorate (NHPD)';
COMMENT ON COLUMN product_license.license_number IS 'The license number from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.license_date IS 'The license date from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.revised_date IS 'The revised date from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.time_receipt IS 'The time receipt from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.date_start IS 'The date the product was first marketed in Canada';
COMMENT ON COLUMN product_license.sub_submission_type_code IS 'The submission type code from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.sub_submission_type_desc IS 'The submission type description from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.flag_primary_name IS 'Flag indicating if the product name is the primary name';
COMMENT ON COLUMN product_license.flag_product_status IS 'Flag indicating if the product is currently marketed in Canada';
COMMENT ON COLUMN product_license.flag_attested_monograph IS 'Flag indicating if the product has an attested monograph';
COMMENT ON COLUMN product_license.product_name_id IS 'The product name ID from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.product_name IS 'The product name from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN product_license.company_id IS 'The company ID from the Canadian Food Inspection Agency (CFIA)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
COMMENT ON COLUMN company.company_name_id IS NULL;
COMMENT ON TABLE company IS NULL;
COMMENT ON COLUMN company.company_name IS NULL;
COMMENT ON COLUMN dosage_form.name IS NULL;
COMMENT ON TABLE dosage_form IS NULL;
COMMENT ON TABLE product_license IS NULL;
COMMENT ON COLUMN product_license.lnhpd_id IS NULL;
COMMENT ON COLUMN product_license.license_number IS NULL;
COMMENT ON COLUMN product_license.license_date IS NULL;
COMMENT ON COLUMN product_license.revised_date IS NULL;
COMMENT ON COLUMN product_license.time_receipt IS NULL;
COMMENT ON COLUMN product_license.date_start IS NULL;
COMMENT ON COLUMN product_license.sub_submission_type_code IS NULL;
COMMENT ON COLUMN product_license.sub_submission_type_desc IS NULL;
COMMENT ON COLUMN product_license.flag_primary_name IS NULL;
COMMENT ON COLUMN product_license.flag_product_status IS NULL;
COMMENT ON COLUMN product_license.flag_attested_monograph IS NULL;
COMMENT ON COLUMN product_license.product_name_id IS NULL;
COMMENT ON COLUMN product_license.product_name IS NULL;
COMMENT ON COLUMN product_license.company_id IS NULL;
-- +goose StatementEnd
