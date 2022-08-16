package interpreter

import "strings"

type NodeVisitor func(*Node, *Node, int)

func LoggingVisitor(node *Node, parent *Node, acc int) {
	var spaces = strings.Repeat(" ", acc)

	_log.Debug(spaces, node.Type, ":", node.Value)

}

func ParentVisitor(node *Node, parent *Node, acc int) {
	node.Parent = parent
}

func VisitNodes(node *Node, parent *Node, visitor NodeVisitor, acc int) {

	if node != nil {

		visitor(node, parent, acc)

		if len(node.Children) > 0 {
			acc = acc + 2

			for _, child := range node.Children {
				VisitNodes(child, node, visitor, acc)
			}
		}
	}
}

func ChildrenToString(node *Node) string {
	if node == nil {
		return ""
	}

	acc := ""

	if len(node.Children) > 0 {
		for _, node := range node.Children {
			if len(acc) > 0 {
				acc += ", "
			}

			acc += node.Value
		}
	}

	acc = "[" + acc + "]"

	return acc
}
