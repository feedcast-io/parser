package scorer

import (
	"github.com/feedcast-io/parser/resources"
	"testing"
)

func TestScorer_GetScore(t *testing.T) {
	scorer := Scorer{
		Weight: 0.9,
	}
	scorer.HandleProduct(&resources.Product{
		Id:          "1234",
		Title:       "My first title",
		Description: "My description",
		Link:        "http://example.com",
		Images:      "http://example.com/image.jpg",
		Gtin:        "123456789",
	})

	// No error/warnings, score should be max 100*weight
	if 0.9 != scorer.GetScore() {
		t.Errorf("GetScore returned wrong score. Expected 100, got %f", scorer.GetScore())
	}

	// Add 4 products with warnings
	for range []int{1, 2, 3, 4} {
		scorer.HandleProduct(&resources.Product{
			Id:          "1234",
			Title:       "da",
			Description: "My description",
			Link:        "http://example.com",
			Images:      "http://example.com/image.jpg",
			Gtin:        "dfklfdsm",
		})
	}

	if 4 != scorer.warnings {
		t.Errorf("GetScore returned wrong warnings. Expected 4, got %d", scorer.warnings)
	}

	if 0.9 == scorer.GetScore() {
		t.Errorf("GetScore returned wrong score. Expected lest than 0.9, got %f", scorer.GetScore())
	}

	// 1 good products + 4 warnings
	// Warning penalty = 0.2 * 4/5 = 16%
	// Final score = 90% - 16% = 74%
	if 0.74 != scorer.GetScore() {
		t.Errorf("GetScore returned wrong score. Expected 0.74, got %f", scorer.GetScore())
	}

	// Add 10 products with blocking error
	for range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		scorer.HandleProduct(&resources.Product{
			Id:          "1234",
			Title:       "da",
			Description: "My description",
			Gtin:        "dfklfdsm",
		})
	}

	if 4 != scorer.warnings {
		t.Errorf("GetScore returned wrong warnings. Expected 4, got %d", scorer.warnings)
	}
	if 10 != scorer.errors {
		t.Errorf("GetScore returned wrong errors. Expected 10, got %d", scorer.errors)
	}
	// 1 good products + 4 warnings + 10 errors
	// Max score : 0.9 * (total - errors) / total = 0.9 * 5 / 15 = 30%
	// Warning penalty = 0.2 * 4/5 = 16%
	// Final score = 30% - 16% = 14%
	if 14 != int(100*scorer.GetScore()) {
		t.Errorf("GetScore returned wrong score. Expected 0.14, got %f", scorer.GetScore())
	}
}
