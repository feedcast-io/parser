package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/feedcast-io/parser/resources"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetFromGoogleSheet(sheetUrl string) (chan []resources.Product, chan error) {
	chErr := make(chan error)
	chProducts := make(chan []resources.Product)
	batchSize := 100

	u, e := url.Parse(sheetUrl)

	if nil != e {
		chErr <- e
		return nil, chErr
	}

	go func() {
		var items []resources.Product
		// Extract sheet id '1sfUBk87240G50bePN0k-VAGtk-FN37_X1kqazn4LlKc'
		// from url /spreadsheets/d/1sfUBk87240G50bePN0k-VAGtk-FN37_X1kqazn4LlKc/edit
		spreadsheetId := strings.ReplaceAll(u.Path, "/edit", "")
		spreadsheetId = strings.ReplaceAll(spreadsheetId, "/spreadsheets/d/", "")

		service := getService()

		spreadsheet, err := service.Spreadsheets.Get(spreadsheetId).Do()

		if nil != err {
			chErr <- err
			return
		}

		r, err := service.Spreadsheets.Values.Get(spreadsheetId, spreadsheet.Sheets[0].Properties.Title).Do()

		if nil != err {
			chErr <- err
			return
		}

		var headers []string

		for _, row := range r.Values {
			if len(headers) == 0 {
				for _, cell := range row {
					if str, ok := cell.(string); ok {
						headers = append(headers, str)
					}
				}

				if len(headers) == 0 {
					chErr <- errors.New("Unable to extract headers from sheet")
					return
				}

				headers = sanitizeCsvHeaders(headers)
			} else {
				record := make(map[string]string)

				for i, cell := range row {
					if str, ok := cell.(string); ok && len(str) > 0 {
						record[headers[i]] = str
					}
				}

				if _, ok := record["id"]; ok {
					var item resources.Product
					data, _ := json.Marshal(record)
					if nil == json.Unmarshal(data, &item) {
						items = append(items, item)
						if len(items) == batchSize {
							chProducts <- items
							items = make([]resources.Product, 0)
						}
					}
				}
			}
		}

		if len(items) > 0 {
			chProducts <- items
		}

		close(chProducts)
		close(chErr)
	}()

	return chProducts, chErr
}

func getService() *sheets.Service {
	c := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly"},
		Endpoint:     google.Endpoint,
	}

	token := &oauth2.Token{
		RefreshToken: os.Getenv("GOOGLE_SHEETS_TOKEN"),
		Expiry:       time.Now().AddDate(0, 0, -1),
	}

	service, _ := sheets.NewService(
		context.Background(),
		option.WithTokenSource(c.TokenSource(context.TODO(), token)),
	)

	return service
}
