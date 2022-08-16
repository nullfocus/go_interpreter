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
