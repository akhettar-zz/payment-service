package logger

import (
	"log"
	"os"
)

var (
	// Trace the trace logger
	Trace *log.Logger

	// Info the info logger
	Info *log.Logger

	// Warning the warning logger
	Warning *log.Logger

	// Error the error logger
	Error *log.Logger
)

func init() {
	Trace = log.New(os.Stdout, "TRACE:", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
