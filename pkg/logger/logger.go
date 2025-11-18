package logger

import (
	"log"
	"os"
)

func New() *log.Logger {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	return l
}
