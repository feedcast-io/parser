package parser

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
const acceptLanguage = "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7"

// Download remote file to a random name & returns resource
// File deletion will not be automatic, must be handled
func downloadFile(url string) (*os.File, error) {
	client := &http.Client{
		Timeout: 300 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", acceptLanguage)

	res, err := client.Do(req)

	if nil != err {
		return nil, err
	}

	if http.StatusForbidden == res.StatusCode {
		encoder := base64.StdEncoding
		encoded := encoder.EncodeToString([]byte(url))
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.feedcast.io/proxy/%s", encoded), nil)
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept-Language", acceptLanguage)
		res, err = client.Do(req)
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("Http status %d for %s", res.StatusCode, url))
	}

	defer res.Body.Close()

	tempFile, err := os.CreateTemp(os.TempDir(), "feedcast-downloader")

	if nil != err {
		return nil, err
	}

	if _, err = io.Copy(tempFile, res.Body); nil != err {
		return nil, err
	}

	tempFile.Seek(0, io.SeekStart)

	if st, err := tempFile.Stat(); nil == err && 0 == st.Size() {
		return nil, errors.New("downloaded file has no content")
	}

	return tempFile, nil
}
