package pet

import (
	"fmt"
	"github.com/louis296/pet/convert"
	"github.com/louis296/pet/dpetk"
	"github.com/suyashkumar/dicom"
	"os"
	"testing"
)

func TestConvertor930_Convert(t *testing.T) {
	dataset, _ := dpetk.ParseFile("test.bin", false)
	c := convert.Convertor930{Source: dataset}
	dicomDataset, _ := c.Convert()
	f, _ := os.Create("out.dcm")
	err := dicom.Write(f, dicomDataset, dicom.DefaultMissingTransferSyntax(), dicom.SkipVRVerification())
	if err != nil {
		fmt.Println(err)
	}
	f.Close()
}
