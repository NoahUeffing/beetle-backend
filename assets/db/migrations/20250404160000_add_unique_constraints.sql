-- +goose Up
-- +goose StatementBegin

-- Add unique constraints for companies table
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'companies_company_id_key'
    ) THEN
        ALTER TABLE companies
        ADD CONSTRAINT companies_company_id_key UNIQUE (company_id);
    END IF;
END $$;

-- Add unique constraint for dosage_forms table
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'dosage_forms_name_key'
    ) THEN
        ALTER TABLE dosage_forms
        ADD CONSTRAINT dosage_forms_name_key UNIQUE (name);
    END IF;
END $$;

-- Add unique constraint for submission_types table
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'submission_types_code_key'
    ) THEN
        ALTER TABLE submission_types
        ADD CONSTRAINT submission_types_code_key UNIQUE (code);
    END IF;
END $$;

-- Add unique constraint for product_licenses table
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'product_licenses_lnhpd_id_key'
    ) THEN
        ALTER TABLE product_licenses
        ADD CONSTRAINT product_licenses_lnhpd_id_key UNIQUE (lnhpd_id);
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_company_id_key;
ALTER TABLE dosage_forms DROP CONSTRAINT IF EXISTS dosage_forms_name_key;
ALTER TABLE submission_types DROP CONSTRAINT IF EXISTS submission_types_code_key;
ALTER TABLE product_licenses DROP CONSTRAINT IF EXISTS product_licenses_lnhpd_id_key;
-- +goose StatementEnd 