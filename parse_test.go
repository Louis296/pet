package pet

import (
	"fmt"
	"github.com/suyashkumar/dicom"
	"testing"
)

func TestParseFile930(t *testing.T) {
	ParseFile930("test.bin")
}

func TestParseFile(t *testing.T) {
	ds, err := dicom.ParseFile("out.dcm", nil)
	fmt.Println(err)
	fmt.Println(ds)
}
