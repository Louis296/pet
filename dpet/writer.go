package dpet

import (
	"compress/flate"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

func Write(dataset *Dataset, writer io.Writer) error {
	fileType := dataset.Header.Content.PublicInfo.FileType
	drive := dataset.Header.Content.ScannerInfo.Device
	err := writeHead(dataset.Header, writer)
	if err != nil {
		return err
	}

	switch drive {
	case "930":
		return writeData930(fileType, dataset, writer)
	case "e180":
		return writeDataE180(fileType, dataset, writer)
	}

	return UnknownDrive
}

func writeData930(fileType FileType, data interface{}, writer io.Writer) error {
	fw, err := flate.NewWriter(writer, flate.BestCompression)
	if err != nil {
		return err
	}
	switch fileType {
	case FileType_RawData:
		rawData, _ := data.(*RawData930)
		return writeRawData930(rawData, fw)
	case FileType_ListModeCoin:
		listMode, _ := data.(*ListModeCoinData930)
		return writeListModeCoinData930(listMode, fw)
	case FileType_Mich:
		//todo
	}
	return UnknownFileType
}

func writeDataE180(fileType FileType, data interface{}, writer io.Writer) error {
	fw, err := flate.NewWriter(writer, flate.BestCompression)
	if err != nil {
		return err
	}
	switch fileType {
	case FileType_RawData:
		rawData, _ := data.(*RawDataE180)
		return writeRawDataE180(rawData, fw)
	case FileType_ListModeCoin:
		listMode, _ := data.(*ListModeCoinDataE180)
		return writeListModeCoinDataE180(listMode, fw)
	case FileType_Mich:
		//todo
	}
	return UnknownFileType
}

func writeHead(header *Header, writer io.Writer) error {
	_, err := writer.Write(MagicNumber)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint16(MarshallMethodProto))
	if err != nil {
		return err
	}
	content, err := proto.Marshal(header.Content)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint32(len(content)))
	if err != nil {
		return err
	}
	_, err = writer.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func writeBinaryData(data []byte, writer io.Writer) error {
	fw, err := flate.NewWriter(writer, flate.BestCompression)
	if err != nil {
		return err
	}
	_, err = fw.Write(data)
	if err != nil {
		return err
	}
	return fw.Flush()
}

func writeRawDataE180(data *RawDataE180, w io.Writer) (err error) {
	for _, info := range data.BDMInfos {
		err = binary.Write(w, binary.LittleEndian, info.BDMIndex)
		err = binary.Write(w, binary.LittleEndian, info.IP)
		err = binary.Write(w, binary.LittleEndian, info.Port)
		err = binary.Write(w, binary.LittleEndian, info.GroupNum)
		err = binary.Write(w, binary.LittleEndian, info.GroupIndex)
		err = binary.Write(w, binary.LittleEndian, info.DataLen)
		if err != nil {
			return
		}
		for _, item := range info.Content {
			err = binary.Write(w, binary.LittleEndian, item.HeadAndDU)
			err = binary.Write(w, binary.LittleEndian, item.BDM)
			err = binary.Write(w, binary.LittleEndian, item.Time)
			err = binary.Write(w, binary.LittleEndian, item.X)
			err = binary.Write(w, binary.LittleEndian, item.Y)
			err = binary.Write(w, binary.LittleEndian, item.Energy)
			err = binary.Write(w, binary.LittleEndian, item.TemperatureInt)
			err = binary.Write(w, binary.LittleEndian, item.TemperatureAndTail)
			if err != nil {
				return
			}
		}
	}
	return nil
}

func writeListModeCoinDataE180(data *ListModeCoinDataE180, w io.Writer) (err error) {
	return nil
}

func writeRawData930(data *RawData930, w io.Writer) (err error) {
	return nil
}

func writeListModeCoinData930(data *ListModeCoinData930, w io.Writer) (err error) {
	return nil
}
