package interpreter

import "main/lib"

var _log lib.Logger = nil

func InitInterpreter(log lib.Logger) {
	_log = log
}

func Interpret(program string) error {

	_log.Debug("")

	_log.Debug("Interpreting program:", program)

	_log.Debug("Scanning...")

	tokens, err := Scan(program)

	if err != nil {
		return err
	}

	_log.Debug("Parsing...")
	root, err := Parse(tokens)

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
