package interpreter

import (
	"errors"
	"strconv"
	"strings"

	"main/lib"
)

var _log lib.Logger = nil

func InitInterpreter(log lib.Logger) {
	_log = log
}

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

func Parse(program string) error {

	_log.Debug("")

	_log.Debug("Parsing program:", program)

	_log.Debug("Tokenizing...")

	tokens, err := Tokenize(program)

	if err != nil {
		return err
	}

	_log.Debug("Reading Expressions...")
	root, err := ReadHead(tokens)

	if err != nil {
		return err
	}

	_log.Debug("Logging...")
	VisitNodes(root, nil, LoggingVisitor, 0)
	VisitNodes(root, nil, ParentVisitor, 0)

	_log.Debug("Evaluating...")
	val, err := Evaluate(root, 0)

	if err != nil {
		return err
	}

	_log.Debug("val: ", val.Value)

	return nil
}

func Tokenize(program string) (*lib.Queue, error) {

	if len(program) == 0 {
		return nil, errors.New("program is empty")
	}

	program = strings.ReplaceAll(program, "(", " ( ")
	program = strings.ReplaceAll(program, ")", " ) ")

	fields := strings.Fields(program)

	var tokens lib.Queue = lib.Queue{}

	for i := 0; i < len(fields); i++ {
		tokens.Push(fields[i])
	}

	//_log.Debug("tokens:", tokens)

	return &tokens, nil
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

func ReadTail(tokens *lib.Queue) ([]*Node, error) {
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
			child, err = ReadHead(tokens)
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

func ReadHead(tokens *lib.Queue) (*Node, error) {
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
		newNode, err = ReadHead(tokens)

		if err != nil {
			return nil, err
		}

		if newNode == nil {
			//_log.Debug("new node is nil")
			return nil, nil
		}

		tail, err := ReadTail(tokens)

		if err != nil {
			return nil, err
		}

		newNode.Children = append(newNode.Children, tail...)

	} else {
		newNode, err = ReadAtom(head)

		if err != nil {
			return nil, err
		}

		newNode.Children, err = ReadTail(tokens)

		if err != nil {
			return nil, err
		}
	}

	//_log.Debug("returning", newNode.Value, ChildrenToString(newNode))

	return newNode, nil
}

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

type EvaluationFunc func(node *Node, acc int) (*Node, error)

func Add(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)
	total := "0"

	if len(node.Children) == 0 {
		_log.Debug(spaces, "func +")
		return node, nil
	}

	for _, node := range node.Children {
		totalVal, _ := strconv.Atoi(total)

		bNode, err := Evaluate(node, acc+2)

		if err != nil {
			return nil, err
		}

		bVal, err := strconv.Atoi(bNode.Value)

		if err != nil {
			return nil, err
		}

		total = strconv.Itoa(totalVal + bVal)

	}

	_log.Debug(spaces, " val:", total)

	return &Node{
		Value: total,
		Type:  Literal,
	}, nil
}

func Subtract(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)
	total := "0"
	first := true

	if len(node.Children) == 0 {
		return node, nil
	}

	_log.Debug(spaces, "subtracting")

	for _, child := range node.Children {

		childEvaled, err := Evaluate(child, acc+2)

		if err != nil {
			_log.Debug("error evaluating during subtraction")
			return nil, err
		}

		val := childEvaled.Value

		if err != nil {
			_log.Debug("error converting 1 during subtraction", err)
			return nil, err
		}

		bVal, err := strconv.Atoi(val)

		if err != nil {
			_log.Debug("error evaluating 2 during subtraction", err)
			return nil, err
		}

		//todo: optimize all this int<>string conversion!
		totalVal, err := strconv.Atoi(total)

		if err != nil {
			_log.Debug("error converting total during subtraction", err)
			return nil, err
		}

		if first {
			total = strconv.Itoa(totalVal + bVal)
			first = false

		} else {
			total = strconv.Itoa(totalVal - bVal)
		}
	}

	//_log.Debug(spaces, " val:", total)

	return &Node{
		Value: total,
		Type:  Literal,
	}, nil
}

func First(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)

	if node == nil || len(node.Children) == 0 {
		return nil, nil

	}

	head := node.Children[0]

	_log.Debug(spaces, "first'ing")

	return head, nil
}

func Rest(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)

	if node == nil || len(node.Children) == 0 {
		return nil, nil

	}

	var newNode *Node = &Node{
		Value:    "list",
		Type:     List,
		Children: node.Children[1:],
	}

	_log.Debug(spaces, "rest'ing", ChildrenToString(newNode))

	return newNode, nil
}

func Equals(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)
	truth := true

	_log.Debug(spaces, "equalizing")

	if node != nil && len(node.Children) > 1 {

		headEvaled, err := Evaluate(node.Children[0], acc+2)

		if err != nil {
			return nil, err
		}

		headVal := headEvaled.Value

		for i := 1; i < len(node.Children); i++ {
			child := node.Children[i]

			childEvaled, err := Evaluate(child, acc+2)

			if err != nil {
				return nil, err
			}

			childVal := childEvaled.Value

			truth = truth && (headVal == childVal)

			if !truth {
				break
			}
		}
	}

	_log.Debug(spaces, " val:", truth)

	return &Node{
		Value: strconv.FormatBool(truth),
		Type:  Literal,
	}, nil
}

func If(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)

	if node == nil || len(node.Children) < 3 {
		return nil, nil
	}

	truthyNode, err := Evaluate(node.Children[0], acc+2)

	if err != nil {
		return nil, err
	}

	if truthyNode.Type != Literal {
		return nil, nil
	}

	_log.Debug("testing ", truthyNode.Value)

	if strings.ToLower(truthyNode.Value) == "true" {
		_log.Debug(spaces, " val:", "if(true)")
		return Evaluate(node.Children[1], acc+2)
	} else {
		_log.Debug(spaces, " val:", "if(false)")
		return Evaluate(node.Children[2], acc+2)
	}
}

func Concat(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)

	if node == nil || len(node.Children) == 0 {
		return nil, nil
	}

	var children []*Node = []*Node{}

	_log.Debug(spaces, "concatanating", len(node.Children), "children")

	for _, child := range node.Children {
		childEvaled, err := Evaluate(child, acc+2)

		if err != nil {
			return nil, err
		}

		if childEvaled.Type == List {
			children = append(children, childEvaled.Children...)

		} else {
			children = append(children, childEvaled)

		}

	}

	var newNode *Node = &Node{
		Value:    "list",
		Type:     List,
		Children: children,
	}

	return newNode, nil
}

func SymbolEval(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)

	if node == nil {
		return nil, nil
	}

	var parent *Node = node.Parent

	var found *Node = nil

	_log.Debug(spaces, "looking up symbol", node.Value)

	for parent != nil {
		_log.Debug(spaces, " ...", parent.Value)
		foundNode := parent.Scope[node.Value]

		if foundNode != nil {
			_log.Debug(spaces, " found in", parent.Value)

			found = foundNode
			break
		}

		parent = parent.Parent
	}

	if found == nil {
		_log.Debug(spaces, " unknown symbol", node.Value, ", returning as-is")

		return node, nil
	}

	//save children from found node, to reset it after eval'ing
	savedChildren := found.Children

	//take the children and then see if evaluating is intersting?
	found.Children = append(found.Children, node.Children...)

	_log.Debug(spaces, "evaluating", found.Value, "with [", len(found.Children), "] children")
	newNode, err := Evaluate(found, acc+2)

	if err != nil {
		return nil, err
	}

	//restore children from before, just in case it was altered
	found.Children = savedChildren

	return newNode, nil

}

func AddToNodeScope(node *Node, symbolName string, symbolValue *Node) {
	if node.Scope == nil {
		node.Scope = make(map[string]*Node)
	}

	node.Scope[symbolName] = symbolValue
}

func Let(node *Node, acc int) (*Node, error) {
	spaces := strings.Repeat(" ", acc)

	if node == nil || len(node.Children) < 3 {
		return nil, nil
	}

	symbolName := node.Children[0].Value

	symbolValue := node.Children[1]

	AddToNodeScope(node, symbolName, symbolValue)

	_log.Debug(spaces, "defining", symbolName, "=", symbolValue.Type, symbolValue.Value)

	childEvaled, err := Evaluate(node.Children[2], acc+2)

	if err != nil {
		return nil, err
	}

	return childEvaled, nil

}

func Lambda(node *Node, acc int) (*Node, error) {
	//(lambda, list of params, body, application)

	spaces := strings.Repeat(" ", acc)

	if node == nil || len(node.Children) < 3 {
		return nil, nil
	}

	//this is the definition of the lambda
	if len(node.Children) == 2 {
		return node, nil
	}

	//node.Children > 2

	//this is the application of the lambda

	paramList, err := Evaluate(node.Children[0], acc+2) //params - must be a list of symbols

	if err != nil {
		return nil, err
	}

	if paramList.Type != List {
		err := spaces + "passed params are not a list!"
		_log.Debug(err)
		return nil, errors.New(err)
	}

	body := node.Children[1] //body - can be whatever (i think?)

	for i := 2; i < len(node.Children); i++ {

		//evaluate the parameter passed to the lambda
		//todo: figure out if we should eval now or later
		val, err := Evaluate(node.Children[2], acc+2)

		if err != nil {
			return nil, err
		}

		//todo: should probably evaulate this before using it...
		symbol := paramList.Children[i-2]

		if symbol.Type != Symbol {
			err := spaces + "paramter" + symbol.Value + "is not a symbol"
			_log.Debug(err)

			return nil, errors.New(err)
		}

		//now let's each child to each symbol in the list

		symbolName := symbol.Value
		symbolValue := val

		AddToNodeScope(body, symbolName, symbolValue)
	}

	return Evaluate(body, acc+2)
}

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
