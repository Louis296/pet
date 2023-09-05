package dpet

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"os"
)

func ParseFile(path string) (*Dataset, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return Parse(buf)
}

func Parse(buf *bytes.Buffer) (*Dataset, error) {
	magicNumber := buf.Next(4)
	dataset := &Dataset{}
	if string(magicNumber) != string(MagicNumber) {
		return nil, WrongFileTypeError
	}
	header := &Header{}
	header.MarshalMethod = binary.LittleEndian.Uint16(buf.Next(2))
	switch header.MarshalMethod {
	case MarshallMethodProto:
		header.DataLen = binary.LittleEndian.Uint32(buf.Next(4))
		header.Content = &PetFileHeader{}
		err := proto.Unmarshal(buf.Next(int(header.DataLen)), header.Content)
		if err != nil {
			return nil, UnmarshalError
		}
	default:
		return nil, UnknownMarshalMethod
	}
	dataset.Header = header

	// deflate解压
	fr := flate.NewReader(buf)
	decompressBuf := bytes.NewBuffer(nil)
	decompressBuf.ReadFrom(fr)

	switch header.Content.PublicInfo.FileType {
	case FileType_RawData:
		dataset.Data = parseRawData(decompressBuf)
	case FileType_ListModeCoin:
		dataset.Data = parseListModeCoinData(decompressBuf)
	case FileType_Mich:
		dataset.Data = parseMichData(decompressBuf)
	}
	return dataset, nil
}

func parseRawData(buf *bytes.Buffer) *RawDataE180 {
	var infos []*BDMInfo
	for buf.Len() > 0 {
		info := &BDMInfo{
			BDMIndex:   buf.Next(1)[0],
			IP:         binary.LittleEndian.Uint16(buf.Next(2)),
			Port:       binary.LittleEndian.Uint16(buf.Next(2)),
			GroupNum:   buf.Next(1)[0],
			GroupIndex: buf.Next(1)[0],
			DataLen:    binary.LittleEndian.Uint32(buf.Next(4)),
		}
		var content []*BDMInfoBody
		for i := 0; i < (int(info.DataLen) / BDMInfoBodyByteLen); i++ {
			item := &BDMInfoBody{
				HeadAndDU:          buf.Next(1)[0],
				BDM:                buf.Next(1)[0],
				Time:               buf.Next(8),
				X:                  buf.Next(1)[0],
				Y:                  buf.Next(1)[0],
				Energy:             buf.Next(2),
				TemperatureInt:     int8(buf.Next(1)[0]),
				TemperatureAndTail: buf.Next(1)[0],
			}
			content = append(content, item)
		}
		info.Content = content
		infos = append(infos, info)
	}
	return &RawDataE180{BDMInfos: infos}
}

func parseListModeCoinData(buf *bytes.Buffer) *ListModeCoinDataE180 {
	var pairs []CoinPair
	for buf.Len() > 0 {
		pair := [2]*CoinInfo{
			{
				GlobalCrystalIndex: binary.LittleEndian.Uint32(buf.Next(4)),
				Energy:             readFloat32(buf),
				TimeValue:          readFloat64(buf),
			},
			{
				GlobalCrystalIndex: binary.LittleEndian.Uint32(buf.Next(4)),
				Energy:             readFloat32(buf),
				TimeValue:          readFloat64(buf),
			},
		}
		pairs = append(pairs, pair)
	}
	return &ListModeCoinDataE180{CoinPairs: pairs}
}

func parseMichData(buf *bytes.Buffer) []float32 {
	var res []float32
	for buf.Len() > 0 {
		res = append(res, readFloat32(buf))
	}
	return res
}

func readFloat32(buf *bytes.Buffer) float32 {
	var res float32
	binary.Read(buf, binary.LittleEndian, &res)
	return res
}

func readFloat64(buf *bytes.Buffer) float64 {
	var res float64
	binary.Read(buf, binary.LittleEndian, &res)
	return res
}
