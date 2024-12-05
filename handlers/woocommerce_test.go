package handlers

import (
	"github.com/feedcast-io/parser/internal"
	"github.com/feedcast-io/parser/resources"
	"os"
	"testing"
)

func TestGetFromWoocommerce(t *testing.T) {
	config := resources.WoocommerceConfigImpl{
		Store:     os.Getenv("WOOCOMMERCE_STORE"),
		ApiKey:    os.Getenv("WOOCOMMERCE_KEY"),
		ApiSecret: os.Getenv("WOOCOMMERCE_SECRET"),
	}

	p, e := GetFromWoocommerce(config)

	internal.TestProductList(t, "woocommerce", p, e)

	config.ProductLimit = 20
	p, e = GetFromWoocommerce(config)
	batch, _ := internal.GetParserResult(p, e)

	if config.ProductLimit != len(batch) {
		t.Errorf("expected %d products, got %d", config.ProductLimit, len(batch))
	}
}
