package script

import (
	"chapter13/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Script struct {
	Cmds []Command
}

func New(cmds ...Command) *Script {
	return &Script{cmds}
}

func (s Script) String() string {
	builder := strings.Builder{}

	for _, cmd := range s.Cmds {
		if cmd.IsOpCode {
			builder.WriteString(cmd.Code.String())
			builder.WriteByte('\n')
		} else {
			builder.WriteString(hex.EncodeToString(cmd.Elem))
			builder.WriteByte('\n')
		}
	}

	return builder.String()
}

func (s Script) RawSerialize() ([]byte, error) {
	result := []byte{}

	for _, cmd := range s.Cmds {
		// 스크립트 명령어가 연산자에 해당하는 경우: 해당 연산자를 1바이트 리틀엔디언으로 직렬화
		if cmd.IsOpCode {
			result = append(result, utils.IntToLittleEndian(cmd.Code.Int(), 1)...)
			continue
		}

		// 스크립트 명령어가 []byte 타입인 경우: 원소의 길이에 따라 직렬화
		length := len(cmd.Elem)
		if length < 75 { // 원소의 길이가 75보다 작은 경우: 해당 길이를 1바이트 리틀엔디언으로 직렬화
			result = append(result, utils.IntToLittleEndian(length, 1)...)
		} else if length > 75 && length < 0x100 { // 원소의 길이가 75보다 크고 0x100보다 작은 경우: OP_PUSHDATA1에 해당하므로 76을 추가하고 길이를 1바이트 리틀엔디언으로 직렬화
			result = append(result, utils.IntToLittleEndian(76, 1)...) // OP_PUSHDATA1
			result = append(result, utils.IntToLittleEndian(length, 1)...)
		} else if length >= 0x100 && length <= 520 { // 원소의 길이가 0x100보다 크거나 같고 520보다 작거나 같은 경우: OP_PUSHDATA2에 해당하므로 77을 추가하고 길이를 2바이트 리틀엔디언으로 직렬화
			result = append(result, utils.IntToLittleEndian(77, 1)...)
			result = append(result, utils.IntToLittleEndian(length, 2)...)
		} else { // 그 외의 경우: 에러 반환
			return nil, errors.New("too long cmd")
		}

		result = append(result, cmd.Elem...) // 원소 추가
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
func (s Script) Evaluate(z []byte, witness [][]byte) (bool, error) {
	cmds := s.Cmds
	stack := []any{}    // 스택
	altstack := []any{} // 대체 스택

	// 스크립트 명령어를 순회하면서 스택에 데이터를 추가하거나 연산을 수행
	for len(cmds) > 0 {
		cmd := cmds[0]  // 스크립트 명령어 집합의 첫 번째 원소
		cmds = cmds[1:] // 스크립트 명령어 집합의 첫 번째 원소 제거

		var res strings.Builder

		for _, item := range stack {
			switch item.(type) {
			case int:
				res.WriteString(fmt.Sprintf("%d ", item.(int)))
			case []byte:
				if len(item.([]byte)) == 0 {
					res.WriteString("0 ")
				} else {
					res.WriteString(hex.EncodeToString(item.([]byte)))
					res.WriteString(" ")
				}
			}
		}

		fmt.Println("stack:", res.String(), "\ncmds:", cmds, "\ncmd:", cmd)
		fmt.Println()

		// 스크립트 명령어가 연산자에 해당하는 경우
		if cmd.IsOpCode {
			operation := OpCodeFuncs[cmd.Code] // 연산자에 해당하는 함수 가져오기

			switch {
			case cmd.Code >= 99 && cmd.Code <= 100:
				fn, ok := operation.(func(*[]any, *[]Command) bool) // 연산자에 해당하는 함수가 func(*[]any, *[]Command) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("failed to cast evaluate func: %s", cmd.Code.String())
				}

				if !fn(&stack, &cmds) { // stack과 cmds를 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", cmd.Code.String())
				}
			case cmd.Code >= 107 && cmd.Code <= 108:
				fn, ok := operation.(func(*[]any, *[]any) bool) // 연산자에 해당하는 함수가 func(*[]any, *[]any) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("failed to cast evaluate func: %s", cmd.Code.String())
				}

				if !fn(&stack, &altstack) { // stack과 altstack을 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", cmd.Code.String())
				}
			case cmd.Code >= 172 && cmd.Code <= 175:
				fn, ok := operation.(func(*[]any, []byte) bool) // 연산자에 해당하는 함수가 func(*[]any, []byte) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("failed to cast evaluate func: %s", cmd.Code.String())
				}

				if !fn(&stack, z) { // stack과 z를 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", cmd.Code.String())
				}
			default:
				fn, ok := operation.(func(*[]any) bool) // 연산자에 해당하는 함수가 func(*[]any) bool 타입인지 확인
				if !ok {
					return false, fmt.Errorf("failed to cast evaluate func: %s", cmd.Code.String())
				}

				if !fn(&stack) { // stack을 인자로 연산자에 해당하는 함수 호출
					return false, fmt.Errorf("failed to evaluate %s", cmd.Code.String())
				}
			}
			continue
		}

		// 스크립트 명령어가 []byte 타입인 경우: 스택에 추가
		stack = append(stack, cmd.Elem)

		// p2sh 스크립트인 경우: 스크립트 명령어 집합의 첫 번째 원소가 OP_HASH160, 두 번째 원소가 20바이트의 데이터, 세 번째 원소가 OP_EQUAL인지 확인
		if len(cmds) == 3 && cmds[0].Code == OpCodeHash160 && len(cmds[1].Elem) == 20 && cmds[2].Code == OpCodeEqual {
			// 스크립트 명령어 집합의 세 개의 원소를 제거 (OP_HASH160, 20바이트의 데이터, OP_EQUAL)
			cmds = cmds[1:]
			h160 := cmds[0].Elem
			cmds = cmds[2:]

			fmt.Println("Detect P2SH Script")

			if !OpHash160(&stack) { // OP_HASH160 연산자 수행
				return false, errors.New("failed to evaluate OP_HASH160")
			}

			stack = append(stack, h160) // 스택에 20바이트의 데이터 추가

			if !OpEqual(&stack) { // OP_EQUAL 연산자 수행
				return false, errors.New("failed to evaluate OP_EQUAL")
			}

			if !OpVerify(&stack) { // OP_VERIFY 연산자 수행
				return false, errors.New("failed to evaluate OP_VERIFY")
			}

			// 스크립트 명령어 집합에 리딤 스크립트를 추가
			rawRedeem := append(utils.EncodeVarint(len(cmd.Elem)), cmd.Elem...)

			redeemScript, _, err := Parse(rawRedeem)
			if err != nil {
				return false, err
			}

			cmds = append(cmds, redeemScript.Cmds...)
		}

		// p2wpkh 스크립트인 경우: 스택의 원소가 0과 20바이트의 데이터인지 확인
		if len(stack) == 2 && len(stack[0].([]byte)) == 0 && len(stack[1].([]byte)) == 20 {
			fmt.Println("Detect P2WPKH Script")

			h160 := stack[1].([]byte) // 스택의 두 번째 원소를 20바이트의 데이터로 변환
			// 스택의 원소를 모두 제거
			stack = []any{}

			// 명령어 집합에 witness를 추가
			for _, item := range witness {
				cmds = append(cmds, NewElem(item))
			}
			// 명령어 집합에 p2pkh 스크립트를 추가
			cmds = append(cmds, NewP2PKHScript(h160).Cmds...)
		}
	}

	fmt.Println(stack)

	// 스택이 비어있는 경우: 스크립트가 유효하지 않음
	if len(stack) == 0 {
		return false, nil
	}

	switch last := stack[len(stack)-1].(type) {
	// 스택의 마지막 원소가 int 타입인 경우: 해당 원소가 0이 아닌 경우 스크립트가 유효함
	case int:
		if last == 0 {
			return false, nil
		}
	// 스택의 마지막 원소가 []byte 타입인 경우: 해당 원소가 비어있지 않은 경우 스크립트가 유효함
	case []byte:
		if len(last) == 0 {
			return false, nil
		}
	// 그 외의 경우: 에러 반환
	default:
		return false, errors.New("invalid stack element")
	}

	return true, nil
}
