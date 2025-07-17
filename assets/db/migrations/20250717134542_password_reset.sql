-- +goose Up
-- +goose StatementBegin
CREATE TABLE password_reset_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code INT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    confirmed bool DEFAULT false,
    expiry TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
SELECT autoupdate_timestamp('password_reset_codes');

COMMENT ON TABLE password_reset_codes IS 'Submission types from the Canadian Food Inspection Agency (CFIA)';
COMMENT ON COLUMN password_reset_codes.confirmed IS 'If the code was used to reset email';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE password_reset_codes;
-- +goose StatementEnd
