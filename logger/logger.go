package logger

import (
	"fmt"
	"io"
	"log"

	"github.com/fatih/color"
)

const (
	_ = iota
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
)

var (
	logger *log.Logger

	logLevel   = LevelWarning
	levelNames = map[string]int{
		"error":   LevelError,
		"warning": LevelWarning,
		"notice":  LevelNotice,
		"info":    LevelInfo,
		"debug":   LevelDebug,
	}

	logColor bool
)

func Init(out io.Writer, color bool) {
	logger = log.New(out, "", log.LstdFlags|log.Lmicroseconds)
	logColor = color
}

func GetLevelByName(name string) (int, error) {
	level, ok := levelNames[name]
	if !ok {
		return 0, fmt.Errorf("invalid level `%s'", name)
	}

	return level, nil
}

func SetLevel(level int) {
	logLevel = level

	if logLevel < LevelError || logLevel > LevelDebug {
		logLevel = LevelInfo
	}
}

func SetOutput(output io.Writer) {
	log.SetOutput(output)
}

func Log(level int, context, format string, v ...interface{}) {
	var criticity string

	if level > logLevel {
		return
	}

	switch level {
	case LevelError:
		if logColor {
			criticity = color.RedString("ERROR")
		} else {
			criticity = "ERROR"
		}
	case LevelWarning:
		if logColor {
			criticity = color.YellowString("WARNING")
		} else {
			criticity = "WARNING"
		}
	case LevelNotice:
		if logColor {
			criticity = color.MagentaString("NOTICE")
		} else {
			criticity = "NOTICE"
		}
	case LevelInfo:
		if logColor {
			criticity = color.BlueString("INFO")
		} else {
			criticity = "INFO"
		}
	case LevelDebug:
		if logColor {
			criticity = color.CyanString("DEBUG")
		} else {
			criticity = "DEBUG"
		}
	}

	logger.Printf(
		"%s: %s",
		fmt.Sprintf("%s: %s", criticity, context),
		fmt.Sprintf(format, v...),
	)
}

func Error(context, format string, v ...interface{}) {
	Log(LevelError, context, format, v...)
}

func Warning(context, format string, v ...interface{}) {
	Log(LevelWarning, context, format, v...)
}

func Notice(context, format string, v ...interface{}) {
	Log(LevelNotice, context, format, v...)
}

func Info(context, format string, v ...interface{}) {
	Log(LevelInfo, context, format, v...)
}

func Debug(context, format string, v ...interface{}) {
	Log(LevelDebug, context, format, v...)
}
