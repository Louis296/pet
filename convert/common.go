package convert

import "github.com/suyashkumar/dicom"

// NewDicomValue 格式化非数组数据，并提供给dicom写入工具
func NewDicomValue(data interface{}) (dicom.Value, error) {
	switch data.(type) {
	case uint64:
		return dicom.NewValue([]int{int(data.(uint64))})
	case string:
		return dicom.NewValue([]string{data.(string)})
	case byte:
		return dicom.NewValue([]byte{data.(byte)})
	case float64:
		return dicom.NewValue([]float64{data.(float64)})
	case []float32:
		var list []float64
		for _, v := range data.([]float32) {
			list = append(list, float64(v))
		}
		return dicom.NewValue(list)
	case []uint32:
		var list []int
		for _, v := range data.([]uint32) {
			list = append(list, int(v))
		}
		return dicom.NewValue(list)
	default:
		return dicom.NewValue(data)
	}
}
