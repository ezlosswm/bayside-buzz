package components

import (
	"strconv"
	"strings"
)

func ToSlice(s string) []string {
    return strings.Split(s, ",")
}

func NewURL(eventId int32) string {
   id := int(eventId) 
   return strings.Join([]string{"/event/", strconv.Itoa(id)}, "")
}
