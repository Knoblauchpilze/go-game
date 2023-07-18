package logger

import (
	"fmt"
)

func writeAndReturn(msg string) string {
	return fmt.Sprintf("%s\n", msg)
}

func writeColoredAndSeparate(msg string, color int) string {
	return fmt.Sprintf("\033[1;%dm%s\033[0m ", color, msg)
}
