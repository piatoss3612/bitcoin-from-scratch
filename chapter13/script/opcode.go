package script

import (
	"chapter13/ecc"
	"chapter13/utils"
	"log"
)

type OpCode int

const (
	OpCode0                   OpCode = 0x00
	OpCode1Negate             OpCode = 0x4f
	OpCode1                   OpCode = 0x51
	Opcode2                   OpCode = 0x52
	Opcode3                   OpCode = 0x53
	Opcode4                   OpCode = 0x54
	Opcode5                   OpCode = 0x55
	Opcode6                   OpCode = 0x56
	Opcode7                   OpCode = 0x57
	Opcode8                   OpCode = 0x58
	Opcode9                   OpCode = 0x59
	Opcode10                  OpCode = 0x5a
	Opcode11                  OpCode = 0x5b
	Opcode12                  OpCode = 0x5c
	Opcode13                  OpCode = 0x5d
	Opcode14                  OpCode = 0x5e
	Opcode15                  OpCode = 0x5f
	OpCode16                  OpCode = 0x60
	OpCodeNop                 OpCode = 0x61
	OpCodeIf                  OpCode = 0x63
	OpCodeNotIf               OpCode = 0x64
	OpCodeVerify              OpCode = 0x69
	OpCodeReturn              OpCode = 0x6a
	OpCodeToAltStack          OpCode = 0x6b
	OpCodeFromAltStack        OpCode = 0x6c
	OpCode2Drop               OpCode = 0x6d
	OpCode2Dup                OpCode = 0x6e
	OpCode3Dup                OpCode = 0x6f
	OpCode2Over               OpCode = 0x70
	OpCode2Rot                OpCode = 0x71
	OpCode2Swap               OpCode = 0x72
	OpCodeIfDup               OpCode = 0x73
	OpCodeDepth               OpCode = 0x74
	OpCodeDrop                OpCode = 0x75
	OpCodeDup                 OpCode = 0x76
	OpCodeNip                 OpCode = 0x77
	OpCodeOver                OpCode = 0x78
	OpCodePick                OpCode = 0x79
	OpCodeRoll                OpCode = 0x7a
	OpCodeRot                 OpCode = 0x7b
	OpCodeSwap                OpCode = 0x7c
	OpCodeTuck                OpCode = 0x7d
	OpCodeSize                OpCode = 0x82
	OpCodeEqual               OpCode = 0x87
	OpCodeEqualVerify         OpCode = 0x88
	OpCode1Add                OpCode = 0x8b
	OpCode1Sub                OpCode = 0x8c
	OpCodeNegate              OpCode = 0x8f
	OpCodeAbs                 OpCode = 0x90
	OpCodeNot                 OpCode = 0x91
	OpCode0NotEqual           OpCode = 0x92
	OpCodeAdd                 OpCode = 0x93
	OpCodeSub                 OpCode = 0x94
	OpCodeMul                 OpCode = 0x95
	OpCodeBoolAnd             OpCode = 0x9a
	OpCodeBoolOr              OpCode = 0x9b
	OpCodeNumEqual            OpCode = 0x9c
	OpCodeNumEqualVerify      OpCode = 0x9d
	OpCodeNumNotEqual         OpCode = 0x9e
	OpCodeLessThan            OpCode = 0x9f
	OpCodeGreaterThan         OpCode = 0xa0
	OpCodeLessThanOrEqual     OpCode = 0xa1
	OpCodeGreaterThanOrEqual  OpCode = 0xa2
	OpCodeMin                 OpCode = 0xa3
	OpCodeMax                 OpCode = 0xa4
	OpCodeWithin              OpCode = 0xa5
	OpCodeRipemd160           OpCode = 0xa6
	OpCodeSha1                OpCode = 0xa7
	OpCodeSha256              OpCode = 0xa8
	OpCodeHash160             OpCode = 0xa9
	OpCodeHash256             OpCode = 0xaa
	OpCodeCheckSig            OpCode = 0xac
	OpCodeCheckSigVerify      OpCode = 0xad
	OpCodeCheckMultiSig       OpCode = 0xae
	OpCodeCheckMultiSigVerify OpCode = 0xaf
	OpCodeNop1                OpCode = 0xb0
	OpCodeCheckLockTimeVerify OpCode = 0xb1
	OpCodeCheckSequenceVerify OpCode = 0xb2
	OpCodeNop4                OpCode = 0xb3
	OpCodeNop5                OpCode = 0xb4
	OpCodeNop6                OpCode = 0xb5
	OpCodeNop7                OpCode = 0xb6
	OpCodeNop8                OpCode = 0xb7
	OpCodeNop9                OpCode = 0xb8
	OpCodeNop10               OpCode = 0xb9
)

var OpCodeNames = map[OpCode]string{
	0:   "OP_0",
	76:  "OP_PUSHDATA1",
	77:  "OP_PUSHDATA2",
	78:  "OP_PUSHDATA4",
	79:  "OP_1NEGATE",
	81:  "OP_1",
	82:  "OP_2",
	83:  "OP_3",
	84:  "OP_4",
	85:  "OP_5",
	86:  "OP_6",
	87:  "OP_7",
	88:  "OP_8",
	89:  "OP_9",
	90:  "OP_10",
	91:  "OP_11",
	92:  "OP_12",
	93:  "OP_13",
	94:  "OP_14",
	95:  "OP_15",
	96:  "OP_16",
	97:  "OP_NOP",
	99:  "OP_IF",
	100: "OP_NOTIF",
	103: "OP_ELSE",
	104: "OP_ENDIF",
	105: "OP_VERIFY",
	106: "OP_RETURN",
	107: "OP_TOALTSTACK",
	108: "OP_FROMALTSTACK",
	109: "OP_2DROP",
	110: "OP_2DUP",
	111: "OP_3DUP",
	112: "OP_2OVER",
	113: "OP_2ROT",
	114: "OP_2SWAP",
	115: "OP_IFDUP",
	116: "OP_DEPTH",
	117: "OP_DROP",
	118: "OP_DUP",
	119: "OP_NIP",
	120: "OP_OVER",
	121: "OP_PICK",
	122: "OP_ROLL",
	123: "OP_ROT",
	124: "OP_SWAP",
	125: "OP_TUCK",
	130: "OP_SIZE",
	135: "OP_EQUAL",
	136: "OP_EQUALVERIFY",
	139: "OP_1ADD",
	140: "OP_1SUB",
	143: "OP_NEGATE",
	144: "OP_ABS",
	145: "OP_NOT",
	146: "OP_0NOTEQUAL",
	147: "OP_ADD",
	148: "OP_SUB",
	149: "OP_MUL",
	154: "OP_BOOLAND",
	155: "OP_BOOLOR",
	156: "OP_NUMEQUAL",
	157: "OP_NUMEQUALVERIFY",
	158: "OP_NUMNOTEQUAL",
	159: "OP_LESSTHAN",
	160: "OP_GREATERTHAN",
	161: "OP_LESSTHANOREQUAL",
	162: "OP_GREATERTHANOREQUAL",
	163: "OP_MIN",
	164: "OP_MAX",
	165: "OP_WITHIN",
	166: "OP_RIPEMD160",
	167: "OP_SHA1",
	168: "OP_SHA256",
	169: "OP_HASH160",
	170: "OP_HASH256",
	171: "OP_CODESEPARATOR",
	172: "OP_CHECKSIG",
	173: "OP_CHECKSIGVERIFY",
	174: "OP_CHECKMULTISIG",
	175: "OP_CHECKMULTISIGVERIFY",
	176: "OP_NOP1",
	177: "OP_CHECKLOCKTIMEVERIFY",
	178: "OP_CHECKSEQUENCEVERIFY",
	179: "OP_NOP4",
	180: "OP_NOP5",
	181: "OP_NOP6",
	182: "OP_NOP7",
	183: "OP_NOP8",
	184: "OP_NOP9",
	185: "OP_NOP10",
}

var OpCodeFuncs = map[OpCode]any{
	OpCode0:                   Op0,
	OpCode1Negate:             Op1Negate,
	OpCode1:                   Op1,
	Opcode2:                   Op2,
	Opcode3:                   Op3,
	Opcode4:                   Op4,
	Opcode5:                   Op5,
	Opcode6:                   Op6,
	Opcode7:                   Op7,
	Opcode8:                   Op8,
	Opcode9:                   Op9,
	Opcode10:                  Op10,
	Opcode11:                  Op11,
	Opcode12:                  Op12,
	Opcode13:                  Op13,
	Opcode14:                  Op14,
	Opcode15:                  Op15,
	OpCode16:                  Op16,
	OpCodeNop:                 OpNop,
	OpCodeIf:                  OpIf,
	OpCodeNotIf:               OpNotIf,
	OpCodeVerify:              OpVerify,
	OpCodeReturn:              OpReturn,
	OpCodeToAltStack:          OpToAltStack,
	OpCodeFromAltStack:        OpFromAltStack,
	OpCode2Drop:               Op2Drop,
	OpCode2Dup:                Op2Dup,
	OpCode3Dup:                Op3Dup,
	OpCode2Over:               Op2Over,
	OpCode2Rot:                Op2Rot,
	OpCode2Swap:               Op2Swap,
	OpCodeIfDup:               OpIfDup,
	OpCodeDepth:               OpDepth,
	OpCodeDrop:                OpDrop,
	OpCodeDup:                 OpDup,
	OpCodeNip:                 OpNip,
	OpCodeOver:                OpOver,
	OpCodePick:                OpPick,
	OpCodeRoll:                OpRoll,
	OpCodeRot:                 OpRot,
	OpCodeSwap:                OpSwap,
	OpCodeTuck:                OpTuck,
	OpCodeSize:                OpSize,
	OpCodeEqual:               OpEqual,
	OpCodeEqualVerify:         OpEqualVerify,
	OpCode1Add:                Op1Add,
	OpCode1Sub:                Op1Sub,
	OpCodeNegate:              OpNegate,
	OpCodeAbs:                 OpAbs,
	OpCodeNot:                 OpNot,
	OpCode0NotEqual:           Op0NotEqual,
	OpCodeAdd:                 OpAdd,
	OpCodeSub:                 OpSub,
	OpCodeMul:                 OpMul,
	OpCodeBoolAnd:             OpBoolAnd,
	OpCodeBoolOr:              OpBoolOr,
	OpCodeNumEqual:            OpNumEqual,
	OpCodeNumEqualVerify:      OpNumEqualVerify,
	OpCodeNumNotEqual:         OpNumNotEqual,
	OpCodeLessThan:            OpLessThan,
	OpCodeGreaterThan:         OpGreaterThan,
	OpCodeLessThanOrEqual:     OpLessThanOrEqual,
	OpCodeGreaterThanOrEqual:  OpGreaterThanOrEqual,
	OpCodeMin:                 OpMin,
	OpCodeMax:                 OpMax,
	OpCodeWithin:              OpWithin,
	OpCodeRipemd160:           OpRipemd160,
	OpCodeSha1:                OpSha1,
	OpCodeSha256:              OpSha256,
	OpCodeHash160:             OpHash160,
	OpCodeHash256:             OpHash256,
	OpCodeCheckSig:            OpCheckSig,
	OpCodeCheckSigVerify:      OpCheckSigVerify,
	OpCodeCheckMultiSig:       OpCheckMultiSig,
	OpCodeCheckMultiSigVerify: OpCheckMultiSigVerify,
}

func (o OpCode) Valid() bool {
	_, ok := OpCodeNames[o]
	return ok
}

func (o OpCode) String() string {
	s, ok := OpCodeNames[o]
	if !ok {
		return "UNKNOWN"
	}
	return s
}

func (o OpCode) Int() int {
	return int(o)
}

func EncodeNum(num int) []byte {
	if num == 0 {
		return []byte{}
	}

	absNum := func(num int) int {
		if num < 0 {
			return -num
		}

		return num
	}(num)
	negative := num < 0

	result := []byte{}

	for absNum > 0 {
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	if result[len(result)-1]&0x80 != 0 {
		if negative {
			result = append(result, 0x80)
		} else {
			result = append(result, 0)
		}
	} else if negative {
		result[len(result)-1] |= 0x80
	}

	return result
}

func DecodeNum(b []byte) int {
	if len(b) == 0 {
		return 0
	}

	result := 0
	negative := false

	bigEndian := utils.ReverseBytes(b)

	if bigEndian[0]&0x80 != 0 {
		negative = true
		result = int(bigEndian[0] & 0x7f)
	} else {
		result = int(bigEndian[0])
	}

	for i := 1; i < len(bigEndian); i++ {
		result <<= 8
		result |= int(bigEndian[i])
	}

	if negative {
		return -result
	}

	return result
}

func Op0(s *[]any) bool {
	*s = append(*s, EncodeNum(0))
	return true
}

func Op1Negate(s *[]any) bool {
	*s = append(*s, EncodeNum(-1))
	return true
}

func Op1(s *[]any) bool {
	*s = append(*s, EncodeNum(1))
	return true
}

func Op2(s *[]any) bool {
	*s = append(*s, EncodeNum(2))
	return true
}

func Op3(s *[]any) bool {
	*s = append(*s, EncodeNum(3))
	return true
}

func Op4(s *[]any) bool {
	*s = append(*s, EncodeNum(4))
	return true
}

func Op5(s *[]any) bool {
	*s = append(*s, EncodeNum(5))
	return true
}

func Op6(s *[]any) bool {
	*s = append(*s, EncodeNum(6))
	return true
}

func Op7(s *[]any) bool {
	*s = append(*s, EncodeNum(7))
	return true
}

func Op8(s *[]any) bool {
	*s = append(*s, EncodeNum(8))
	return true
}

func Op9(s *[]any) bool {
	*s = append(*s, EncodeNum(9))
	return true
}

func Op10(s *[]any) bool {
	*s = append(*s, EncodeNum(10))
	return true
}

func Op11(s *[]any) bool {
	*s = append(*s, EncodeNum(11))
	return true
}

func Op12(s *[]any) bool {
	*s = append(*s, EncodeNum(12))
	return true
}

func Op13(s *[]any) bool {
	*s = append(*s, EncodeNum(13))
	return true
}

func Op14(s *[]any) bool {
	*s = append(*s, EncodeNum(14))
	return true
}

func Op15(s *[]any) bool {
	*s = append(*s, EncodeNum(15))
	return true
}

func Op16(s *[]any) bool {
	*s = append(*s, EncodeNum(16))
	return true
}

func OpNop(s *[]any) bool {
	return true
}

func OpIf(s *[]any, i *[]Command) bool {
	if len(*s) < 1 {
		return false
	}

	trueItems := []Command{}
	falseItems := []Command{}
	currentItems := &trueItems
	found := false
	numEndIfsNeeded := 1

	for len(*i) > 0 {
		item := (*i)[0]
		*i = (*i)[1:]

		if item.IsOpCode {
			if item.Code > 98 && item.Code < 101 {
				numEndIfsNeeded++
				*currentItems = append(*currentItems, item)
			} else if numEndIfsNeeded == 1 && item.Code == 103 {
				currentItems = &falseItems
			} else if item.Code == 104 {
				if numEndIfsNeeded == 1 {
					found = true
					break
				} else {
					numEndIfsNeeded--
					*currentItems = append(*currentItems, item)
				}
			} else {
				*currentItems = append(*currentItems, item)
			}
		} else {
			*currentItems = append(*currentItems, item)
		}
	}

	if !found {
		return false
	}

	condition := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	if DecodeNum(condition.([]byte)) == 0 {
		*i = append(*i, falseItems...)
	} else {
		*i = append(*i, trueItems...)
	}

	return true
}

func OpNotIf(s *[]any, i *[]Command) bool {
	if len(*s) < 1 {
		return false
	}

	trueItems := []Command{}
	falseItems := []Command{}
	currentItems := &trueItems
	found := false
	numEndIfsNeeded := 1

	for len(*i) > 0 {
		item := (*i)[0]
		*i = (*i)[1:]

		if item.IsOpCode {
			if item.Code > 98 && item.Code < 101 {
				numEndIfsNeeded++
				*currentItems = append(*currentItems, item)
			} else if numEndIfsNeeded == 1 && item.Code == 103 {
				currentItems = &falseItems
			} else if item.Code == 104 {
				if numEndIfsNeeded == 1 {
					found = true
					break
				} else {
					numEndIfsNeeded--
					*currentItems = append(*currentItems, item)
				}
			} else {
				*currentItems = append(*currentItems, item)
			}
		} else {
			*currentItems = append(*currentItems, item)
		}
	}

	if !found {
		return false
	}

	condition := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	if DecodeNum(condition.([]byte)) == 0 {
		*i = append(*i, trueItems...)
	} else {
		*i = append(*i, falseItems...)
	}

	return true
}

func OpVerify(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		if DecodeNum(element) == 0 {
			return false
		}
		return true
	default:
		return false
	}
}

func OpReturn(s *[]any) bool {
	return false
}

func OpToAltStack(s, a *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	*a = append(*a, element)
	return true
}

func OpFromAltStack(s, a *[]any) bool {
	if len(*a) < 1 {
		return false
	}

	element := (*a)[len(*a)-1]
	*a = (*a)[:len(*a)-1]
	*s = append(*s, element)
	return true
}

func Op2Drop(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	*s = (*s)[:len(*s)-2]
	return true
}

func Op2Dup(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	*s = append(*s, (*s)[len(*s)-2:]...)
	return true
}

func Op3Dup(s *[]any) bool {
	if len(*s) < 3 {
		return false
	}

	*s = append(*s, (*s)[len(*s)-3:]...)
	return true
}

func Op2Over(s *[]any) bool {
	if len(*s) < 4 {
		return false
	}

	*s = append(*s, (*s)[len(*s)-4:len(*s)-2]...)
	return true
}

func Op2Rot(s *[]any) bool {
	if len(*s) < 6 {
		return false
	}

	*s = append(*s, (*s)[len(*s)-6:len(*s)-4]...)
	return true
}

func Op2Swap(s *[]any) bool {
	if len(*s) < 4 {
		return false
	}

	(*s)[len(*s)-4], (*s)[len(*s)-3], (*s)[len(*s)-2], (*s)[len(*s)-1] = (*s)[len(*s)-2], (*s)[len(*s)-1], (*s)[len(*s)-4], (*s)[len(*s)-3]
	return true
}

func OpIfDup(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]

	switch element := element.(type) {
	case []byte:
		if DecodeNum(element) != 0 {
			*s = append(*s, element)
		}
		return true
	default:
		return false
	}
}

func OpDepth(s *[]any) bool {
	*s = append(*s, EncodeNum(len(*s)))
	return true
}

func OpDrop(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	*s = (*s)[:len(*s)-1]
	return true
}

func OpDup(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = append(*s, element)
	return true
}

func OpNip(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	*s = append((*s)[:len(*s)-2], (*s)[len(*s)-1])
	return true
}

func OpOver(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	*s = append(*s, (*s)[len(*s)-2])
	return true
}

func OpPick(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		num := DecodeNum(element)
		if len(*s) < num+1 {
			return false
		}
		*s = append(*s, (*s)[len(*s)-num-1])
		return true
	default:
		return false
	}
}

func OpRoll(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		num := DecodeNum(element)
		if len(*s) < num+1 {
			return false
		}
		if num == 0 {
			return true
		}
		temp := (*s)[len(*s)-num-1]
		*s = append((*s)[:len(*s)-num-1], (*s)[len(*s)-num:]...)
		*s = append(*s, temp)
		return true
	default:
		return false
	}
}

func OpRot(s *[]any) bool {
	if len(*s) < 3 {
		return false
	}

	*s = append((*s)[:len(*s)-3], (*s)[len(*s)-2], (*s)[len(*s)-1], (*s)[len(*s)-3])
	return true
}

func OpSwap(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	(*s)[len(*s)-2], (*s)[len(*s)-1] = (*s)[len(*s)-1], (*s)[len(*s)-2]
	return true
}

func OpTuck(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	*s = append((*s)[:len(*s)-2], (*s)[len(*s)-1], (*s)[len(*s)-2], (*s)[len(*s)-1])
	return true
}

func OpSize(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]

	switch element := element.(type) {
	case []byte:
		*s = append(*s, EncodeNum(len(element)))
		return true
	default:
		return false
	}
}

func OpEqual(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if len(x1) != len(x2) {
				*s = append(*s, EncodeNum(0))
				return true
			}

			for i := 0; i < len(x1); i++ {
				if x1[i] != x2[i] {
					*s = append(*s, EncodeNum(0))
					return true
				}
			}

			*s = append(*s, EncodeNum(1))
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpEqualVerify(s *[]any) bool {
	return OpEqual(s) && OpVerify(s)
}

func Op1Add(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		num := DecodeNum(element)
		*s = append((*s)[:len(*s)-1], EncodeNum(num+1))
		return true
	default:
		return false
	}
}

func Op1Sub(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		num := DecodeNum(element)
		*s = append((*s)[:len(*s)-1], EncodeNum(num-1))
		return true
	default:
		return false
	}
}

func OpNegate(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		num := DecodeNum(element)
		*s = append((*s)[:len(*s)-1], EncodeNum(-num))
		return true
	default:
		return false
	}
}

func OpAbs(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		num := DecodeNum(element)
		if num < 0 {
			num = -num
		}
		*s = append((*s)[:len(*s)-1], EncodeNum(num))
		return true
	default:
		return false
	}
}

func OpNot(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		if DecodeNum(element) == 0 {
			*s = append(*s, EncodeNum(1))
		} else {
			*s = append(*s, EncodeNum(0))
		}
		return true
	default:
		return false
	}
}

func Op0NotEqual(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		if DecodeNum(element) == 0 {
			*s = append(*s, EncodeNum(0))
		} else {
			*s = append(*s, EncodeNum(1))
		}
		return true
	default:
		return false
	}
}

func OpAdd(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			num1 := DecodeNum(x1)
			num2 := DecodeNum(x2)
			*s = append(*s, EncodeNum(num1+num2))
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpSub(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			num1 := DecodeNum(x1)
			num2 := DecodeNum(x2)
			*s = append(*s, EncodeNum(num2-num1))
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpMul(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			num1 := DecodeNum(x1)
			num2 := DecodeNum(x2)
			*s = append(*s, EncodeNum(num1*num2))
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpBoolAnd(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x1)&DecodeNum(x2) != 0 {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpBoolOr(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x1)|DecodeNum(x2) != 0 {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpNumEqual(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x1) == DecodeNum(x2) {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpNumEqualVerify(s *[]any) bool {
	return OpNumEqual(s) && OpVerify(s)
}

func OpNumNotEqual(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x1) != DecodeNum(x2) {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpLessThan(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x2) < DecodeNum(x1) {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpGreaterThan(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x2) > DecodeNum(x1) {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpLessThanOrEqual(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x2) <= DecodeNum(x1) {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpGreaterThanOrEqual(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x2) >= DecodeNum(x1) {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpMin(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x1) < DecodeNum(x2) {
				*s = append(*s, x1)
			} else {
				*s = append(*s, x2)
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpMax(s *[]any) bool {
	if len(*s) < 2 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			if DecodeNum(x1) > DecodeNum(x2) {
				*s = append(*s, x1)
			} else {
				*s = append(*s, x2)
			}
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpWithin(s *[]any) bool {
	if len(*s) < 3 {
		return false
	}

	x1 := (*s)[len(*s)-1]
	x2 := (*s)[len(*s)-2]
	x3 := (*s)[len(*s)-3]
	*s = (*s)[:len(*s)-3]

	switch x1 := x1.(type) {
	case []byte:
		switch x2 := x2.(type) {
		case []byte:
			switch x3 := x3.(type) {
			case []byte:
				if DecodeNum(x3) <= DecodeNum(x2) && DecodeNum(x2) < DecodeNum(x1) {
					*s = append(*s, EncodeNum(1))
				} else {
					*s = append(*s, EncodeNum(0))
				}
				return true
			default:
				return false
			}
		default:
			return false
		}
	default:
		return false
	}
}

func OpRipemd160(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		hash := utils.Ripemd160(element)
		*s = append(*s, hash)
		return true
	default:
		return false
	}
}

func OpSha1(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		hash := utils.Sha1(element)
		*s = append(*s, hash)
		return true
	default:
		return false
	}
}

func OpSha256(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		hash := utils.Sha256(element)
		*s = append(*s, hash)
		return true
	default:
		return false
	}
}

func OpHash160(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		hash := utils.Hash160(element)
		*s = append(*s, hash)
		return true
	default:
		return false
	}
}

func OpHash256(s *[]any) bool {
	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	switch element := element.(type) {
	case []byte:
		hash := utils.Hash256(element)
		*s = append(*s, hash)
		return true
	default:
		return false
	}
}

func OpCheckSig(s *[]any, z []byte) bool {
	if len(*s) < 2 {
		return false
	}

	pubKey := (*s)[len(*s)-1]
	derSig := (*s)[len(*s)-2]
	*s = (*s)[:len(*s)-2]

	switch pubKey := pubKey.(type) {
	case []byte:
		switch derSig := derSig.(type) {
		case []byte:

			point, err := ecc.ParsePoint(pubKey)
			if err != nil {
				return false
			}

			sig, err := ecc.ParseSignature(derSig)
			if err != nil {
				return false
			}

			ok, err := point.Verify(z, sig)
			if err != nil {
				return false
			}

			if ok {
				*s = append(*s, EncodeNum(1))
			} else {
				*s = append(*s, EncodeNum(0))
			}

			return true
		default:
			return false
		}
	default:
		return false
	}
}

func OpCheckSigVerify(s *[]any, z []byte) bool {
	return OpCheckSig(s, z) && OpVerify(s)
}

func OpCheckMultiSig(s *[]any, z []byte) bool {
	if len(*s) < 1 {
		return false
	}

	encN, ok := (*s)[len(*s)-1].([]byte) // 인코딩된 n
	if !ok {
		return false
	}
	*s = (*s)[:len(*s)-1]

	n := DecodeNum(encN) // n

	if len(*s) < n+1 {
		return false
	}

	pubKeys := make([][]byte, n) // pubkeys

	for i := 0; i < n; i++ {
		pubKey, ok := (*s)[len(*s)-1].([]byte)
		if !ok {
			return false
		}
		*s = (*s)[:len(*s)-1]
		pubKeys[i] = pubKey
	}

	encM, ok := (*s)[len(*s)-1].([]byte) // 인코딩된 m
	if !ok {
		return false
	}
	*s = (*s)[:len(*s)-1]

	m := DecodeNum(encM) // m

	if len(*s) < m+1 {
		return false
	}

	derSigs := make([][]byte, m) // der sigs

	for i := 0; i < m; i++ {
		derSig, ok := (*s)[len(*s)-1].([]byte)
		if !ok {
			return false
		}

		*s = (*s)[:len(*s)-1]
		derSigs[i] = derSig[:len(derSig)-1] // remove the sighash type
	}

	*s = (*s)[:len(*s)-1] // pop off the 0

	points := make([]ecc.Point, n)   // points
	sigs := make([]ecc.Signature, m) // sigs

	for i := 0; i < n; i++ {
		point, err := ecc.ParsePoint(pubKeys[i])
		if err != nil {
			log.Println("line 1538:", err)
			return false
		}
		points[i] = point
	}

	for i := 0; i < m; i++ {
		sig, err := ecc.ParseSignature(derSigs[i])
		if err != nil {
			log.Println("line 1547:", err)
			return false
		}
		sigs[i] = sig
	}

	// check that all the signatures are valid
	for _, sig := range sigs {
		for len(points) > 0 {
			point := points[0]
			points = points[1:]

			ok, err := point.Verify(z, sig)
			if err != nil {
				log.Println("line 1561:", err)
				return false
			}

			if ok {
				break
			}
		}
	}

	*s = append(*s, EncodeNum(1))

	return true
}

func OpCheckMultiSigVerify(s *[]any, z []byte) bool {
	return OpCheckMultiSig(s, z) && OpVerify(s)
}

func OpCheckLockTimeVerify(s *[]any, locktime, seqNo int) bool {
	if seqNo == 0xffffffff {
		return false
	}

	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]

	switch element := element.(type) {
	case []byte:
		decoded := DecodeNum(element)
		if decoded < 0 {
			return false
		}

		if decoded < 500000000 && locktime < 500000000 {
			return false
		}

		if locktime < decoded {
			return false
		}
		return true
	default:
		return false
	}
}

func OpCheckSequenceVerify(s *[]any, version, seqNo int) bool {
	if seqNo&0x80000000 == 0x80000000 {
		return false
	}

	if len(*s) < 1 {
		return false
	}

	element := (*s)[len(*s)-1]

	switch element := element.(type) {
	case []byte:
		decoded := DecodeNum(element)
		if decoded < 0 {
			return false
		}

		if decoded&0x80000000 == 0x80000000 {
			if version < 2 {
				return false
			}

			if seqNo&0x80000000 == 0x80000000 {
				return false
			}

			if decoded&0x7fffffff != seqNo&0x7fffffff {
				return false
			}

			if decoded&0xffff > seqNo&0xffff {
				return false
			}
		}

		return true
	default:
		return false
	}
}
