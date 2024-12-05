package handlers

import (
	"encoding/csv"
	"encoding/json"
	"github.com/feedcast-io/parser/resources"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
)

func GetFromCsv(localFile *os.File) (chan []resources.Product, chan error) {
	chProducts := make(chan []resources.Product)
	chErrors := make(chan error)

	go func() {
		items := make([]resources.Product, 0)
		batchSize := 100

		possibleSeparators := []rune{
			'|',
			',',
			'\t',
			'~',
		}

		for sepFound, i, maxLoop := false, 0, len(possibleSeparators); !sepFound && i < maxLoop; i++ {
			localFile.Seek(0, io.SeekStart)

			csvReader := csv.NewReader(localFile)
			csvReader.Comma = possibleSeparators[i]
			csvReader.LazyQuotes = true
			csvReader.TrimLeadingSpace = true
			header, _ := csvReader.Read()
			header = sanitizeCsvHeaders(header)

			sepFound = len(header) > 1 && slices.Contains(header, "id")

			if sepFound {
				csvReader.FieldsPerRecord = len(header)
				csvReader.TrimLeadingSpace = false

				for row, err := csvReader.Read(); len(row) > 1; {
					if nil != err {
						slog.Debug("CSV error", "file", localFile.Name(), "err", err.Error())
					}

					if len(row) == len(header) {
						record := make(map[string]string)
						for i, col := range header {
							if len(col) > 0 {
								record[col] = row[i]
							}
						}

						sanitizeProduct(record)

						encoded, _ := json.Marshal(record)

						var item resources.Product
						json.Unmarshal(encoded, &item)
						items = append(items, item)
						if len(items) == batchSize {
							chProducts <- items
							items = make([]resources.Product, 0)
						}
					}
					row, err = csvReader.Read()
				}
			}
		}

		if len(items) > 0 {
			chProducts <- items
		}

		close(chProducts)
		close(chErrors)
	}()

	return chProducts, chErrors
}

func sanitizeCsvHeaders(header []string) []string {
	sanitized := make([]string, len(header))

	var translations = map[string]string{
		"nom":             "title",
		"image":           "image_link",
		"prix":            "price",
		"lien":            "link",
		"price.value":     "price",
		"saleprice.value": "sale_price",
		"i_d":             "id",
		"g_t_i_n":         "gtin",
	}

	snakeCase := func(s string) string {
		var result string

		for _, v := range s {
			if len(result) > 0 && '_' != result[len(result)-1] {
				if v >= 'A' && v <= 'Z' {
					result += "_"
				} else if v >= '0' && v <= '9' {
					result += "_"
				}
			}
			result += string(v)
		}

		return strings.Trim(strings.TrimSpace(strings.ToLower(result)), "_")
	}

	for i, val := range header {
		sanitized[i] = strings.ReplaceAll(val, " ", "_")
		sanitized[i] = strings.ReplaceAll(sanitized[i], "-", "_")
		sanitized[i] = strings.ReplaceAll(sanitized[i], "g:", "")
		sanitized[i] = snakeCase(sanitized[i])
		sanitized[i] = strings.TrimLeft(sanitized[i], "_")

		for k, v := range translations {
			if k == sanitized[i] {
				sanitized[i] = v
			}
		}
	}

	return sanitized
}

func sanitizeProduct(product map[string]string) {
	id, _ := product["id"]
	offerId, _ := product["offer_id"]
	if len(id) == 0 || len(offerId) > 0 {
		product["id"] = offerId
	}
}
