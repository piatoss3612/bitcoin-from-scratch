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
			builder.WriteString(OpCode(cmd).String())
			builder.WriteString(" ")
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

func (s *Script) Evaluate(z []byte) bool {
	cmds := s.Cmds
	stack := []any{}
	altstack := []any{}

	for len(cmds) > 0 {
		cmd := cmds[0]
		cmds = cmds[1:]

		switch cmd := cmd.(type) {
		case int:
			operation := OpCodeFuncs[OpCode(cmd)]

			if cmd > 98 && cmd < 101 {
				fn, ok := operation.(func(*[]any, *[]any) bool)
				if !ok {
					return false
				}

				if !fn(&stack, &cmds) {
					return false
				}
			} else if cmd > 106 && cmd < 109 {
				fn, ok := operation.(func(*[]any, *[]any) bool)
				if !ok {
					return false
				}

				if !fn(&stack, &altstack) {
					return false
				}
			} else if cmd > 171 && cmd < 176 {
				fn, ok := operation.(func(*[]any, []byte) bool)
				if !ok {
					return false
				}

				if !fn(&stack, z) {
					return false
				}
			} else {
				fn, ok := operation.(func(*[]any) bool)
				if !ok {
					return false
				}

				if !fn(&stack) {
					return false
				}
			}
		case []byte:
			stack = append(stack, cmd)
		default:
			return false
		}
	}

	if len(stack) == 0 {
		return false
	}

	switch popped := stack[len(stack)-1].(type) {
	case int:
		if popped == 0 {
			return false
		}
	case []byte:
		if len(popped) == 0 {
			return false
		}
	default:
		return false
	}

	return true
}
