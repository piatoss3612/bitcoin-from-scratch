package script

import "chapter06/utils"

type OpCode int

type OpCodeFunc func(*Script) bool

const (
	OpCode0           OpCode = 0x00
	OpCode1           OpCode = 0x51
	OpCode16          OpCode = 0x60
	OpCodeDup         OpCode = 0x76
	OpCodeEqualVerify OpCode = 0x88
	OpCodeAdd         OpCode = 0x93
	OpCodeHash160     OpCode = 0xa9
	OpCodeCheckSig    OpCode = 0xac
)

func (o OpCode) Valid() bool {
	return o == OpCode0 || (o >= OpCode1 && o <= OpCode16) || o == OpCodeDup ||
		o == OpCodeAdd || o == OpCodeEqualVerify || o == OpCodeHash160 || o == OpCodeCheckSig
}

func (o OpCode) String() string {
	switch o {
	case OpCode0:
		return "OP_0"
	case OpCode1:
		return "OP_1"
	case OpCode16:
		return "OP_16"
	case OpCodeDup:
		return "OP_DUP"
	case OpCodeEqualVerify:
		return "OP_EQUALVERIFY"
	case OpCodeAdd:
		return "OP_ADD"
	case OpCodeHash160:
		return "OP_HASH160"
	case OpCodeCheckSig:
		return "OP_CHECKSIG"
	default:
		return ""
	}
}

func OpDup(s *Script) bool {
	if len(s.Cmds) < 1 {
		return false
	}

	cmd := s.Cmds[len(s.Cmds)-1]
	s.Cmds = append(s.Cmds, cmd)
	return true
}

func OpHash256(s *Script) bool {
	if len(s.Cmds) < 1 {
		return false
	}

	element := s.Cmds[len(s.Cmds)-1]

	switch elType := element.(type) {
	case []byte:
		hash := utils.Hash256(elType)
		s.Cmds[len(s.Cmds)-1] = hash
		return true
	default:
		return false
	}
}
