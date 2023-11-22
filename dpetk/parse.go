package dpetk

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strconv"
)

func ParseFile(path string, parseData bool) (*DataSet, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	return Parse(file, parseData)
}

func Parse(reader io.Reader, parseData bool) (*DataSet, error) {
	p := &Parser{
		reader:    reader,
		byteOrder: binary.LittleEndian,
		modifyStr: true,

		parseData: parseData,
	}
	return p.parse(), nil
}

type Parser struct {
	reader    io.Reader
	byteOrder binary.ByteOrder

	// 是否移除string末尾的空字符
	modifyStr bool

	// 是否对数据区进行解析，如不解析则将数据放置在dataset中缓冲区
	parseData bool
}

func (p *Parser) parse() *DataSet {
	dataSet := &DataSet{}
	dataSet.PublicInfo = p.parsePublicInfo()
	dataSet.DeviceInfo = p.parseDeviceInfo()
	switch dataSet.PublicInfo.Type {
	case RawDataType:
		dataSet.AcquisitionInfo = p.parseAcquisitionInfo()
		dataSet.DataInfo = p.parseDataInfo()

		if p.parseData {
			dataSet.RawData = p.parseRawData()
		}
	case ListmodeDataType:
		dataSet.AcquisitionInfo = p.parseAcquisitionInfo()
		dataSet.DataInfo = p.parseDataInfo()

		if p.parseData {
			dataSet.ListmodeData = p.parseListmodeData()
		}
	case MichDataType:
		dataSet.AcquisitionInfo = p.parseAcquisitionInfo()
		dataSet.DataInfo = p.parseDataInfo()

		if p.parseData {
			dataSet.MichData = p.parseMichData()
		}
	case EnergyCalibrationMap:
		dataSet.DataInfo = p.parseDataInfo()
	case TimeCalibrationMap:
		dataSet.DataInfo = p.parseDataInfo()
	case EnergySpectrumData:
		dataSet.DataInfo = p.parseDataInfo()
	default:
		dataSet.AcquisitionInfo = p.parseAcquisitionInfo()
		dataSet.ImageInfo = p.parseImageInfo()
		dataSet.DataInfo = p.parseDataInfo()
	}
	if !p.parseData {
		dataSet.DataBuf = bytes.NewBuffer(nil)
		io.Copy(dataSet.DataBuf, p.reader)
	}
	return dataSet
}

func (p *Parser) parsePublicInfo() *PublicInfo {
	// skip magic keys
	_, _ = p.nextString(16)
	return &PublicInfo{
		HeaderCRC:       p.mustNextUint16(),
		Length:          p.mustNextUint32(),
		Type:            p.mustNextUint16(),
		SoftwareVersion: p.mustNextString(16),
		HeaderLength:    p.mustNextUint32(),
	}
}

func (p *Parser) parseDeviceInfo() *DeviceInfo {
	return &DeviceInfo{
		Length:            p.mustNextUint32(),
		Device:            p.mustNextString(16),
		Serial:            p.mustNextString(16),
		AxisDetectors:     p.mustNextUint16(),
		TransDetectors:    p.mustNextUint16(),
		DetectorsRings:    p.mustNextUint16(),
		DetectorsChannels: p.mustNextUint16(),
		IpCounts:          p.mustNextUint16(),
		IpStart:           p.mustNextUint16(),
		ChannelCounts:     p.mustNextUint16(),
		ChannelStart:      p.mustNextUint16(),
		MvtThresholds:     p.mustNextFloat32Slice(8),
		MvtParameters:     p.mustNextFloat32Slice(3),
	}
}

func (p *Parser) parseAcquisitionInfo() *AcquisitionInfo {
	return &AcquisitionInfo{
		Length:             p.mustNextUint32(),
		Isotope:            p.mustNextUint16(),
		Activity:           p.mustNextFloat32(),
		InjectTime:         p.mustNextString(16),
		Time:               p.mustNextString(16),
		Duration:           p.mustNextUint16(),
		TimeWindow:         p.mustNextFloat32(),
		DelayWindow:        p.mustNextFloat32(),
		XTalkWindow:        p.mustNextFloat32(),
		EnergyWindow:       []uint32{p.mustNextUint32(), p.mustNextUint32()},
		PositionWindow:     p.mustNextUint16(),
		Corrected:          p.mustNextUint16(),
		TablePosition:      p.mustNextFloat32(),
		TableHeight:        p.mustNextFloat32(),
		PETCTSpacing:       p.mustNextFloat32(),
		TableCount:         p.mustNextUint16(),
		TableIndex:         p.mustNextUint16(),
		ScanLengthPerTable: p.mustNextFloat32(),
		PatientID:          p.mustNextString(64),
		StudyID:            p.mustNextString(64),
		PatientName:        p.mustNextString(128),
		PatientSex:         p.mustNextString(8),
		PatientHeight:      p.mustNextFloat32(),
		PatientWeight:      p.mustNextFloat32(),
	}
}

func (p *Parser) parseImageInfo() *ImageInfo {
	return &ImageInfo{
		Length:               p.mustNextUint32(),
		ImageSizeRows:        p.mustNextUint16(),
		ImageSizeCols:        p.mustNextUint16(),
		ImageSizeSlices:      p.mustNextUint16(),
		ImageRowPixelSize:    p.mustNextFloat32(),
		ImageColumnPixelSize: p.mustNextFloat32(),
		ImageSliceThickness:  p.mustNextFloat32(),
		ReconMethod:          p.mustNextString(16),
		MaxRingDiffNum:       p.mustNextUint16(),
		SubsetNum:            p.mustNextUint16(),
		IterNum:              p.mustNextUint16(),
		AttnCalibration:      p.mustNextUint16(),
		ScatCalibration:      p.mustNextUint16(),
		ScatPara:             p.mustNextFloat32Slice(6),
		TVPara:               p.mustNextFloat32Slice(2),
		PetCtFovOffset:       p.mustNextFloat32Slice(3),
		CtRotationAngle:      p.mustNextFloat32(),
		SeriesNumber:         p.mustNextUint16(),
		ReconSoftwareVersion: p.mustNextString(16),
		PromptsCounts:        p.mustNextUint32(),
		DelayCounts:          p.mustNextUint32(),
	}
}

func (p *Parser) parseDataInfo() *DataInfo {
	return &DataInfo{
		Length:     p.mustNextUint32(),
		DataLength: p.mustNextUint32(),
		CRC:        p.mustNextUint16(),
	}
}

func (p *Parser) parseRawData() []RawDataItem {
	var res []RawDataItem
	for {
		data, err := p.nextUint8Slice(1152)
		if err != nil {
			break
		}
		res = append(res, RawDataItem{
			Data: data,
			IP:   toIPStr(p.mustNextUint16()),
		})
	}
	return res
}

func (p *Parser) parseListmodeData() []ListmodeDataItem {
	var res []ListmodeDataItem
	for {
		ip, err := p.nextUint16()
		if err != nil {
			break
		}
		ch := p.mustNextUint16()
		res = append(res, ListmodeDataItem{
			IP:       toIPStr(ip),
			XTalk:    ch&(1<<15) != 0,
			Reserved: uint8((ch >> 12) & (1<<3 - 1)),
			Channel:  ch & (1<<12 - 1),
			Energy:   p.mustNextFloat32(),
			Time:     p.mustNextFloat64(),
		})
	}
	return res
}

func (p *Parser) parseMichData() []uint16 {
	var res []uint16
	for {
		v, err := p.nextUint16()
		if err != nil {
			break
		}
		res = append(res, v)
	}
	return res
}

func (p *Parser) nextUint16() (uint16, error) {
	var res uint16
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (p *Parser) nextUint32() (uint32, error) {
	var res uint32
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (p *Parser) nextFloat32() (float32, error) {
	var res float32
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (p *Parser) nextFloat64() (float64, error) {
	var res float64
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (p *Parser) nextString(l int) (string, error) {
	res := make([]byte, l)
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return "", err
	}
	if p.modifyStr {
		return modifyStringByFirstBlank(res), nil
	}
	return string(res), nil
}

func (p *Parser) nextFloat32Slice(l int) ([]float32, error) {
	res := make([]float32, l)
	for i := range res {
		v, err := p.nextFloat32()
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func (p *Parser) nextUint8Slice(l int) ([]uint8, error) {
	res := make([]uint8, l)
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Parser) mustNextUint16() uint16 {
	res, err := p.nextUint16()
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextUint32() uint32 {
	res, err := p.nextUint32()
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextFloat32() float32 {
	res, err := p.nextFloat32()
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextFloat32Slice(l int) []float32 {
	res, err := p.nextFloat32Slice(l)
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextString(l int) string {
	res, err := p.nextString(l)
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextFloat64() float64 {
	res, err := p.nextFloat64()
	if err != nil {
		panic(err)
	}
	return res
}

// modifyString 将bytes转为string，并移除末尾的空字符
func modifyString(bs []byte) string {
	i := len(bs) - 1
	for i >= 0 {
		if bs[i] != 0 && bs[i] != ' ' {
			break
		}
		i--
	}
	return string(bs[:i+1])
}

// modifyStringByFirstBlank 将bytes转为string，并从第一个空白字符处截断
func modifyStringByFirstBlank(bs []byte) string {
	i := 0
	for i < len(bs) {
		if bs[i] == 0 {
			break
		}
		i++
	}
	for i >= 0 {
		if bs[i] != ' ' {
			break
		}
		i--
	}
	return string(bs[:i])
}

func toIPStr(ip uint16) string {
	bs := []byte(ipPrefix)
	bs = append(bs, []byte(strconv.Itoa(int(ip>>8)))...)
	bs = append(bs, '.')
	bs = append(bs, []byte(strconv.Itoa(int(ip&(1<<8-1))))...)
	return string(bs)
}
