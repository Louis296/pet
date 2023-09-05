package convert

import (
	"github.com/louis296/pet/dpetk"
	ptag "github.com/louis296/pet/tag"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	"reflect"
)

type Convertor930 struct {
	Source *dpetk.DataSet
	target dicom.Dataset
}

func (c *Convertor930) writeDicomElement(st interface{}, tagGroup uint16) {
	v := reflect.ValueOf(st)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		e := &dicom.Element{Tag: tag.Tag{Group: tagGroup, Element: uint16(i + 1)}}
		switch f.Kind() {
		case reflect.Uint16:
			e.Value, _ = NewDicomValue(f.Uint())
			e.ValueRepresentation = tag.VRUInt16List
			e.RawValueRepresentation = "US"
		case reflect.Uint32:
			e.Value, _ = NewDicomValue(f.Uint())
			e.ValueRepresentation = tag.VRUInt32List
			e.RawValueRepresentation = "UL"
		case reflect.Float32:
			e.Value, _ = NewDicomValue(f.Float())
			e.ValueRepresentation = tag.VRFloat32List
			e.RawValueRepresentation = "FL"
		case reflect.String:
			e.Value, _ = NewDicomValue(f.String())
			e.ValueRepresentation = tag.VRStringList
			e.RawValueRepresentation = "SH"
		// 对于数组直接写入，与dicom标准的SQ结构有差异
		case reflect.Slice:
			e.Value, _ = NewDicomValue(f.Interface())
			e.ValueRepresentation = tag.VRSequence
			e.RawValueRepresentation = "SQ"
		}
		c.target.Elements = append(c.target.Elements, e)
	}
}

func (c *Convertor930) Convert() (dicom.Dataset, error) {
	c.writeDicomElement(*c.Source.PublicInfo, ptag.PublicInfoGroup)
	c.writeDicomElement(*c.Source.DeviceInfo, ptag.DeviceInfoGroup)
	if c.Source.DeviceInfo != nil {
		c.writeDicomElement(*c.Source.AcquisitionInfo, ptag.AcquisitionInfoGroup)
	}
	if c.Source.ImageInfo != nil {
		c.writeDicomElement(*c.Source.ImageInfo, ptag.ImageInfoGroup)
	}
	if c.Source.DataInfo != nil {
		c.writeDicomElement(*c.Source.DataInfo, ptag.DataInfoGroup)
	}
	return c.target, nil
}
