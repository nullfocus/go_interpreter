package interpreter

import (
	"strconv"
	"strings"

	"main/lib"
)

func Parse(tokens *lib.Queue) (*Node, error) {
	return ParseHead(tokens)
}

func ParseHead(tokens *lib.Queue) (*Node, error) {
	//_log.Debug("ReadHead: ", tokens)

	head, empty := tokens.Pop()

	if empty {
		return nil, nil
	}

	//_log.Debug("  head: ", head)

	if head == ")" {
		return nil, nil
	}

	var newNode *Node = nil
	var err error = nil

	if head == ")" {
		return nil, nil
	} else if head == "(" {
		newNode, err = ParseHead(tokens)

		if err != nil {
			return nil, err
		}

		if newNode == nil {
			//_log.Debug("new node is nil")
			return nil, nil
		}

		tail, err := ParseTail(tokens)

		if err != nil {
			return nil, err
		}

		newNode.Children = append(newNode.Children, tail...)

	} else {
		newNode, err = ReadAtom(head)

		if err != nil {
			return nil, err
		}

		newNode.Children, err = ParseTail(tokens)

		if err != nil {
			return nil, err
		}
	}

	//_log.Debug("returning", newNode.Value, ChildrenToString(newNode))

	return newNode, nil
}

func ParseTail(tokens *lib.Queue) ([]*Node, error) {
	//_log.Debug("ReadTail: ", tokens)

	var children []*Node = []*Node{}

	for {
		head, empty := tokens.Pop()

		if empty {
			return children, nil
		}

		//_log.Debug("  head of tail:", head)

		if head == ")" {
			return children, nil
		}

		var child *Node = nil
		var err error = nil

		if head == "(" {
			child, err = ParseHead(tokens)
		} else {
			child, err = ReadAtom(head)
		}

		if err != nil {
			return nil, err
		}

		if child != nil {
			//_log.Debug("  appending", child.Type, ":", child.Value)
			children = append(children, child)
		}
	}
}

//https://stackoverflow.com/questions/45686163/how-to-write-isnumeric-function-in-golang
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isBoolean(s string) bool {
	val := strings.ToLower(s)
	return val == "true" || val == "false"
}

func ReadAtom(token string) (*Node, error) {
	var nodeType NodeType
	var evalFunc EvaluationFunc

	if isNumeric(token) || isBoolean(token) {
		nodeType = Literal

	} else {
		switch strings.ToLower(token) {

		case "first":
			nodeType = Function
			evalFunc = First

		case "rest":
			nodeType = Function
			evalFunc = Rest

		case "+":
			nodeType = Function
			evalFunc = Add
		case "-":
			nodeType = Function
			evalFunc = Subtract

		case "eq":
			nodeType = Function
			evalFunc = Equals

		case "if":
			nodeType = Function
			evalFunc = If

		case "concat":
			nodeType = Function
			evalFunc = Concat

		case "let":
			nodeType = Function
			evalFunc = Let

		case "lambda":
			nodeType = Function
			evalFunc = Lambda

		default:
			nodeType = Symbol
			evalFunc = SymbolEval

		}
	}

	//_log.Debug("ReadAtom:", token, "=>", nodeType)

	newNode := &Node{
		Value:    token,
		Type:     nodeType,
		Children: nil,
		EvalFunc: evalFunc,
	}

	return newNode, nil
}
