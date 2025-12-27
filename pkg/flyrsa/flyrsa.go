package flyrsa

import (
	"encoding/binary"
	"fmt"
)

func removePadding(paddedData []byte) ([]byte, error) {
	if len(paddedData) == 0 {
		return nil, fmt.Errorf("invalid padded data")
	}
	if paddedData[0] != 0x01 {
		return nil, fmt.Errorf("invalid padding")
	}
	dataLen := binary.BigEndian.Uint32(paddedData[1:5])
	if int(dataLen) > len(paddedData)-5 {
		return nil, fmt.Errorf("invalid data length")
	}
	start := len(paddedData) - int(dataLen)
	return paddedData[start:], nil
}
