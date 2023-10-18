package dpet

import (
	"bytes"
	"github.com/louis296/pet/dpetk"
	"reflect"
)

func ParseFrom930(buf *bytes.Buffer, parseData bool) (*Dataset, error) {
	reader := bytes.NewReader(buf.Bytes())
	dataset930, err := dpetk.Parse(reader, parseData)
	if err != nil {
		return nil, err
	}
	return convertFrom930(dataset930), nil
}

func convertFrom930(dataset *dpetk.DataSet) *Dataset {
	petHeader := &PetFileHeader{
		PublicInfo: &PublicInfo{
			FileType:           FileType(dataset.PublicInfo.Type),
			DataTransferSyntax: DataTransferSyntax_Deflate,
			MD5:                "",
		},
		ScanInfo:        &ScanInfo{},
		AcquisitionInfo: &AcquisitionInfo{},
		ScannerInfo:     &ScannerInfo{},
		CoincidenceInfo: &CoincidenceInfo{},
		ImageInfo:       &ImageInfo{},
	}
	switch petHeader.PublicInfo.FileType {
	case FileType_RawData:
		petHeader.CoincidenceInfo = nil
		petHeader.ImageInfo = nil
	case FileType_ListModeCoin:
		petHeader.ImageInfo = nil
	case FileType_Mich:
		petHeader.ImageInfo = nil
	case FileType_EnergyCalibrationMap:
		fallthrough
	case FileType_TimeCalibrationMap:
		fallthrough
	case FileType_EnergySpectrumData:
		fallthrough
	case FileType_PositionTable:
		fallthrough
	case FileType_EnergyMap:
		petHeader.ScanInfo = nil
		petHeader.AcquisitionInfo = nil
		petHeader.CoincidenceInfo = nil
		petHeader.ImageInfo = nil
	}
	fillStructByStructFieldName(petHeader.AcquisitionInfo, dataset.AcquisitionInfo)
	fillStructByStructFieldName(petHeader.ScannerInfo, dataset.DeviceInfo)
	fillStructByStructFieldName(petHeader.ImageInfo, dataset.ImageInfo)

	return &Dataset{
		Header: &Header{
			MarshalMethod: MarshallMethodProto,
			DataLen:       0,
			Content:       petHeader,
		},
		DataBuf: dataset.DataBuf,
	}
}

// fillStructByStructFieldName 通过字段名进行结构体自动装配
func fillStructByStructFieldName(target, source interface{}) {
	if !reflect.ValueOf(target).IsValid() || !reflect.ValueOf(source).IsValid() ||
		reflect.ValueOf(target).IsNil() || reflect.ValueOf(source).IsNil() {
		return
	}
	sourceT := reflect.TypeOf(source).Elem()
	sourceV := reflect.ValueOf(source).Elem()
	targetT := reflect.TypeOf(target).Elem()
	targetV := reflect.ValueOf(target).Elem()
	for k := 0; k < sourceT.NumField(); k++ {
		name := sourceT.Field(k).Name
		value := sourceV.Field(k)
		if t, ok := targetT.FieldByName(name); ok {
			targetV.FieldByName(name).Set(value.Convert(t.Type))
		}
	}
}
