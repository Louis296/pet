package dpet

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"os"
)

func ParseFile(path string, opt ...ParseOption) (*Dataset, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return Parse(buf, opt...)
}

func Parse(buf *bytes.Buffer, opt ...ParseOption) (*Dataset, error) {
	option := genParseOption(opt...)

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

	if option.onlyHeader {
		return dataset, nil
	}

	// deflate解压
	fr := flate.NewReader(buf)
	decompressBuf := bytes.NewBuffer(nil)
	decompressBuf.ReadFrom(fr)

	if option.notParseData {
		dataset.DataBuf = decompressBuf
		return dataset, nil
	}
	switch header.Content.ScannerInfo.Device {
	case File930:
		parseData930(decompressBuf, header.Content.PublicInfo.FileType)
	case FileE180:
		parseDataE180(decompressBuf, header.Content.PublicInfo.FileType)
	}
	return dataset, nil
}

func parseData930(buf *bytes.Buffer, fileType FileType) interface{} {
	switch fileType {
	case FileType_RawData:
		return parseRawData930(buf)
	case FileType_ListModeCoin:
		return parseListModeCoinData930(buf)
	case FileType_Mich:
		return parseMichData930(buf)
	}
	return nil
}

func parseDataE180(buf *bytes.Buffer, fileType FileType) interface{} {
	switch fileType {
	case FileType_RawData:
		return parseRawDataE180(buf)
	case FileType_ListModeCoin:
		return parseListModeCoinDataE180(buf)
	case FileType_Mich:
		return parseMichDataE180(buf)
	}
	return nil
}

func parseRawDataE180(buf *bytes.Buffer) *RawDataE180 {
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

func parseListModeCoinDataE180(buf *bytes.Buffer) *ListModeCoinDataE180 {
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

func parseMichDataE180(buf *bytes.Buffer) []float32 {
	var res []float32
	for buf.Len() > 0 {
		res = append(res, readFloat32(buf))
	}
	return res
}

func parseRawData930(buf *bytes.Buffer) *RawData930 {
	var res []RawDataItem930
	for {
		data := buf.Next(1152)
		if len(data) == 0 {
			break
		}
		res = append(res, RawDataItem930{
			Data: data,
			IP:   binary.LittleEndian.Uint16(buf.Next(2)),
		})
	}
	return &RawData930{List: res}
}

func parseListModeCoinData930(buf *bytes.Buffer) *ListModeCoinData930 {
	var res []ListModeDataItem930
	for {
		rawIp := buf.Next(2)
		if len(rawIp) == 0 {
			break
		}
		ch := binary.LittleEndian.Uint16(buf.Next(2))
		res = append(res, ListModeDataItem930{
			IP:       binary.LittleEndian.Uint16(rawIp),
			XTalk:    ch&(1<<15) != 0,
			Reserved: uint8((ch >> 12) & (1<<3 - 1)),
			Channel:  ch & (1<<12 - 1),
			Energy:   readFloat32(buf),
			Time:     readFloat64(buf),
		})
	}
	return &ListModeCoinData930{List: res}
}

func parseMichData930(buf *bytes.Buffer) []uint16 {
	var res []uint16
	for buf.Len() > 0 {
		res = append(res, binary.LittleEndian.Uint16(buf.Next(2)))
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
