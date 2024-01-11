package script

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
