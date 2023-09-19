module github.com/louis296/pet

go 1.19

require (
	github.com/suyashkumar/dicom v1.0.6
	google.golang.org/protobuf v1.26.0
	gopkg.in/ini.v1 v1.67.0
)

require (
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/text v0.3.8 // indirect
)

replace github.com/suyashkumar/dicom v1.0.6 => github.com/louis296/dicom-go v1.0.6-0.20220430142348-fed8ad17cc0b
