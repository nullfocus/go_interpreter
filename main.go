package main

import (
	"io/ioutil"
	"main/interpreter"
	"main/lib"
	"os"
)

var _log lib.Logger = nil

func main() {
	_log = lib.InitLogging()
	interpreter.InitInterpreter(_log)

	filename := os.Args[1]

	contents, _ := ioutil.ReadFile(filename)

	Interpret(string(contents))

	//testInterpreter(_log)
}

func testInterpreter(log lib.Logger) {
	Interpret("(+ 1 2)")

	Interpret("(- 4 3)")

	Interpret("(+ 5 (- 6 13))")

	Interpret("(((")

	Interpret(")))")

	Interpret("(((+) 4)  3)")

	Interpret("(+ (+ 1 2) (- 3 2))")

	Interpret("(first 1 2 3 4 5)")

	Interpret("(rest 1 2 3 4 5)")

	Interpret("(eq 1 1)")

	Interpret("(eq 1 2)")

	Interpret("(if (true) 1 2)")

	Interpret("(if (false) 1 2)")

	Interpret("(if (eq (- -6 7) (+ -14 1)) 55 66)")

	Interpret("(concat 1 2 3 4 5 6)")

	Interpret("(concat 1 2 (rest 2 3 4 5) 6)")

	Interpret("(let x 4 (+ x 1))")

	Interpret("(let x 4 (let func + (func x 1))")

	Interpret("(let x + (x 1 1))")

	Interpret("((lambda (concat x) (+ x 1)) 1)")

	Interpret("(+ x 1 )")

	Interpret(`
	 (let subtract-one 
		(lambda 
			(concat x) 
			(- x 1)) 
		(+ 1 (subtract-one 3)))`)

	Interpret(`
		(let x 
			(if (eq 1 0) + -) 
			(x 1 1))`)

	Interpret(`((if (eq 1 0) + -)  1 0)`)

	Interpret(`
	(let f
		(lambda 
			(concat x)
			(+ x 1))
		(+ (f 2) (f 3)))`)
}

func Interpret(program string) {
	err := interpreter.Interpret(program)

	if err != nil {
		_log.Debug(err.Error())
	}
}
