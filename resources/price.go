package resources

import (
	"regexp"
	"strconv"
	"strings"
)

var currencies = map[string]string{
	"$":   "USD",
	"€":   "EUR",
	"£":   "GBP",
	"£GB": "GBP",
	"CHF": "CHF",
	"PLN": "PLN",
}

type Price struct {
	Price        float32
	CurrencyCode string
	IsDefined    bool
}

func (p *Price) ParseString(raw string) {
	p.Price = 0
	p.CurrencyCode = ""
	p.IsDefined = false

	raw = strings.ToUpper(raw)
	raw = strings.ReplaceAll(raw, " ", "")

	// EUR 1144,0 -> EUR 1144.0
	if strings.Contains(raw, ",") && !strings.Contains(raw, ".") {
		raw = strings.ReplaceAll(raw, ",", ".")
	}

	for k, v := range currencies {
		if strings.HasPrefix(raw, k) {
			raw = strings.ReplaceAll(raw, k, "") + " " + k
		}

		if strings.HasPrefix(raw, v) {
			raw = strings.ReplaceAll(raw, v, "") + " " + v
		}

		if strings.Contains(raw, k) && !strings.Contains(raw, " "+k) {
			raw = strings.ReplaceAll(raw, k, " "+k)
		}

		if strings.Contains(raw, v) && !strings.Contains(raw, " "+v) {
			raw = strings.ReplaceAll(raw, v, " "+v)
		}
	}

	// Unknown currency (UAD, XAT, ...)
	if !strings.Contains(raw, " ") {
		re, _ := regexp.Compile(`([A-Z]{3})`)
		raw = strings.TrimSpace(re.ReplaceAllString(raw, " $1 "))
	}

	parts := strings.Split(strings.TrimSpace(raw), " ")

	if len(parts) < 2 {
		parts = append(parts, "EUR")
	}

	price, lastError := strconv.ParseFloat(parts[0], 16)
	currency := parts[1]

	if nil == lastError {
		p.Price = float32(price)
		p.IsDefined = true
		p.CurrencyCode = currency

		if converted, ok := currencies[p.CurrencyCode]; ok {
			p.CurrencyCode = converted
		}
	}
}
