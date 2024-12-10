package resources

import "testing"

func TestPrice_ParseString(t *testing.T) {
	p := Price{}

	raw := "12.34 EUR"
	p.ParseString(raw)
	validPrice(t, p, raw, 12.34, "EUR")

	raw = "456 CHF"
	p.ParseString(raw)
	validPrice(t, p, raw, 456, "CHF")

	raw = "32.17USD"
	p.ParseString(raw)
	validPrice(t, p, raw, 32.17, "USD")

	raw = "3333$"
	p.ParseString(raw)
	validPrice(t, p, raw, 3333, "USD")

	raw = "19.90 £ GB"
	p.ParseString(raw)
	validPrice(t, p, raw, 19.90, "GBP")

	raw = "12.18€"
	p.ParseString(raw)
	validPrice(t, p, raw, 12.18, "EUR")

	raw = "24.36£"
	p.ParseString(raw)
	validPrice(t, p, raw, 24.36, "GBP")

	raw = "EUR 12.34"
	p.ParseString(raw)
	validPrice(t, p, raw, 12.34, "EUR")

	raw = "$4.56"
	p.ParseString(raw)
	validPrice(t, p, raw, 4.56, "USD")

	raw = "chf3.37"
	p.ParseString(raw)
	validPrice(t, p, raw, 3.37, "CHF")

	raw = "4.99"
	p.ParseString(raw)
	validPrice(t, p, raw, 4.99, "EUR")

	raw = "EUR 1144,00"
	p.ParseString(raw)
	validPrice(t, p, raw, 1144, "EUR")

	raw = "1 690,00 €"
	p.ParseString(raw)
	validPrice(t, p, raw, 1690, "EUR")

	raw = "1691€"
	p.ParseString(raw)
	validPrice(t, p, raw, 1691, "EUR")

	raw = "1692"
	p.ParseString(raw)
	validPrice(t, p, raw, 1692, "EUR")

	raw = "6995.00XAT"
	p.ParseString(raw)
	validPrice(t, p, raw, 6995, "XAT")

	raw = "6.847,80 EUR"
	p.ParseString(raw)
	validPrice(t, p, raw, 6847.80, "EUR")

	raw = "EUR"
	p.ParseString(raw)
	if p.IsDefined {
		t.Errorf("Price should not have been defined from %s", raw)
	}
}

func validPrice(t *testing.T, testValue Price, raw string, expectedAmount float32, expectedCurrency string) {
	if !testValue.IsDefined {
		t.Errorf("Price should be defined from '%s'", raw)
	}

	if testValue.Price != expectedAmount {
		t.Errorf("Invalid price from '%s'. Expected: %f, current: %f", raw, expectedAmount, testValue.Price)
	}

	if testValue.CurrencyCode != expectedCurrency {
		t.Errorf("Invalid currency from '%s'. Expected: %s, current: %s", raw, expectedCurrency, testValue.CurrencyCode)
	}
}
