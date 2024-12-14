-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_tags (
    eventId INTEGER NOT NULL, 
    tagId INTEGER NOT NULL, 
    PRIMARY KEY (eventId, tagId), 
    FOREIGN KEY (eventId) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (tagId) REFERENCES tags(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_tags;
-- +goose StatementEnd
