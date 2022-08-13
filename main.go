package main

import (
	"main/lib"
	"main/interpreter"
)

var _log lib.Logger = nil

func main() {
	_log = lib.InitLogging()

	testInterpreter(_log)
}

func testInterpreter(log lib.Logger) {
	interpreter.InitInterpreter(log)

	Parse("(+ 1 2)")

	Parse("(- 4 3)")

	Parse("(+ 5 (- 6 13))")

	Parse("(((")

	Parse(")))")

	Parse("(((+) 4)  3)")

	Parse("(+ (+ 1 2) (- 3 2))")

	Parse("(first 1 2 3 4 5)")

	Parse("(rest 1 2 3 4 5)")

	Parse("(eq 1 1)")

	Parse("(eq 1 2)")

	Parse("(if (true) 1 2)")

	Parse("(if (false) 1 2)")

	Parse("(if (eq (- -6 7) (+ -14 1)) 55 66)")

	Parse("(concat 1 2 3 4 5 6)")

	Parse("(concat 1 2 (rest 2 3 4 5) 6)")

	Parse("(let x 4 (+ x 1))")

	Parse("(let x 4 (let func + (func x 1))")

	Parse("(let x + (x 1 1))")

	Parse("((lambda (concat x) (+ x 1)) 1)")

	Parse("(+ x 1 )")

	Parse(`
	 (let subtract-one 
		(lambda 
			(concat x) 
			(- x 1)) 
		(+ 1 (subtract-one 3)))`)

	Parse(`
		(let x 
			(if (eq 1 0) + -) 
			(x 1 1))`)

	Parse(`((if (eq 1 0) + -)  1 0)`)

	Parse(`
	(let f
		(lambda 
			(concat x)
			(+ x 1))
		(+ (f 2) (f 3)))`)
}

func Parse(program string) {
	err := interpreter.Parse(program)

	if err != nil {
		_log.Debug(err.Error())
	}
}


