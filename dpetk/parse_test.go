package dpetk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	p := &Parser{
		reader: bytes.NewReader([]byte{
			// uint16: 65535
			0xff, 0xff,
			// uint32: 1
			0x00, 0x00, 0x00, 0x01,
			// float32: 0.1
			0x3d, 0xcc, 0xcc, 0xcd,
			// []float32: {0.1,0.2}
			0x3d, 0xcc, 0xcc, 0xcd, 0x3e, 0x4c, 0xcc, 0xcd,
			// string: "test [nil][nil]"
			0x74, 0x65, 0x73, 0x74, 0x20, 0x00, 0x00,
			// string: " [nil]"
			0x20, 0x00, 0x00,
		}),
		byteOrder: binary.BigEndian,
		modifyStr: true,
	}
	fmt.Println(p.mustNextUint16())
	fmt.Println(p.mustNextUint32())
	fmt.Println(p.mustNextFloat32())
	fmt.Println(p.mustNextFloat32Slice(2))
	fmt.Println(p.mustNextString(7))
	fmt.Println(p.mustNextString(3))
}
