-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_license (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lnhpd_id INT NOT NULL UNIQUE,
    license_number INT,
    license_date DATE,
    revised_date DATE,
    time_receipt DATE,
    date_start DATE,
    product_name_id INT,
    product_name TEXT,
    dosage_form TEXT,
    company_id INT,
    company_name TEXT,
    company_name_id INT,
    sub_submission_type_code INT,
    sub_submission_type_desc TEXT,
    flag_primary_name BOOLEAN,
    flag_product_status BOOLEAN,
    flag_attested_monograph BOOLEAN,
    
    
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);
SELECT autoupdate_timestamp('product_license');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_license;
-- +goose StatementEnd
