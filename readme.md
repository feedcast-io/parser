# Feedcast Parser

Feedcast parser fetches products from xml, csv, Google Sheets for Woocommerce sources.

## Usage

```go
package app

import (
	"github.com/feedcast-io/parser"
	"github.com/feedcast-io/parser/resources"
)

// For a Google Sheets, XML, or CSV url
productsChannel, errorChannel := GetProducts(resources.Config{
    Url: "<feed_url>",
})

// For a Woocommerce site via the API
config := resources.WoocommerceConfigImpl{
    Store:     "<store_url>",
    ApiKey:    "<api_key>",
    ApiSecret: "<api_secret>",
}

productsChannel, errorChannel := GetProducts(config)

// todo: process batch & errors
```

## Configuration

| Variable             | Description                  | Mandatory ?                      |
|----------------------|------------------------------|----------------------------------|
| GOOGLE_CLIENT_ID     | Client Id for Sheets API     | Yes (if using google sheets url) |
| GOOGLE_CLIENT_SECRET | Client Secret for Sheets API | Yes (if using google sheets url) |
| GOOGLE_SHEETS_TOKEN  | Refresh token for Sheets API | Yes (if using google sheets url) |
| WOOCOMMERCE_STORE    | Woocommerce url              | For unit tests only              |
| WOOCOMMERCE_KEY      | Woocommerce Api Key          | For unit tests only              |
| WOOCOMMERCE_SECRET   | Woocommerce Api Secret       | For unit tests only              |