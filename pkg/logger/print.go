package logger

import (
	"fmt"
	"io"
)

func writeAndReturnTo(msg string, out io.Writer) {
	// Voluntarily ignoring return values.
	fmt.Fprintf(out, "%s\n", msg)
}

func writeColoredAndSeparateTo(msg string, color int, out io.Writer) {
	// Voluntarily ignoring return values.
	fmt.Fprintf(out, "\033[1;%dm%s\033[0m ", color, msg)
}
