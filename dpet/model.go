package dpet

import "bytes"

var MagicNumber = []byte{'D', 'P', 'E', 'T'}

type Dataset struct {
	Header *Header
	Data   interface{}

	// 如果选择不解析data，数据体会以二进制形式存放在此处
	DataBuf *bytes.Buffer
}

type Header struct {
	MarshalMethod uint16
	DataLen       uint32
	Content       *PetFileHeader
}

// RawDataE180 E180原始数据
type RawDataE180 struct {
	BDMInfos []*BDMInfo
}

type BDMInfo struct {
	BDMIndex   uint8
	IP         uint16
	Port       uint16
	GroupNum   uint8
	GroupIndex uint8
	DataLen    uint32
	Content    []*BDMInfoBody
}

type BDMInfoBody struct {
	HeadAndDU          uint8
	BDM                uint8
	Time               []uint8
	X                  uint8
	Y                  uint8
	Energy             []uint8
	TemperatureInt     int8
	TemperatureAndTail uint8
}

// RawData930 930原始数据
type RawData930 struct {
	List []RawDataItem930
}

type RawDataItem930 struct {
	Data []uint8
	IP   uint16
}

// ListModeCoinDataE180 E180符合信息
type ListModeCoinDataE180 struct {
	CoinPairs []CoinPair
}

type CoinPair [2]*CoinInfo

type CoinInfo struct {
	GlobalCrystalIndex uint32
	Energy             float32
	TimeValue          float64
}

// ListModeCoinData930 930符合信息
type ListModeCoinData930 struct {
	List []ListModeDataItem930
}

type ListModeDataItem930 struct {
	IP       uint16
	XTalk    bool
	Reserved uint8
	Channel  uint16
	Energy   float32
	Time     float64
}

const (
	MarshallMethodProto = iota
)

const (
	BDMInfoBodyByteLen = 16
)
const (
	File930  = "DigitMI-930"
	FileE180 = "180"
)
