package dpet

import (
	"fmt"
	"github.com/louis296/pet/dpetk"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	dataset, _ := dpetk.ParseFile("../resource/test-raw-data.bin", false)
	newDataset := convertFrom930(dataset)

	out, _ := os.Create("../resource/test-gen-raw-data.pet")
	err := Write(newDataset, out)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
}
