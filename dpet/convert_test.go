package dpet

import (
	"fmt"
	"github.com/louis296/pet/dpetk"
	"testing"
)

func TestConvertFrom930(t *testing.T) {
	res := convertFrom930(&dpetk.DataSet{
		PublicInfo: &dpetk.PublicInfo{Type: 0},
		DeviceInfo: &dpetk.DeviceInfo{
			Length:            0,
			Device:            "930",
			Serial:            "121212",
			AxisDetectors:     2,
			TransDetectors:    3,
			DetectorsRings:    4,
			DetectorsChannels: 5,
			IpCounts:          6,
			IpStart:           7,
			ChannelCounts:     8,
			ChannelStart:      9,
			MvtThresholds:     []float32{1.2, 2.1},
			MvtParameters:     []float32{1.1, 2.2},
		},
		AcquisitionInfo: nil,
		ImageInfo:       nil,
		DataInfo:        nil,
		RawData:         nil,
		ListmodeData:    nil,
		MichData:        nil,
		ImageData:       nil,
		DataBuf:         nil,
	})
	fmt.Println(res)
}

func TestConvertFrom930File(t *testing.T) {
	dataset, err := dpetk.ParseFile("../resource/test-raw-data.bin", false)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	datasetNew := convertFrom930(dataset)
	fmt.Println(datasetNew)
}
