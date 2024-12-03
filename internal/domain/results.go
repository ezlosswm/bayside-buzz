package domain

import "fmt"

type Results struct {
	TotalOrganizers string
}

func NewResults(totalOrg int64) *Results {
	return &Results{
		TotalOrganizers: fmt.Sprintf("%v", totalOrg),
	}
}
