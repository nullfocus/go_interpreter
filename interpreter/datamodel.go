package interpreter

type NodeType string

const (
	Literal  NodeType = "literal"
	Function NodeType = "function"
	List     NodeType = "list"
	Symbol   NodeType = "symbol"
)

//ast node
type Node struct {
	Value    string
	Type     NodeType
	Children []*Node
	Parent   *Node

	EvalFunc EvaluationFunc
	Scope    map[string]*Node
}


// --- new approach ---

type ExpressionType int

const (
	CellType  ExpressionType = 0
	ValueType ExpressionType = 1
)

type Expression interface {
	Type() ExpressionType
}

type Cell struct {
	Head Expression
	Tail Expression
}

func (c Cell) Type() ExpressionType {
	return CellType
}

type Value struct {
}

func (c Value) Type() ExpressionType {
	return ValueType
}
