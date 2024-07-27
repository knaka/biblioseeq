package log

import (
	"io"
	stdlog "log"
)

var logger = stdlog.New(io.Discard, prefix, stdlog.LstdFlags)

func Println(v ...interface{}) {
	logger.Println(v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}
