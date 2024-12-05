package sanitizers

import (
	"errors"
	"io"
	"os"
	"strings"
)

type RssFeed struct{}

func (s RssFeed) CanProcess(fileHeader string) bool {
	return strings.Contains(fileHeader, "<rss") && strings.Contains(fileHeader, "<channel")
}

func (s RssFeed) Process(file *os.File) (*os.File, error) {
	content := make([]byte, 2048)
	var err error

	st, err := file.Stat()
	if err != nil {
		return nil, err
	}

	tailSize := st.Size()
	if tailSize > 2048 {
		tailSize = 2048
	}

	if _, err = file.Seek(-tailSize, io.SeekEnd); nil != err {
		return nil, errors.New("failed to seek to end of rss feed")
	}

	if _, err = file.Read(content); nil != err {
		return nil, err
	}

	str := string(content)

	if !strings.Contains(str, "</rss>") {
		var newFile *os.File

		if newFile, err = os.CreateTemp(os.TempDir(), "feedcast-rss-sanitizer"); nil != err {
			return nil, err
		}

		if _, err := file.Seek(0, io.SeekStart); nil != err {
			return nil, err
		}

		if _, err := io.Copy(newFile, file); nil != err {
			return nil, err
		}

		if _, err := newFile.Seek(0, io.SeekEnd); nil != err {
			return nil, err
		}

		if !strings.Contains(str, "</channel>") {
			if _, err := newFile.Write([]byte("\n</channel>")); nil != err {
				return nil, err
			}
		}

		if _, err := newFile.Write([]byte("\n</rss>")); nil != err {
			return nil, err
		}

		return newFile, nil
	}

	return nil, nil
}
