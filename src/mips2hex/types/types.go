package types

type ItemKind int

const (
	Instr ItemKind = iota
	Word
)

type Item struct {
	Kind     ItemKind
	Raw      string   // 清理过注释和空白的原始文本
	Tokens   []string // tokenized
	LineNo   int
	OrigLine string
	Size     uint32
}

type InstrType int

const (
	RType InstrType = iota
	IType
	JType
	Special
)

type Instruction struct {
	Type   InstrType
	Opcode uint32
	Funct  uint32
}
