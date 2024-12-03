package sanitizers

import (
	"io"
	"os"
	"strings"
)

type CharacterCleaner struct{}

func (s CharacterCleaner) CanProcess(fileHeader string) bool {
	return true
}

func (s CharacterCleaner) Process(file *os.File) (*os.File, error) {
	var err error
	var newTempFile *os.File

	if newTempFile, err = os.CreateTemp(os.TempDir(), "feedcast-cr-only"); nil != err {
		return nil, err
	}

	if _, err = file.Seek(0, io.SeekStart); nil != err {
		return nil, err
	}

	for {
		buffer := make([]byte, 1_048_576)
		if _, err := file.Read(buffer); io.EOF == err {
			break
		}

		str := string(buffer)
		// Replace CR with LF for CR files only
		str = strings.ReplaceAll(str, "\r\n", "\n")
		str = strings.ReplaceAll(str, "\n\r", "\n")
		str = strings.ReplaceAll(str, "\r", "\n")
		// Remove backspaces
		str = strings.ReplaceAll(str, string(rune(8)), "")
		str = strings.ReplaceAll(str, string(rune(0)), "")

		if _, err = newTempFile.Write([]byte(str)); nil != err {
			return nil, err
		}
	}

	return newTempFile, nil
}
