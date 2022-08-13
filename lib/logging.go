package lib

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger interface {
	Debug(v ...any)
	Write(b []byte) (n int, err error)
}

type simpleLog struct {
	log    *log.Logger
	writer io.Writer
}

//implement Debug interface
func (sl simpleLog) Debug(v ...any) {
	sl.log.Print(v...)
	fmt.Println(v...)
}

//implement Writer interface
func (sl simpleLog) Write(b []byte) (n int, err error) {
	return sl.writer.Write(b)
}

func InitLogging() Logger {
	logFile, err := os.OpenFile("scratch.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		log.Fatalln(err)
	}

	var logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds)

	logger.Print("logging initalized")

	sl := simpleLog{log: logger, writer: logFile}

	return &sl
}
