package handlers

import (
	"github.com/feedcast-io/parser/internal"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFromXml(t *testing.T) {
	files, err := filepath.Glob("../../samples/*.xml")

	if err != nil {
		t.Fatal("Fatal", err)
	}

	for _, filename := range files {
		r, err := os.Open(filename)
		if err != nil {
			t.Error(err)
		}
		defer r.Close()

		p, e := GetFromXml(r)

		internal.TestProductList(t, filepath.Base(filename), p, e)
	}
}

func TestGetHugeXmlUnknownFormat(t *testing.T) {
	file, err := os.Open("../../samples/huge-xml.txt")
	if nil != err {
		t.Fatal("Fatal", err)
	}

	p, e := internal.GetParserResult(GetFromXml(file))

	if nil != e {
		t.Error(e)
	}

	if len(p) > 0 {
		t.Error("product size should be 0 when parsing unknown xml format")
	}
}
