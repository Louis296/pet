package dpet

import (
	"bytes"
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
	if dataset.DataBuf != nil {
		return writeBinaryData(dataset.DataBuf, writer)
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
	defer fw.Flush()
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
		mich, _ := data.([]uint16)
		return writeMichData930(mich, fw)
	}
	return UnknownFileType
}

func writeDataE180(fileType FileType, data interface{}, writer io.Writer) error {
	fw, err := flate.NewWriter(writer, flate.BestCompression)
	defer fw.Flush()
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
		mich, _ := data.([]float32)
		return writeMichDataE180(mich, fw)
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

func writeBinaryData(buf *bytes.Buffer, writer io.Writer) error {
	fw, err := flate.NewWriter(writer, flate.BestCompression)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, buf)
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
	for _, pair := range data.CoinPairs {
		err = binary.Write(w, binary.LittleEndian, pair[0].GlobalCrystalIndex)
		err = binary.Write(w, binary.LittleEndian, pair[0].Energy)
		err = binary.Write(w, binary.LittleEndian, pair[0].TimeValue)
		err = binary.Write(w, binary.LittleEndian, pair[1].GlobalCrystalIndex)
		err = binary.Write(w, binary.LittleEndian, pair[1].Energy)
		err = binary.Write(w, binary.LittleEndian, pair[1].TimeValue)
		if err != nil {
			return
		}
	}
	return nil
}

func writeMichDataE180(data []float32, w io.Writer) (err error) {
	for i := range data {
		err = binary.Write(w, binary.LittleEndian, data[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func writeRawData930(data *RawData930, w io.Writer) (err error) {
	for _, item := range data.List {
		_, err = w.Write(item.Data)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.LittleEndian, item.IP)
		if err != nil {
			return
		}
	}
	return nil
}

func writeListModeCoinData930(data *ListModeCoinData930, w io.Writer) (err error) {
	for _, item := range data.List {
		err = binary.Write(w, binary.LittleEndian, item.IP)
		ch := uint16(item.Reserved)<<12 + item.Channel
		if item.XTalk {
			ch &= 1 << 15
		}
		err = binary.Write(w, binary.LittleEndian, ch)
		err = binary.Write(w, binary.LittleEndian, item.Energy)
		err = binary.Write(w, binary.LittleEndian, item.Time)
		if err != nil {
			return
		}
	}
	return nil
}

func writeMichData930(data []uint16, w io.Writer) (err error) {
	for i := range data {
		err = binary.Write(w, binary.LittleEndian, data[i])
		if err != nil {
			return err
		}
	}
	return nil
}
