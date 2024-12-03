package resources

import (
	"testing"
)

func TestProduct_AvailabilityId(t *testing.T) {
	p := Product{}

	if 1 != p.AvailabilityId() {
		t.Errorf("Invalid default availability_id (%d, expected: 1)", p.AvailabilityId())
	}

	p.RawAvailability = "out of stock"
	if 2 != p.AvailabilityId() {
		t.Errorf("Invalid 'out of stock' availability_id (%d, expected: 2)", p.AvailabilityId())
	}

	p.RawAvailability = "backorder"
	if 3 != p.AvailabilityId() {
		t.Errorf("Invalid 'out of stock' availability_id (%d, expected: 3)", p.AvailabilityId())
	}
}

func TestProduct_ConditionId(t *testing.T) {
	p := Product{}

	testValues := map[string]int8{
		"":            1,
		"neuf":        1,
		"new":         1,
		"REFURBISHED": 2,
		"occasion":    3,
		"used":        3,
	}

	for val, expected := range testValues {
		if p.RawCondition = val; expected != p.ConditionId() {
			t.Errorf("Invalid default condition_id '%d' for value '%s' (expected: %d)", p.ConditionId(), val, expected)
		}
	}
}

func TestProduct_GetAgeGroupId(t *testing.T) {
	p := Product{}

	testValues := map[string]int8{
		"newborn": 1,
		"infant":  2,
		"toddler": 3,
		"kids":    4,
		"adult":   5,
	}

	for val, expected := range testValues {
		if p.AgeGroup = val; expected != *p.GetAgeGroupId() {
			t.Errorf("Invalid default age_group_id '%d' for value '%s' (expected: %d)", *p.GetAgeGroupId(), val, expected)
		}
	}
}

func TestProduct_GetQuantity(t *testing.T) {
	p := Product{
		Quantity: 123,
	}

	if 123 != *p.GetQuantity() {
		t.Errorf("Unexpected default quantity '%d' (expected: %d)", 123, *p.GetQuantity())
	}

	p.Quantity = "456"

	if 456 != *p.GetQuantity() {
		t.Errorf("Unexpected default quantity '%d' (expected: %d)", 456, *p.GetQuantity())
	}

	p.Quantity = ""

	if nil != p.GetQuantity() {
		t.Errorf("Quantity should be null, current: '%v'", *p.GetQuantity())
	}
}

func TestProduct_GetGenderId(t *testing.T) {
	p := Product{}

	testValues := map[string]int8{
		"m":      1,
		"h":      1,
		"male":   1,
		"f":      2,
		"female": 2,
		"femme":  2,
		"unisex": 3,
	}

	for val, expected := range testValues {
		if p.Gender = val; expected != *p.GetGenderId() {
			t.Errorf("Invalid default gender_id '%d' for value '%s' (expected: %d)", *p.GetGenderId(), val, expected)
		}
	}
}

func TestProduct_SalePrice(t *testing.T) {
	p := Product{}

	if p.Price().IsDefined {
		t.Errorf("Product price should be not defined")
	}

	p.RawPrice = "12.34 USD"
	price := p.Price()

	if !price.IsDefined {
		t.Errorf("Product price should defined")
	}

	if 12.34 != price.Price {
		t.Errorf("Product price invalid: %f (expected: %f)", price.Price, 12.34)
	}

	if "USD" != price.CurrencyCode {
		t.Errorf("Product currency invalid: %s (expected: %s)", price.CurrencyCode, "USD")
	}

	p.RawPrice = "56.78"
	price = p.Price()
	if !price.IsDefined {
		t.Errorf("Product price should defined")
	}

	if 56.78 != price.Price {
		t.Errorf("Product price invalid: %f (expected: %f)", price.Price, 56.78)
	}

	if "EUR" != price.CurrencyCode {
		t.Errorf("Product currency invalid: %s (expected: %s)", price.CurrencyCode, "EUR")
	}
}

func TestProduct_HasIdentifier(t *testing.T) {
	p := Product{
		RawIdentifierExists: "yes",
	}

	if computed := p.HasIdentifier(); nil == computed || 1 != *computed {
		t.Errorf("HasIdentifier invalid: %d (expected: %d)", *computed, 1)
	}

	p.RawIdentifierExists = "false"
	if computed := p.HasIdentifier(); nil == computed || 0 != *computed {
		t.Errorf("HasIdentifier invalid: %d (expected: %d)", *computed, 0)
	}

	p.RawIdentifierExists = ""
	if computed := p.HasIdentifier(); nil != computed {
		t.Errorf("HasIdentifier invalid: %d (expected: %s)", *computed, "null")
	}
}
