package script

import (
	"bytes"
	"chapter06/utils"
	"errors"
)

func Parse(b []byte) (*Script, int, error) {
	length, read := utils.ReadVarint(b)

	buf := bytes.NewBuffer(b[read:])

	var cmds []any
	var count int

	for count < length {
		current := buf.Next(1)
		count += 1

		currentByte := current[0]
		if currentByte >= 1 && currentByte <= 75 {
			dataLength := int(currentByte)
			data := buf.Next(dataLength)
			count += dataLength
			cmds = append(cmds, data)
		} else if currentByte == 76 {
			dataLength := utils.LittleEndianToInt(buf.Next(1))
			data := buf.Next(dataLength)
			cmds = append(cmds, data)
			count += dataLength + 1
		} else if currentByte == 77 {
			dataLength := utils.LittleEndianToInt(buf.Next(2))
			data := buf.Next(dataLength)
			cmds = append(cmds, data)
			count += dataLength + 2
		} else {
			opCode := int(currentByte)
			cmds = append(cmds, opCode)
		}
	}

	if count != length {
		return nil, 0, errors.New("parse error: length mismatch")
	}

	return New(cmds...), read + length, nil
}
