package resources

import "testing"

func TestShipping_ParseString(t *testing.T) {
	s := Shipping{}
	raw := "FR:::12.34 EUR"
	s.ParseString(raw)
	validShipping(t, s, raw, "FR", 12.34, "EUR")

	raw = "US:$4.33"
	s.ParseString(raw)
	validShipping(t, s, raw, "US", 4.33, "USD")

	raw = ""
	s.ParseString(raw)

	if s.IsDefined {
		t.Errorf("Shipping should not be defined from '%s'", raw)
	}
}

func validShipping(t *testing.T, testValue Shipping, raw string, expectedCountry string, expectedPrice float32, expectedCurrency string) {
	if !testValue.IsDefined {
		t.Errorf("Shipping should be defined from '%s'", raw)
	}
	if expectedCountry != testValue.Country {
		t.Errorf("Invalid shipping country from '%s'. Expected: %s, current: %s", raw, expectedCountry, testValue.Country)
	}
	if expectedCurrency != testValue.Price.CurrencyCode {
		t.Errorf("Invalid shipping currency from '%s'. Expected: %s, current: %s", raw, expectedCurrency, testValue.Price.CurrencyCode)
	}
	if expectedPrice != testValue.Price.Price {
		t.Errorf("Invalid shipping price from '%s'. Expected: %f, current: %f", raw, expectedPrice, testValue.Price.Price)
	}
}

func TestShipping_FromObject(t *testing.T) {
	s := Shipping{}
	raw := ShippingObject{
		Country: "IT",
		Price:   "eur12.56",
	}

	s.FromObject(raw)
	validShippingFromObject(t, s, raw)

	raw = ShippingObject{
		Country: "ES",
		Price:   "€34.22",
	}

	s.FromObject(raw)
	validShippingFromObject(t, s, raw)

	raw = ShippingObject{}
	s.FromObject(raw)
	if s.IsDefined {
		t.Errorf("Shipping should not be defined from empty object")
	}
}

func validShippingFromObject(t *testing.T, testValue Shipping, sourceObject ShippingObject) {
	price := Price{}
	price.ParseString(sourceObject.Price)

	if !testValue.IsDefined {
		t.Errorf("Shipping should be defined")
	}

	if sourceObject.Country != testValue.Country {
		t.Errorf("Invalid shipping country. Expected: %s, current: %s", sourceObject.Country, testValue.Country)
	}
	if price.CurrencyCode != testValue.Price.CurrencyCode {
		t.Errorf("Invalid shipping currency. Expected: %s, current: %s", price.CurrencyCode, testValue.Price.CurrencyCode)
	}
	if price.Price != testValue.Price.Price {
		t.Errorf("Invalid shipping price. Expected: %f, current: %f", price.Price, testValue.Price.Price)
	}
}
