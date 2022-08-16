package interpreter

import (
	"errors"
	"main/lib"
	"strings"
)

func Scan(program string) (*lib.Queue, error) {

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
