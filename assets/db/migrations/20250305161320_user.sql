-- +goose Up
-- +goose StatementBegin

-- All timestamps in UTC
SET TIME ZONE 'UTC';

-- Set up a function to automatically update `updated_at` fields on tables
-- See https://x-team.com/blog/automatic-timestamps-with-postgresql/
CREATE OR REPLACE FUNCTION trigger_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
	IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
		NEW.modified = CURRENT_TIMESTAMP; 
		RETURN NEW;
	ELSE
		RETURN OLD;
	END IF;
END;
$$ language 'plpgsql';

-- Create a function which we can use to set-up automatic timestamp updating
-- for any table.
CREATE OR REPLACE FUNCTION autoupdate_timestamp(table_name TEXT)
	RETURNS VOID
	LANGUAGE plpgsql AS
$func$
BEGIN
	EXECUTE format('
		CREATE TRIGGER update_timestamp_%I
		BEFORE UPDATE ON %1$I
		FOR EACH ROW
		EXECUTE PROCEDURE trigger_update_timestamp();
	', table_name);
END
$func$;

-- Members
CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	username TEXT UNIQUE NOT NULL,
    first_name TEXT,
    last_name TEXT,
	gender TEXT,
	date_of_birth DATE,
	country TEXT, -- ISO Code
	city TEXT,

	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
SELECT autoupdate_timestamp('users');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP FUNCTION trigger_update_timestamp;
-- +goose StatementEnd
