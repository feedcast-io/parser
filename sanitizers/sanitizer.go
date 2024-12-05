package sanitizers

import (
	"os"
)

type Sanitizer interface {
	CanProcess(fileHeader string) bool
	Process(file *os.File) (*os.File, error)
}
