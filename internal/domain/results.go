package domain

import "fmt"

type Results struct {
	TotalOrganizers string
	TotalEvents string
}

func NewResults(totalOrg, totalEvents int64) *Results {
	return &Results{
		TotalOrganizers: fmt.Sprintf("%v", totalOrg),
		TotalEvents: fmt.Sprintf("%v", totalEvents),
	}
}
