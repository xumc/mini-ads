package engine

type FieldType int

const (
	_ FieldType = iota
	Number
	Bytes
	Bool
)

type DicType int

const (
	_ DicType = iota
	Hash
	Trie
)

type Field struct {
	Name    string
	Type    FieldType
	DicType DicType
}

type Schema []Field
