package interpreter

import "strings"

func Evaluate(node *Node, acc int) (*Node, error) {

	var spaces = strings.Repeat(" ", acc)

	//_log.Debug(spaces, "Evaluating...")

	if node == nil {
		_log.Debug(spaces, "node is nil")
		return &Node{Value: "nil"}, nil

	}

	switch node.Type {
	case Function:
		_log.Debug(spaces, "function: ", node.Value)
		node, err := node.EvalFunc(node, acc)

		if err != nil {
			return nil, err
		}

		return node, nil
	case Literal:
		_log.Debug(spaces, "literal:", node.Value)

		return node, nil
	case List:
		_log.Debug(spaces, "list:", ChildrenToString(node))

		return node, nil

	case Symbol:
		_log.Debug(spaces, "symbol", node.Value)

		node, err := node.EvalFunc(node, acc)

		if err != nil {
			return nil, err
		}

		return node, nil

	default:
		_log.Debug(spaces, "nil")
		return nil, nil
	}
}


