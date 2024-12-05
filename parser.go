package parser

import (
	"bytes"
	"fmt"
	"github.com/feedcast-io/parser/handlers"
	"github.com/feedcast-io/parser/resources"
	"github.com/feedcast-io/parser/sanitizers"
	"io"
	"log/slog"
	"os"
	"strings"
)

func GetProducts(config resources.Config) (chan []resources.Product, chan error) {
	// Check woocommerce & google sheets first
	if config.Woocommerce != nil && len(config.Woocommerce.GetStore()) > 0 {
		return handlers.GetFromWoocommerce(config.Woocommerce)
	} else if strings.Contains(config.Url, "https://docs.google.com/spreadsheets") {
		return handlers.GetFromGoogleSheet(config.Url)
	}

	// Other remote file source : download and process locally
	localFile, err := downloadFile(config.Url)
	if err == nil {
		defer localFile.Close()
		defer os.Remove(localFile.Name())
	} else {
		return getChanResultError(err)
	}

	return GetProductFromLocalFile(localFile)
}

func GetProductFromLocalFile(localFile *os.File) (chan []resources.Product, chan error) {
	sanitizedFile, header, err := getSanitizedFile(localFile)
	defer os.Remove(sanitizedFile.Name())

	if err != nil {
		return getChanResultError(err)
	}

	sanitizedFile.Seek(0, io.SeekStart)

	if bytes.Contains(header, []byte(`<?xml`)) ||
		bytes.Contains(header, []byte(`<rss`)) ||
		bytes.Contains(header, []byte(`<channel`)) {
		return handlers.GetFromXml(sanitizedFile)
	} else {
		return handlers.GetFromCsv(sanitizedFile)
	}
}

func getSanitizedFile(originFile *os.File) (*os.File, []byte, error) {
	possibleSanitizers := []sanitizers.Sanitizer{
		sanitizers.Win1252Converter{},
		sanitizers.RssFeed{},
		sanitizers.CharacterCleaner{},
	}

	currentReader := originFile
	fileHeader := make([]byte, 2048)
	originFile.Seek(0, io.SeekStart)
	originFile.Read(fileHeader)

	for _, san := range possibleSanitizers {
		if san.CanProcess(string(fileHeader)) {
			newReader, err := san.Process(currentReader)
			if err != nil {
				slog.Warn("Sanitizer error", "type", fmt.Sprintf("%T", san), "error", err.Error())
			} else if nil != newReader {
				if currentReader != originFile {
					defer currentReader.Close()
					defer os.Remove(currentReader.Name())
				}

				currentReader = newReader
			}
		}
	}

	return currentReader, fileHeader, nil
}

func getChanResultError(e error) (chan []resources.Product, chan error) {
	err := make(chan error)
	p := make(chan []resources.Product)

	go func() {
		err <- e
		close(p)
		close(err)
	}()

	return p, err
}
