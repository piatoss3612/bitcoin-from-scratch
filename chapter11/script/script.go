package script

import (
	"chapter11/utils"
	"encoding/hex"
	"errors"
	"fmt"
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
		// 스크립트 명령어가 []byte 타입인 경우: 원소의 길이에 따라 다른 방식으로 직렬화
		case []byte:
			length := len(cmd)
			if length < 75 { // 원소의 길이가 75보다 작은 경우: 해당 길이를 1바이트 리틀엔디언으로 직렬화
				result = append(result, utils.IntToLittleEndian(length, 1)...)
			} else if length > 75 && length < 0x100 { // 원소의 길이가 75보다 크고 0x100보다 작은 경우: OP_PUSHDATA1에 해당하므로 76을 추가하고 길이를 1바이트 리틀엔디언으로 직렬화
				result = append(result, 76)
				result = append(result, utils.IntToLittleEndian(length, 1)...)
			} else if length >= 0x100 && length < 520 { // 원소의 길이가 0x100보다 크거나 같고 520보다 작은 경우: OP_PUSHDATA2에 해당하므로 77을 추가하고 길이를 2바이트 리틀엔디언으로 직렬화
				result = append(result, 77)
				result = append(result, utils.IntToLittleEndian(length, 2)...)
			} else { // 그 외의 경우: 에러 반환
				return nil, errors.New("too long an cmd")
			}
			result = append(result, cmd...) // 직렬화한 데이터를 추가
		// 스크립트 명령어가 int 타입인 경우: 연산자에 해당하므로 리틀엔디언으로 직렬화
		case int:
			result = append(result, utils.IntToLittleEndian(cmd, 1)...)
		}
	}

	return result, nil
}

func (s Script) Serialize() ([]byte, error) {
	result, err := s.RawSerialize() // 직렬화한 데이터
	if err != nil {
		return nil, err
	}

	total := len(result) // 직렬화한 데이터의 전체 길이

	// 직렬화한 데이터의 전체 길이를 가변 정수로 직렬화한 뒤 직렬화한 데이터를 추가하여 반환
	return append(utils.EncodeVarint(total), result...), nil
}

func (s Script) Add(other *Script) *Script {
	cmds := append(s.Cmds, other.Cmds...)
	return New(cmds...)
}

// 스크립트 명령어 집합을 순회하면서 스크립트가 유효한지 확인
func (s *Script) Evaluate(z []byte) (bool, error) {
	cmds := make([]any, len(s.Cmds))
	copy(cmds, s.Cmds)  // 스크립트 명령어 집합 복사
	stack := []any{}    // 스택
	altstack := []any{} // 대체 스택

	// 스크립트 명령어를 순회하면서 스택에 데이터를 추가하거나 연산을 수행
	for len(cmds) > 0 {
		cmd := cmds[0]  // 스크립트 명령어 집합의 첫 번째 원소
		cmds = cmds[1:] // 스크립트 명령어 집합의 첫 번째 원소 제거

		switch cmd := cmd.(type) {
		// 스크립트 명령어가 int 타입인 경우: 연산자에 해당하므로 연산을 수행
		case int:
			operation := OpCodeFuncs[OpCode(cmd)] // 연산자에 해당하는 함수 가져오기

			if cmd > 98 && cmd < 101 {
				fn, ok := operation.(func(*[]any, *[]any) bool) // 연산자에 해당하는 함수가 func(*[]any, *[]any) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("operation is not valid: %s", OpCode(cmd).String())
				}

				if !fn(&stack, &cmds) { // stack과 cmds를 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", OpCode(cmd).String())
				}
			} else if cmd > 106 && cmd < 109 {
				fn, ok := operation.(func(*[]any, *[]any) bool) // 연산자에 해당하는 함수가 func(*[]any, *[]any) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("operation is not valid: %s", OpCode(cmd).String())
				}

				if !fn(&stack, &altstack) { // stack과 altstack을 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", OpCode(cmd).String())
				}
			} else if cmd > 171 && cmd < 176 {
				fn, ok := operation.(func(*[]any, []byte) bool) // 연산자에 해당하는 함수가 func(*[]any, []byte) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("operation is not valid: %s", OpCode(cmd).String())
				}

				if !fn(&stack, z) { // stack과 z를 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", OpCode(cmd).String())
				}
			} else {
				fn, ok := operation.(func(*[]any) bool) // 연산자에 해당하는 함수가 func(*[]any) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("operation is not valid: %s", OpCode(cmd).String())
				}

				if !fn(&stack) { // stack을 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", OpCode(cmd).String())
				}
			}
		// 스크립트 명령어가 []byte 타입인 경우: 스택에 원소를 추가
		case []byte:
			stack = append(stack, cmd)

			// cmds 안에 3개의 명령어가 남아있고 BIP0016에서 규정한 특별 패턴에 해당하는 경우
			if len(cmds) == 3 {
				// cmds의 첫 번째 원소가 OP_HASH160, cmds의 두 번째 원소가 20바이트인 []byte 타입, cmds의 세 번째 원소가 OP_EQUAL인지 확인
				opCodeH160, ok1 := cmds[0].(int)
				h160, ok2 := cmds[1].([]byte)
				opCodeEqual, ok3 := cmds[2].(int)

				if ok1 && ok2 && ok3 && opCodeH160 == 0xa9 && len(h160) == 20 && opCodeEqual == 0x87 {
					cmds = cmds[3:] // cmds에서 3개의 명령어 제거

					if !OpHash160(&stack) {
						return false, errors.New("failed to evaluate OP_HASH160")
					}

					stack = append(stack, h160) // 스택에 h160 추가

					if !OpEqual(&stack) {
						return false, errors.New("failed to evaluate OP_EQUAL")
					}

					if !OpVerify(&stack) {
						return false, errors.New("failed to evaluate OP_VERIFY")
					}

					rawRedeem := append(utils.EncodeVarint(len(cmd)), cmd...)

					redeemScript, _, err := Parse(rawRedeem) // redeemScript 파싱
					if err != nil {
						return false, err
					}

					cmds = append(cmds, redeemScript.Cmds...) // cmds에 스크립트 명령어 집합 추가
				}
			}
		// 그 외의 경우: 에러 반환
		default:
			return false, errors.New("invalid cmd")
		}
	}

	// 스택이 비어있는 경우: 스크립트가 유효하지 않음
	if len(stack) == 0 {
		return false, nil
	}

	switch popped := stack[len(stack)-1].(type) {
	// 스택의 마지막 원소가 int 타입인 경우: 해당 원소가 0이 아닌 경우 스크립트가 유효함
	case int:
		if popped == 0 {
			return false, nil
		}
	// 스택의 마지막 원소가 []byte 타입인 경우: 해당 원소가 비어있지 않은 경우 스크립트가 유효함
	case []byte:
		if len(popped) == 0 {
			return false, nil
		}
	// 그 외의 경우: 에러 반환
	default:
		return false, errors.New("invalid stack")
	}

	return true, nil
}
