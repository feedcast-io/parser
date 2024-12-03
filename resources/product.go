package resources

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

type GetProducts interface {
	GetProducts() []Product
}

type Product struct {
	Id               string   `xml:"id"`
	Title            string   `xml:"title" json:"title"`
	Description      string   `xml:"description" json:"description"`
	Brand            string   `xml:"brand" json:"brand"`
	Mpn              string   `xml:"mpn" json:"mpn"`
	Gtin             string   `xml:"gtin" json:"gtin"`
	Isbn             string   `xml:"isbn" json:"isbn"`
	Ean              string   `xml:"ean" json:"ean"`
	Link             string   `xml:"link" json:"link"`
	MobileLink       string   `xml:"mobile_link" json:"mobile_link"`
	Images           string   `xml:"image_link" json:"image_link"`
	AdditionalImages []string `xml:"additional_image_link"`
	AvailableDate    string   `xml:"availability_date" json:"availability_date"`
	ProductType      string   `xml:"product_type" json:"product_type"`
	Category         string   `xml:"google_product_category" json:"google_product_category"`
	ItemGroupId      string   `xml:"item_group_id" json:"item_group_id"`
	Label0           string   `xml:"custom_label_0" json:"custom_label_0"`
	Label1           string   `xml:"custom_label_1" json:"custom_label_1"`
	Label2           string   `xml:"custom_label_2" json:"custom_label_2"`
	Label3           string   `xml:"custom_label_3" json:"custom_label_3"`
	Label4           string   `xml:"custom_label_4" json:"custom_label_4"`
	Color            string   `xml:"color" json:"color"`
	Material         string   `xml:"material" json:"material"`
	Gender           string   `xml:"gender" json:"gender"`
	AgeGroup         string   `xml:"age_group" json:"age_group"`
	AdsRedirect      string   `xml:"ads_redirect" json:"ads_redirect"`
	Size             string   `xml:"size" json:"size"`
	Quantity         any

	// Fields to transform
	RawAvailability     string          `xml:"availability" json:"availability"`
	RawCondition        string          `xml:"condition" json:"condition"`
	RawIdentifierExists string          `xml:"identifier_exists" json:"identifier_exists"`
	RawIsBundle         string          `xml:"is_bundle" json:"is_bundle"`
	RawPrice            string          `xml:"price" json:"price"`
	RawAdult            string          `xml:"adult" json:"adult"`
	RawSalePrice        string          `xml:"sale_price" json:"sale_price"`
	RawShipping         *ShippingObject `xml:"shipping"`
	RawShippingAsString string          `json:"shipping"`
	RawShippingWeight   string          `xml:"shipping_weight" json:"shipping_weight"`
	RawProductWeight    string          `xml:"product_weight" json:"product_weight"`

	// Alias for other feed formats
	BrandAlt1 string `xml:"manufacturer" json:"manufacturer"`
	ImageAlt1 string `xml:"image" json:"image"`
	PriceAlt1 string `xml:"price_with_vat"`
}

func (p *Product) GetBrand() string {
	if len(p.BrandAlt1) > 0 {
		return p.BrandAlt1
	} else {
		return p.Brand
	}
}

func (p *Product) GetCategoryWithFallback() string {
	if len(p.Category) > 0 {
		return p.Category
	} else {
		return p.ProductType
	}
}

func (p *Product) GetImageLink() string {
	if len(p.Images) > 0 {
		return p.Images
	} else {
		return p.ImageAlt1
	}
}

// Get Price from raw string e.g: 12.34 EUR
func (p *Product) Price() Price {
	price := Price{}

	if len(p.RawPrice) > 0 {
		price.ParseString(p.RawPrice)
	} else if len(p.PriceAlt1) > 0 {
		price.ParseString(p.PriceAlt1)
	}

	return price
}

// Get SalePrice from raw string e.g: 12.34 EUR
func (p *Product) SalePrice() Price {
	price := Price{}
	price.ParseString(p.RawSalePrice)

	return price
}

func (p *Product) GetGenderId() *int8 {
	var value *int8

	switch strings.ToLower(p.Gender) {
	case "homme", "male", "m", "h":
		value = new(int8)
		*value = 1
		break

	case "femme", "female", "f", "w":
		value = new(int8)
		*value = 2
		break

	case "unisex":
		value = new(int8)
		*value = 3
		break
	}

	return value
}

func (p *Product) GetAgeGroupId() *int8 {
	var value *int8

	switch strings.ToLower(p.AgeGroup) {
	case "newborn", "nourrissons":
		value = new(int8)
		*value = 1
		break
	case "infant", "bébés":
		value = new(int8)
		*value = 2
		break
	case "toddler", "tout-petits":
		value = new(int8)
		*value = 3
		break
	case "kids", "enfants":
		value = new(int8)
		*value = 4
		break
	case "adult", "adultes":
		value = new(int8)
		*value = 5
		break
	}

	return value
}

// Get availability from raw string
// 1: in_stock (default)
// 2: out_of_stock
// 3: backorder
func (p *Product) AvailabilityId() int8 {
	avail := int8(1)

	re, _ := regexp.Compile("[^a-z]+")
	test := re.ReplaceAll([]byte(strings.ToLower(p.RawAvailability)), []byte(""))

	if strings.Contains(string(test), "backorder") {
		avail = int8(3)
	} else {
		re, _ = regexp.Compile("(outofstock|non?)")
		if re.Match(test) {
			return int8(2)
		}
	}

	return avail
}

// Get ConditionId from string
// 1: new (default)
// 2: refurbished
// 3: used
func (p *Product) ConditionId() int8 {
	condition := int8(1)

	switch strings.ToLower(p.RawCondition) {
	case "reconditionné", "refurbished":
		condition = int8(2)
		break

	case "used", "occasion":
		condition = int8(3)
		break
	}

	return condition
}

// Get ProductWeight from raw string
func (p *Product) ProductWeight() Weight {
	w := Weight{}

	if len(p.RawProductWeight) > 0 {
		w.ParseString(p.RawProductWeight)
	}

	return w
}

// Get ShippingWeight from raw string
func (p *Product) ShippingWeight() Weight {
	w := Weight{}

	if len(p.RawShippingWeight) > 0 {
		w.ParseString(p.RawShippingWeight)
	}

	return w
}

// Get Shipping from raw string or struct
func (p *Product) Shipping() Shipping {
	shipping := Shipping{}

	if nil != p.RawShipping {
		shipping.FromObject(*p.RawShipping)
	} else if len(p.RawShippingAsString) > 0 {
		shipping.ParseString(p.RawShippingAsString)
	}

	return shipping
}

func (p *Product) HasIdentifier() *int8 {
	return getIntFromString(p.RawIdentifierExists)
}

func (p *Product) IsBundle() *int8 {
	return getIntFromString(p.RawIsBundle)
}

func (p *Product) IsAdult() *int8 {
	return getIntFromString(p.RawAdult)
}

func (p *Product) GetQuantity() *int {
	if i, ok := p.Quantity.(int); ok {
		return &i
	} else if s, ok := p.Quantity.(string); ok {
		if i, err := strconv.Atoi(s); nil == err {
			return &i
		}
	}

	return nil
}

func getIntFromString(value string) *int8 {
	var val int8

	switch strings.ToLower(value) {
	case "no", "non", "false", "0":
		val = 0
		return &val

	case "yes", "true", "oui", "1":
		val = 1
		return &val
	}

	return nil
}

// <mywebstore>
//
//	<products>
//	  <product>
type MyWebstore struct {
	XMLName  xml.Name `xml:"mywebstore"`
	Products struct {
		XMLName xml.Name  `xml:"products"`
		Items   []Product `xml:"product"`
	} `xml:"products"`
	Items []Product `xml:"product"`
}

func (r *MyWebstore) GetProducts() []Product {
	return r.Products.Items
}

// <products>
//
//	<product>
type Products struct {
	XMLName xml.Name  `xml:"products"`
	Items   []Product `xml:"product"`
}

func (r *Products) GetProducts() []Product {
	return r.Items
}

// <root>
//
//	<channel>
//	  <item>
type RootChannel struct {
	XMLName xml.Name `xml:"root"`
	Channel struct {
		XMLName xml.Name  `xml:"channel"`
		Items   []Product `xml:"item"`
	} `xml:"channel"`
}

func (r *RootChannel) GetProducts() []Product {
	return r.Channel.Items
}

// <rss>
//
//	<channel>
//	  <item>...
type RssChannel struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		XMLName xml.Name  `xml:"channel"`
		Items   []Product `xml:"item"`
	} `xml:"channel"`
}

func (rss *RssChannel) GetProducts() []Product {
	return rss.Channel.Items
}

// Format <feed><entry>...
type FeedEntry struct {
	XMLName xml.Name  `xml:"feed"`
	Items   []Product `xml:"entry"`
}

func (f *FeedEntry) GetProducts() []Product {
	return f.Items
}
