package internal

import (
	"github.com/feedcast-io/parser/resources"
	"path/filepath"
	"strings"
	"testing"
)

func GetParserResult(products chan []resources.Product, errors chan error) ([]resources.Product, error) {
	result := make([]resources.Product, 0)
	var e error

	for ok1, ok2 := true, true; ok1 || ok2; {
		var batch []resources.Product
		var err error

		select {
		case batch, ok1 = <-products:
			for _, product := range batch {
				result = append(result, product)
			}
			break
		case err, ok2 = <-errors:
			if ok2 {
				e = err
			}

			break
		}
	}

	return result, e
}

func TestProductList(t *testing.T, fileName string, chProducts chan []resources.Product, chErrors chan error) {
	fileName = filepath.Base(fileName)

	products, e := GetParserResult(chProducts, chErrors)

	if e != nil {
		t.Errorf("[%s] error: %s", fileName, e)
	}

	if total := len(products); 0 == total {
		t.Errorf("[%s] no product found", fileName)
	} else {
		statMissing := map[string]int{
			"id":    0,
			"title": 0,
			"desc":  0,
			"image": 0,
			"link":  0,
			"price": 0,
		}

		for _, p := range products {
			if len(p.Id) == 0 {
				statMissing["id"]++
			}
			if len(p.Title) == 0 {
				statMissing["title"]++
			}
			if len(p.Description) == 0 {
				statMissing["desc"]++
			}
			if p.Price().Price < 0.01 {
				statMissing["price"]++
			}
			if !strings.Contains(p.GetImageLink(), "https://") {
				statMissing["image"]++
			}
			if !strings.Contains(p.Link, "https://") {
				statMissing["link"]++
			}
		}

		for field, missing := range statMissing {
			if missing == total {
				t.Errorf("[%s] missing attribute '%s' in product list", fileName, field)
			} else if missing > 0 {
				t.Errorf("[%s] missing attribute '%s' in product list for %d/%d total products. maybe normal ?", fileName, field, missing, total)
			}
		}
	}
}
