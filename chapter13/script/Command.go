package script

import "fmt"

type Command struct {
	IsOpCode bool
	Code     OpCode
	Elem     []byte
}

func NewOpCode(code OpCode) Command {
	return Command{true, code, nil}
}

func NewElem(elem []byte) Command {
	return Command{false, 0, elem}
}

func (c Command) String() string {
	if c.IsOpCode {
		return c.Code.String()
	}
	return fmt.Sprintf("%x", c.Elem)
}
