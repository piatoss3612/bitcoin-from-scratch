package script

import (
	"chapter06/utils"
	"encoding/hex"
	"errors"
	"strings"
)

type Script struct {
	Cmds []any
}

func New(cmds ...any) *Script {
	return &Script{cmds}
}

func (s Script) String() string {
	builder := strings.Builder{}

	for _, cmd := range s.Cmds {
		switch cmd := cmd.(type) {
		case []byte:
			builder.WriteString(hex.EncodeToString(cmd))
			builder.WriteString(" ")
		case int:
			builder.WriteString("OP_ ")
		}
	}

	return builder.String()
}

func (s Script) RawSerialize() ([]byte, error) {
	result := []byte{}

	for _, cmd := range s.Cmds {
		switch cmd := cmd.(type) {
		case []byte:
			length := len(cmd)
			if length < 75 {
				result = append(result, utils.IntToLittleEndian(length, 1)...)
			} else if length > 75 && length < 0x100 {
				result = append(result, 76)
				result = append(result, utils.IntToLittleEndian(length, 1)...)
			} else if length >= 0x100 && length < 520 {
				result = append(result, 77)
				result = append(result, utils.IntToLittleEndian(length, 2)...)
			} else {
				return nil, errors.New("too long an cmd")
			}
			result = append(result, cmd...)
		case int:
			result = append(result, utils.IntToLittleEndian(cmd, 1)...)
		}
	}

	return result, nil
}

func (s Script) Serialize() ([]byte, error) {
	result, err := s.RawSerialize()
	if err != nil {
		return nil, err
	}

	total := len(result)

	return append(utils.EncodeVarint(total), result...), nil
}

func (s Script) Add(other *Script) *Script {
	cmds := append(s.Cmds, other.Cmds...)
	return New(cmds...)
}
