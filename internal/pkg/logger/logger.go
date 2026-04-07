package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type ErrorType CustomLog

var (
	writer io.Writer = os.Stderr
	DEBUG            = New(writer, "DEBUG: ", log.LstdFlags|log.Lshortfile)
	INFO             = New(writer, "INFO: ", log.LstdFlags)
	ERROR            = ErrorType(New(writer, "ERROR: ", log.LstdFlags|log.Lshortfile))
)

type CustomLog struct {
	Out    io.Writer
	Prefix string
	Flag   int
}

func New(out io.Writer, prefix string, flag int) CustomLog {
	return CustomLog{Out: out, Prefix: prefix, Flag: flag}
}

func (cl *ErrorType) Printf(format string, v ...any) {
	l := log.New(cl.Out, cl.Prefix, cl.Flag)
	err := l.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cl *CustomLog) Printf(format string, v ...any) {
	l := log.New(cl.Out, cl.Prefix, cl.Flag)
	err := l.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cl *CustomLog) Print(v ...any) {
	l := log.New(cl.Out, cl.Prefix, cl.Flag)
	err := l.Output(2, fmt.Sprint(v...))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cl *CustomLog) Println(v ...any) {
	l := log.New(cl.Out, cl.Prefix, cl.Flag)
	err := l.Output(2, fmt.Sprintln(v...))
	if err != nil {
		fmt.Println(err)
		return
	}
}
