package log

import "fmt"

func format(level string, msg string, args ...any) string {
	template := fmt.Sprintf("%s: %s\n", level, msg)
	return fmt.Sprintf(template, args...)
}

func Log(level string, msg string, args ...any) {
	formatted := format(level, msg, args...)
	fmt.Print(formatted)
}

func Debug(msg string, args ...any) {
	Log("DEBUG", msg, args...)
}

func Info(msg string, args ...any) {
	Log("INFO", msg, args...)
}

func Warn(msg string, args ...any) {
	Log("WARN", msg, args...)
}

func Error(msg string, args ...any) {
	Log("ERROR", msg, args...)
}
