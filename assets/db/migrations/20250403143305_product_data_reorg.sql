-- +goose Up
-- +goose StatementBegin
CREATE TABLE company (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id INT,
    company_name TEXT,
    company_name_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
SELECT autoupdate_timestamp('company');

CREATE TABLE dosage_form (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
SELECT autoupdate_timestamp('dosage_form');

ALTER TABLE product_license
ADD COLUMN new_company_id uuid REFERENCES company(id),
ADD COLUMN new_dosage_form_id uuid REFERENCES dosage_form(id);

-- Insert unique companies from product_license into company table
INSERT INTO company (company_id, company_name, company_name_id)
SELECT DISTINCT company_id, company_name, company_name_id
FROM product_license
WHERE company_id IS NOT NULL;

-- Insert unique dosage forms from product_license into dosage_form table
INSERT INTO dosage_form (name)
SELECT DISTINCT dosage_form
FROM product_license
WHERE dosage_form IS NOT NULL;

-- Update product_license to reference the new company records
UPDATE product_license pl
SET new_company_id = c.id
FROM company c
WHERE pl.company_id = c.company_id;

-- Update product_license to reference the new dosage form records
UPDATE product_license pl
SET new_dosage_form_id = df.id
FROM dosage_form df
WHERE pl.dosage_form = df.name;

-- Make the columns NOT NULL after we've populated them
ALTER TABLE product_license
ALTER COLUMN new_company_id SET NOT NULL;

ALTER TABLE product_license
ALTER COLUMN new_dosage_form_id SET NOT NULL;

ALTER TABLE product_license 
DROP COLUMN company_id,
DROP COLUMN company_name,
DROP COLUMN company_name_id,
DROP COLUMN dosage_form;

ALTER TABLE product_license
RENAME COLUMN new_company_id TO company_id;

ALTER TABLE product_license
RENAME COLUMN new_dosage_form_id TO dosage_form_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product_license
RENAME COLUMN company_id TO new_company_id;

ALTER TABLE product_license
RENAME COLUMN dosage_form_id TO new_dosage_form_id;

ALTER TABLE product_license
ADD COLUMN company_id INT,
ADD COLUMN company_name TEXT,
ADD COLUMN company_name_id INT,
ADD COLUMN dosage_form TEXT;

-- Restore company data from company table
UPDATE product_license pl
SET 
    company_id = c.company_id,
    company_name = c.company_name,
    company_name_id = c.company_name_id
FROM company c
WHERE pl.new_company_id = c.id;

-- Restore dosage form data
UPDATE product_license pl
SET dosage_form = df.name
FROM dosage_form df
WHERE pl.new_dosage_form_id = df.id;

ALTER TABLE product_license
DROP COLUMN new_company_id,
DROP COLUMN new_dosage_form_id;

DROP TABLE dosage_form;
DROP TABLE company;
-- +goose StatementEnd
