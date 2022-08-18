package interpreter

import (
	"errors"
	"main/lib"
	"strings"
)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func isParens(r rune) bool {
	return r == '(' || r == ')'
}

func Scan(program string) (*lib.Queue, error) {

	if len(program) == 0 {
		return nil, errors.New("program is empty")
	}

	program = strings.ReplaceAll(program, "(", " ( ")
	program = strings.ReplaceAll(program, ")", " ) ")

	_log.Debug("program:", program)

	var tokens lib.Queue = lib.Queue{}

	for i := 0; i < len(program); i++ {
		cur_char := rune(program[i])

		//_log.Debug(string(cur_char))

		switch cur_char {
		//whitespace
		case ' ', '	', '\r', '\n':
			continue

		//parens
		case '(', ')':
			tokens.Push(string(cur_char))

		//quotes
		case '"':
			tokens.Push(string(cur_char))

			for {
				i++
				if i >= len(program) {
					return nil, errors.New("unexpected EOF while reading quotes")
				}

				cur_char = rune(program[i])

				tokens.Push(string(cur_char))

				if cur_char == '"' {
					break
				}
			}

		//string
		default:
			var str string

			for {
				str += string(cur_char)

				i++

				cur_char = rune(program[i])

				if i >= len(program) || isWhitespace(cur_char) {
					break
				}
			}

			_log.Debug("str", str)

			tokens.Push(str)
		}
	}

	_log.Debug("tokens:", tokens)

	return &tokens, nil
}
