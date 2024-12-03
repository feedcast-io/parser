package sanitizers

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strings"
)

// Convert csv from Windows-1252 to UTF-8
type Win1252Converter struct {
	encoded []string
}

func (s Win1252Converter) CanProcess(fileHeader string) bool {
	if strings.Contains(fileHeader, "<?xml") ||
		strings.Contains(fileHeader, "<rss") ||
		strings.Contains(fileHeader, "<channel") {
		return false
	}

	for _, sub := range s.getWin1252Chars() {
		if strings.Contains(fileHeader, sub) {
			return true
		}
	}

	return false
}

func (s Win1252Converter) Process(file *os.File) (*os.File, error) {
	if _, err := file.Seek(0, io.SeekStart); nil != err {
		return nil, err
	}

	var content = make([]byte, 1_048_576)
	if _, err := file.Read(content); nil != err {
		return nil, err
	}

	tempFile, err := os.CreateTemp(os.TempDir(), "feedcast-win1252")
	if nil != err {
		return nil, err
	}

	transformer := transform.NewReader(file, charmap.Windows1252.NewDecoder())
	file.Seek(0, io.SeekStart)
	for {
		var content = make([]byte, 1024)
		if _, err := transformer.Read(content); io.EOF == err {
			break
		}

		tempFile.Write(content)
	}

	return tempFile, nil
}

// Get a list of specific chars, Windows-1252 encoded
func (s Win1252Converter) getWin1252Chars() []string {
	if 0 == len(s.encoded) {
		utf8Characters := []string{
			"à",
			"é",
			"é",
			"ê",
			"è",
			"î",
			"ï",
			"ù",
			"û",
			"ü",
			"ô",
			"ö",
		}

		s.encoded = make([]string, len(utf8Characters))
		encoder := charmap.Windows1252.NewEncoder()

		for i, e := range utf8Characters {
			s.encoded[i], _ = encoder.String(e)
		}
	}

	return s.encoded
}
