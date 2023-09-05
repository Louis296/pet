package dpet

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"gopkg.in/ini.v1"
	"os"
	"testing"
)

var ToLower ini.NameMapper = func(s string) string {
	return string([]byte{s[0] - 'A' + 'a'}) + s[1:]
}

func TestGenDpetFileHeader(t *testing.T) {
	scanInfo := &ScanInfo{}
	acquisition := &AcquisitionInfo{}
	scannerInfo := &ScannerInfo{}
	coincidenceInfo := &CoincidenceInfo{}
	ini.MapToWithMapper(scanInfo, ToLower, "../resource/raw/PET-CT-config.ini")
	ini.MapToWithMapper(acquisition, ToLower, "../resource/raw/acquisition.ini")
	ini.MapToWithMapper(scannerInfo, ToLower, "../resource/coin&mich/scanner.ini")
	ini.MapToWithMapper(coincidenceInfo, ToLower, "../resource/coin&mich/coin-config.ini")

	header := &PetFileHeader{
		PublicInfo: &PublicInfo{
			FileType:           FileType_RawData,
			DataTransferSyntax: DataTransferSyntax_Deflate,
		},
		ScanInfo:        scanInfo,
		AcquisitionInfo: acquisition,
		ScannerInfo:     scannerInfo,
		CoincidenceInfo: coincidenceInfo,
	}

	out, err := proto.Marshal(header)
	if err != nil {
		t.Failed()
	}
	err = os.WriteFile("header", out, 0644)
	if err != nil {
		t.Failed()
	}
}

func TestReadDpetFileHeader(t *testing.T) {
	header := &PetFileHeader{}
	f, _ := os.ReadFile("header")
	err := proto.Unmarshal(f, header)
	if err != nil {
		t.Failed()
	}
	fmt.Print(header)
}

func TestStruct(t *testing.T) {
	f, _ := os.Create("data")
	defer f.Close()
	data := &RawDataE180{BDMInfos: []*BDMInfo{
		{
			BDMIndex:   1,
			IP:         2,
			Port:       3,
			GroupNum:   4,
			GroupIndex: 5,
			DataLen:    6,
			Content: []*BDMInfoBody{
				{
					HeadAndDU:          1,
					BDM:                2,
					Time:               []uint8{3, 4},
					X:                  5,
					Y:                  6,
					Energy:             []uint8{7, 8},
					TemperatureInt:     9,
					TemperatureAndTail: 10,
				},
			},
		},
		{
			BDMIndex:   1,
			IP:         2,
			Port:       3,
			GroupNum:   4,
			GroupIndex: 5,
			DataLen:    6,
			Content: []*BDMInfoBody{
				{
					HeadAndDU:          1,
					BDM:                2,
					Time:               []uint8{3, 4},
					X:                  5,
					Y:                  6,
					Energy:             []uint8{7, 8},
					TemperatureInt:     9,
					TemperatureAndTail: 10,
				},
			},
		},
	}}
	err := writeRawDataE180(data, f)
	if err != nil {
		t.Failed()
	}
}
