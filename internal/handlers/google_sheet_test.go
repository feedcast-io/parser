package handlers

import (
	"fmt"
	"github.com/feedcast-io/parser/internal"
	"testing"
)

func TestGetFromGoogleSheet(t *testing.T) {
	testUris := []string{
		"https://docs.google.com/spreadsheets/d/144LPOb9LP25E-07IlYnhXYHIO3h3AkDWpjkdF0l4Dlk/edit#gid=0",
		"https://docs.google.com/spreadsheets/d/1Hk6Ih6-WPVQttwUPFlHkRLjuYL9CkXh5GLz6rIEMyuE/edit?usp=sharing",
		"https://docs.google.com/spreadsheets/d/1s7IihIv3_nt5RJ8C5EHSFhUhzGxZXay8eUqNZzqrdcw/edit#gid=0",
	}

	for index, u := range testUris {
		p, err := GetFromGoogleSheet(u)

		internal.TestProductList(t, fmt.Sprintf("google-sheet-%d.csv", index), p, err)
	}
}
