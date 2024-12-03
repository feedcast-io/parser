package handlers

import (
	"github.com/feedcast-io/parser/internal"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFromCsv(t *testing.T) {
	files, err := filepath.Glob("../../samples/*.csv")

	if err != nil {
		t.Fatal("Fatal", err)
	}

	for _, fileName := range files {
		r, e := os.Open(fileName)
		if e != nil {
			t.Fatal("Fatal", e)
		}

		p, er := GetFromCsv(r)

		internal.TestProductList(t, filepath.Base(fileName), p, er)
	}
}
