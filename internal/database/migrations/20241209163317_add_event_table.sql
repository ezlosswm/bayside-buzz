-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    title TEXT NOT NULL, 
    description TEXT NOT NULL, 
    date DATE NOT NULL, 
    freq TEXT CHECK (freq IN ('once', 'daily', 'weekly', 'monthly')) NOT NULL,
    organizer TEXT NOT NULL, 
    imgPath TEXT NOT NULL, 
    userId INTEGER NOT NULL, 
    FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down 
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
