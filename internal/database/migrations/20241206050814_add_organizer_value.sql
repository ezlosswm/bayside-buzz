-- +goose Up
-- +goose StatementBegin
ALTER TABLE organizers ADD COLUMN value TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE organizers DROP COLUMN value;
-- +goose StatementEnd
