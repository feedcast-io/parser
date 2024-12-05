package handlers

import (
	"encoding/xml"
	"github.com/feedcast-io/parser/resources"
	"io"
	"log/slog"
	"os"
)

func GetFromXml(localFile *os.File) (chan []resources.Product, chan error) {
	d := xml.NewDecoder(localFile)

	slog.Debug("Start process file", "file", localFile.Name())

	ch := make(chan []resources.Product)
	err := make(chan error)

	go func() {
		batch := make([]resources.Product, 0)
		batchSize := 80
		for {
			token, e := d.Token()

			if e == io.EOF {
				break
			} else if e != nil {
				err <- e
				break
			}

			switch tok := token.(type) {
			case xml.StartElement:
				if tok.Name.Local == "item" ||
					tok.Name.Local == "entry" ||
					tok.Name.Local == "product" {
					var doc resources.Product

					if nil == d.DecodeElement(&doc, &tok) {
						batch = append(batch, doc)
						if len(batch) == batchSize {
							ch <- batch
							batch = make([]resources.Product, 0)
						}
					}
				}
				break
			}
		}

		if len(batch) > 0 {
			ch <- batch
		}

		close(ch)
		close(err)
	}()

	return ch, err
}
