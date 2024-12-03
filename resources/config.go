package resources

type Config struct {
	Url         string
	Woocommerce WoocommerceConfig
}

type WoocommerceConfig interface {
	GetStore() string
	GetApiKey() string
	GetApiSecret() string
	GetProductLimit() int
}

type WoocommerceConfigImpl struct {
	Store        string
	ApiKey       string
	ApiSecret    string
	ProductLimit int
}

func (w WoocommerceConfigImpl) GetStore() string {
	return w.Store
}

func (w WoocommerceConfigImpl) GetApiKey() string {
	return w.ApiKey
}

func (w WoocommerceConfigImpl) GetApiSecret() string {
	return w.ApiSecret
}

func (w WoocommerceConfigImpl) GetProductLimit() int {
	return w.ProductLimit
}
