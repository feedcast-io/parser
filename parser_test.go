package parser

import (
	"github.com/feedcast-io/parser/internal"
	"github.com/feedcast-io/parser/resources"
	"os"
	"path/filepath"
	"testing"
)

func TestGetProductsFromUrl(t *testing.T) {
	p, e := GetProducts(resources.Config{
		Url: "https://www.champion-direct.com/gmerchantcentera1023cc2badef4fe64943c9fd1257a7c.fr.shop1.xml",
	})

	internal.TestProductList(t, "champion.xml", p, e)
}

func TestGetProductFromLocalFile(t *testing.T) {
	files, _ := filepath.Glob("samples/*")
	if 0 == len(files) {
		t.Error("FATAL: no files found")
	}

	for _, fileName := range files {
		r, err := os.Open(fileName)
		if err != nil {
			t.Error(err)
			continue
		}

		p, e := GetProductFromLocalFile(r)

		internal.TestProductList(t, fileName, p, e)
	}
}
