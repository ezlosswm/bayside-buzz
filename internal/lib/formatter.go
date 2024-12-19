package lib

import (
	"database/sql"
	"strings"
)

func OrganizerToValue(organizer string) sql.NullString {
	s := strings.Split(organizer, " ")
	v := strings.TrimSuffix(strings.ToLower(s[0]), "'s")

	return sql.NullString{String: v, Valid: v != ""}
}

