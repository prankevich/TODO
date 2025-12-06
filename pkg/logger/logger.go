package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func New() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return logger

}
