-- +goose Up
-- +goose StatementBegin
ALTER TABLE company RENAME TO companies;
ALTER TABLE dosage_form RENAME TO dosage_forms;
ALTER TABLE product_license RENAME TO product_licenses;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE companies RENAME TO company;
ALTER TABLE dosage_forms RENAME TO dosage_form;
ALTER TABLE product_licenses RENAME TO product_license;
-- +goose StatementEnd
