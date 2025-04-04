-- +goose Up
-- +goose StatementBegin
CREATE TABLE submission_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code INT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
SELECT autoupdate_timestamp('submission_types');

COMMENT ON TABLE submission_types IS 'Submission types from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN submission_types.code IS 'The submission type code from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN submission_types.description IS 'The submission type description from the Canadian Food Inspection Agency (CFIA)';

ALTER TABLE product_license
ADD COLUMN new_submission_type_id UUID REFERENCES submission_types(id);

INSERT INTO submission_types (code, description)
SELECT DISTINCT sub_submission_type_code, sub_submission_type_desc
FROM product_license
WHERE sub_submission_type_code IS NOT NULL;

UPDATE product_license pl
SET new_submission_type_id = st.id
FROM submission_types st
WHERE pl.sub_submission_type_code = st.code;

ALTER TABLE product_license
ALTER COLUMN new_submission_type_id SET NOT NULL;

ALTER TABLE product_license 
DROP COLUMN sub_submission_type_code,
DROP COLUMN sub_submission_type_desc;

ALTER TABLE product_license
RENAME COLUMN new_submission_type_id TO submission_type_id;

COMMENT ON COLUMN product_license.submission_type_id IS 'Reference to the submission type for this product license';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE product_license
RENAME COLUMN submission_type_id TO new_submission_type_id;

ALTER TABLE product_license
ADD COLUMN sub_submission_type_code INT,
ADD COLUMN sub_submission_type_desc TEXT;

UPDATE product_license pl
SET 
    sub_submission_type_code = st.code,
    sub_submission_type_desc = st.description
FROM submission_types st
WHERE pl.new_submission_type_id = st.id;

ALTER TABLE product_license
DROP COLUMN new_submission_type_id;

DROP TABLE submission_types;
-- +goose StatementEnd
