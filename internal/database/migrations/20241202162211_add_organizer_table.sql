-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizers (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    organizer_name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    img_url TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organizers;
-- +goose StatementEnd
