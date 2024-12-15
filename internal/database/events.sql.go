// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: events.sql

package database

import (
	"context"
	"time"
)

const countEvents = `-- name: CountEvents :one
SELECT COUNT(*) FROM events
`

func (q *Queries) CountEvents(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countEvents)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (
    title, description, date, freq, organizer, imgPath, userId
    ) VALUES (
    ?, ?, ?, ?, ?, ?, ?
) RETURNING id, title, description, date, freq, organizer, imgpath, userid
`

type CreateEventParams struct {
	Title       string
	Description string
	Date        time.Time
	Freq        string
	Organizer   string
	Imgpath     string
	Userid      int64
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, createEvent,
		arg.Title,
		arg.Description,
		arg.Date,
		arg.Freq,
		arg.Organizer,
		arg.Imgpath,
		arg.Userid,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Date,
		&i.Freq,
		&i.Organizer,
		&i.Imgpath,
		&i.Userid,
	)
	return i, err
}

const getEvent = `-- name: GetEvent :one
SELECT id, title, description, date, freq, organizer, imgpath, userid FROM events WHERE id = ?
`

func (q *Queries) GetEvent(ctx context.Context, id int64) (Event, error) {
	row := q.db.QueryRowContext(ctx, getEvent, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Date,
		&i.Freq,
		&i.Organizer,
		&i.Imgpath,
		&i.Userid,
	)
	return i, err
}

const getEventWithTags = `-- name: GetEventWithTags :one
SELECT 
    e.id AS eventId,
    e.title,
    e.description,
    e.date,
    e.freq,
    e.organizer,
    e.imgPath,
    e.userId,
    COALESCE(GROUP_CONCAT(t.name, ','), '') AS tags
FROM 
    events e
LEFT JOIN 
    event_tags et ON e.id = et.eventId
LEFT JOIN 
    tags t ON et.tagId = t.id
WHERE 
    e.id = ?
GROUP BY 
    e.id, e.title, e.description, e.date, e.freq, e.organizer, e.imgPath, e.userId
`

type GetEventWithTagsRow struct {
	Eventid     int64
	Title       string
	Description string
	Date        time.Time
	Freq        string
	Organizer   string
	Imgpath     string
	Userid      int64
	Tags        interface{}
}

func (q *Queries) GetEventWithTags(ctx context.Context, id int64) (GetEventWithTagsRow, error) {
	row := q.db.QueryRowContext(ctx, getEventWithTags, id)
	var i GetEventWithTagsRow
	err := row.Scan(
		&i.Eventid,
		&i.Title,
		&i.Description,
		&i.Date,
		&i.Freq,
		&i.Organizer,
		&i.Imgpath,
		&i.Userid,
		&i.Tags,
	)
	return i, err
}

const getEvents = `-- name: GetEvents :many
SELECT id, title, description, date, freq, organizer, imgpath, userid FROM events
`

func (q *Queries) GetEvents(ctx context.Context) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Date,
			&i.Freq,
			&i.Organizer,
			&i.Imgpath,
			&i.Userid,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventsByOrganizer = `-- name: GetEventsByOrganizer :many
SELECT id, title, description, date, freq, organizer, imgpath, userid FROM events WHERE organizer = ?
`

func (q *Queries) GetEventsByOrganizer(ctx context.Context, organizer string) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getEventsByOrganizer, organizer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Date,
			&i.Freq,
			&i.Organizer,
			&i.Imgpath,
			&i.Userid,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventsWithTags = `-- name: GetEventsWithTags :many
SELECT
    e.id AS eventId,
    e.title,
    e.description,
    e.date,
    e.freq,
    e.organizer,
    e.imgPath,
    e.userId,
    COALESCE(GROUP_CONCAT(t.name, ','), '') AS tags
FROM
    events e
LEFT JOIN
    event_tags et ON e.id = et.eventId
LEFT JOIN
    tags t ON et.tagId = t.id
GROUP BY
    e.id
`

type GetEventsWithTagsRow struct {
	Eventid     int64
	Title       string
	Description string
	Date        time.Time
	Freq        string
	Organizer   string
	Imgpath     string
	Userid      int64
	Tags        interface{}
}

func (q *Queries) GetEventsWithTags(ctx context.Context) ([]GetEventsWithTagsRow, error) {
	rows, err := q.db.QueryContext(ctx, getEventsWithTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEventsWithTagsRow
	for rows.Next() {
		var i GetEventsWithTagsRow
		if err := rows.Scan(
			&i.Eventid,
			&i.Title,
			&i.Description,
			&i.Date,
			&i.Freq,
			&i.Organizer,
			&i.Imgpath,
			&i.Userid,
			&i.Tags,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
