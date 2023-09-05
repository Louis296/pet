package tag

import "github.com/suyashkumar/dicom/pkg/tag"

var PublicInfoGroup uint16 = 0x8080
var DeviceInfoGroup uint16 = 0x8081
var AcquisitionInfoGroup uint16 = 0x8082
var ImageInfoGroup uint16 = 0x8083
var DataInfoGroup uint16 = 0x8084

var HeaderCRC = tag.Tag{Group: PublicInfoGroup, Element: 0x0001}
var PublicInfoLength = tag.Tag{Group: PublicInfoGroup, Element: 0x0002}
var Type = tag.Tag{Group: PublicInfoGroup, Element: 0x0003}
var SoftwareVersion = tag.Tag{Group: PublicInfoGroup, Element: 0x0004}
var HeaderLength = tag.Tag{Group: PublicInfoGroup, Element: 0x0005}
