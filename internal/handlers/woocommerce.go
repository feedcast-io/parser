package handlers

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/feedcast-io/parser/resources"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type valueObject struct {
	Value string `json:"value"`
}

type serverConfig struct {
	Env struct {
		Version      string `json:"version"`
		WpVersion    string `json:"wp_version"`
		PhpVersion   string `json:"php_version"`
		MysqlVersion string `json:"mysql_version"`
	} `json:"environment"`
}

func createRequest(config resources.WoocommerceConfig, url string, params *url.Values) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", strings.Trim(config.GetStore(), "/"), url), nil)
	req.SetBasicAuth(config.GetApiKey(), config.GetApiSecret())
	req.Header.Set("User-Agent", "curl/8.4.0")
	req.Header.Set("Accept", "*/*")
	if nil != params {
		req.URL.RawQuery = params.Encode()
	}

	return req
}

func getResponseObject[T any](req *http.Request) (T, error) {
	client := http.Client{
		Timeout: time.Minute,
	}

	var resp T
	var buf []byte
	res, lastError := client.Do(req)

	if lastError != nil {
		return resp, lastError
	}

	defer res.Body.Close()

	if 200 == res.StatusCode {
		buf, lastError = io.ReadAll(res.Body)
		if lastError != nil {
			return resp, lastError
		}
		lastError = json.Unmarshal(buf, &resp)
	} else if res.StatusCode >= http.StatusBadRequest {
		return resp, errors.New(res.Status)
	}

	return resp, lastError
}

func getDefaultCurrency(config resources.WoocommerceConfig, defaultCurrency string) (string, error) {
	response, err := getResponseObject[valueObject](createRequest(config, "/wp-json/wc/v3/settings/general/woocommerce_currency", nil))

	return cmp.Or(response.Value, defaultCurrency), err
}

func GetFromWoocommerce(config resources.WoocommerceConfig) (chan []resources.Product, chan error) {
	chProducts := make(chan []resources.Product)
	chErrors := make(chan error)

	var defaultCurrency string

	go func() {
		batchSize, totalProducts := 100, 0

		cfg, lastError := getResponseObject[serverConfig](createRequest(config, "/wp-json/wc/v3/system_status", nil))

		if lastError != nil {
			chErrors <- lastError
		} else {
			slog.Info(
				"Fetch woocommerce config success",
				"woo_version",
				cfg.Env.Version,
				"wp_version",
				cfg.Env.WpVersion,
				"php",
				cfg.Env.PhpVersion,
				"mysql",
				cfg.Env.MysqlVersion,
			)
		}

		defaultCurrency, lastError = getDefaultCurrency(config, "EUR")
		if lastError != nil {
			slog.Warn("Unable to get currency", "use_default", defaultCurrency)
		}

		q := url.Values{}
		q.Set("per_page", strconv.Itoa(batchSize))
		q.Set("status", "publish")
		q.Set("catalog_visibility", "visible")

		for page := 1; nil == lastError; page++ {
			slog.Debug("Fetch page from Woocommerce API", "page", page)

			var products []resources.Product
			var wooProducts []apiProduct

			q.Set("page", strconv.Itoa(page))

			wooProducts, lastError = getResponseObject[[]apiProduct](createRequest(config, "/wp-json/wc/v3/products", &q))

			for _, product := range wooProducts {
				var variants []apiProductVariant

				if "variable" == product.Type {
					// Get product variants
					variants, lastError = getResponseObject[[]apiProductVariant](createRequest(config, fmt.Sprintf("/wp-json/wc/v3/products/%d/variations", product.Id), nil))
				}

				// Add Single Product to list (no variants)
				if 0 == len(variants) {
					if checkAddProduct(config.GetProductLimit(), totalProducts) {
						products = append(products, product.GetProduct(apiProductVariant{}, defaultCurrency))
						totalProducts++
					} else {
						lastError = io.EOF
					}
				}

				for _, variant := range variants {
					if checkAddProduct(config.GetProductLimit(), totalProducts) {
						totalProducts++
						products = append(products, product.GetProduct(variant, defaultCurrency))
					} else {
						lastError = io.EOF
					}
				}
			}

			chProducts <- products

			// End when got fewer results from api than expected (last page)
			if len(wooProducts) < batchSize {
				break
			}

			products = make([]resources.Product, 0)
		}

		close(chProducts)
		close(chErrors)
	}()

	return chProducts, chErrors
}

func checkAddProduct(maxExpected, current int) bool {
	return 0 == maxExpected || current < maxExpected
}

type attribute struct {
	Name    string   `json:"name"`
	Options []string `json:"options"`
}

type image struct {
	Src string `json:"src"`
}

func getFloatFromAny(value interface{}) float32 {
	if s, ok := value.(string); ok {
		f, _ := strconv.ParseFloat(s, 32)
		return float32(f)
	} else if s, ok := value.(float32); ok {
		return s
	} else if s, ok := value.(float64); ok {
		return float32(s)
	}

	return 0
}

type apiProduct struct {
	Id               int32    `json:"id"`
	Gtin             string   `json:"sku"`
	Title            string   `json:"name"`
	Link             string   `json:"permalink"`
	Available        string   `json:"stock_status"`
	Description      string   `json:"description"`
	DescriptionShort string   `json:"short_description"`
	Categories       []string `json:"categories>name"`
	Images           []image  `json:"images"`

	// All prices can be either returned as string or float32 from woocomerce api :(
	Price        interface{} `json:"price"`
	RegularPrice interface{} `json:"regular_price"`
	SalePrice    interface{} `json:"sale_price"`

	Weight     string      `json:"weight"`
	Type       string      `json:"type"`
	Attributes []attribute `json:"attributes"`
}

func (p *apiProduct) GetPrice() float32 {
	return getFloatFromAny(p.Price)
}

func (p *apiProduct) GetRegularPrice() float32 {
	return getFloatFromAny(p.RegularPrice)
}

func (p *apiProduct) GetSalePrice() float32 {
	return getFloatFromAny(p.SalePrice)
}

func (wp *apiProduct) GetProduct(variant apiProductVariant, defaultCurrency string) resources.Product {
	result := resources.Product{
		Id:               fmt.Sprintf("%d", wp.Id),
		Gtin:             wp.Gtin,
		Brand:            wp.Brand(),
		Title:            wp.Title,
		Link:             wp.Link,
		Description:      wp.Description,
		RawAvailability:  wp.Available,
		RawProductWeight: wp.Weight,
	}

	if 0 == len(result.Description) {
		result.Description = wp.DescriptionShort
	}

	if wp.GetRegularPrice() > 0 {
		result.RawPrice = fmt.Sprintf("%.2f %s", wp.GetRegularPrice(), defaultCurrency)
	} else {
		result.RawPrice = fmt.Sprintf("%.2f %s", wp.GetPrice(), defaultCurrency)
	}

	if len(wp.Images) > 0 {
		result.Images = wp.Images[0].Src
	}

	if wp.GetSalePrice() > 0 {
		result.RawSalePrice = fmt.Sprintf("%.2f %s", wp.GetSalePrice(), defaultCurrency)
	}

	if len(wp.Categories) > 0 {
		result.Category = strings.Join(wp.Categories, " > ")
	}

	if variant.Id > 0 {
		result.ItemGroupId = result.Id
		result.Id = fmt.Sprintf("%s-%d", result.Id, variant.Id)
		result.Link = variant.Link
		result.Images = variant.Image.Src
		result.RawAvailability = variant.Available
		result.RawPrice = fmt.Sprintf("%.2f %s", variant.GetPrice(), defaultCurrency)
		result.RawProductWeight = variant.Weight

		if gtin := variant.Gtin(); len(gtin) > 0 {
			result.Gtin = gtin
		}

		if variant.GetSalePrice() > 0 {
			result.RawSalePrice = fmt.Sprintf("%.2f %s", variant.GetSalePrice(), defaultCurrency)
		}
	}

	if len(result.RawProductWeight) > 0 {
		result.RawProductWeight = fmt.Sprintf("%s g", result.RawProductWeight)
		result.RawShippingWeight = result.RawProductWeight
	}

	return result
}

func (wp *apiProduct) getAttribute(name string) string {
	value := ""

	for _, attr := range wp.Attributes {
		if strings.ToLower(attr.Name) == name && len(attr.Options) > 0 {
			value = attr.Options[0]
		}
	}

	return value
}

func (wp *apiProduct) Brand() string {
	brand := wp.getAttribute("brand")

	if len(brand) == 0 {
		brand = wp.getAttribute("marque")
	}

	return brand
}

type metadata struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type apiProductVariant struct {
	Id        int32       `json:"id"`
	Metadata  []metadata  `json:"meta_data"`
	Price     interface{} `json:"regular_price"`
	Link      string      `json:"permalink"`
	SalePrice interface{} `json:"sale_price"`
	Available string      `json:"stock_status"`
	Image     image       `json:"image"`
	Weight    string      `json:"weight"`
}

func (v *apiProductVariant) Gtin() string {
	gtin := ""

	for _, m := range v.Metadata {
		if m.Key == "_alg_ean" {
			if val, ok := m.Value.(string); ok {
				gtin = val
			}
		}
	}

	return gtin
}

func (v *apiProductVariant) GetPrice() float32 {
	return getFloatFromAny(v.Price)
}

func (v *apiProductVariant) GetSalePrice() float32 {
	return getFloatFromAny(v.SalePrice)
}
