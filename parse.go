package pet

import (
	"github.com/louis296/pet/dpetk"
)

func ParseFile930(path string) (*dpetk.DataSet, error) {
	return dpetk.ParseFile(path, true)
}
