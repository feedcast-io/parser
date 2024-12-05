package scorer

import (
	"github.com/feedcast-io/parser/resources"
	"regexp"
)

type Scorer struct {
	found    int
	errors   int
	warnings int
	Weight   float32
}

func (s *Scorer) HandleProduct(product *resources.Product) {
	s.found++

	if len(product.Id) > 50 ||
		0 == len(product.Title) ||
		0 == len(product.Link) ||
		0 == len(product.GetImageLink()) {
		s.errors++
	} else if validGtin, _ := regexp.Match("^[0-9\\-]{8,50}$", []byte(product.Gtin)); !validGtin ||
		len(product.Title) > 150 ||
		0 == len(product.Description) ||
		len(product.Description) > 5000 ||
		len(product.ProductType) > 750 ||
		len(product.GetBrand()) > 70 ||
		len(product.Mpn) > 70 ||
		0 == len(product.Gtin) {
		s.warnings++
	}
}

func (s *Scorer) GetScore() float32 {
	return s.GetScoreFromStats(s.found, s.errors, s.warnings, s.Weight)
}

func (s *Scorer) GetScoreFromStats(found, error, warning int, weight float32) float32 {
	score := float32(0)

	if found > 0 {
		online := float32(found - error)
		score = weight * online / float32(found)

		if online > 0 && warning > 0 {
			penalty := 0.2 * float32(warning) / online
			score -= penalty
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}
