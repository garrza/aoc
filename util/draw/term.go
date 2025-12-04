package draw

import (
	"fmt"
	"os"
)

func Clear() {
	fmt.Fprintf(os.Stdout, "\033[2J")
}

func Move(x, y uint) {
	fmt.Fprintf(os.Stdout, "\033[%d;%dH", y+1, x+1)
}

func MoveRel(x, y int) {
	if y < 0 {
		fmt.Fprintf(os.Stdout, "\033[%dA", -y)
	}
	if y > 0 {
		fmt.Fprintf(os.Stdout, "\033[%dB", y)
	}
	if x > 0 {
		fmt.Fprintf(os.Stdout, "\033[%dC", x)
	}
	if x < 0 {
		fmt.Fprintf(os.Stdout, "\033[%dD", -x)
	}
}

func Write(char string) {
	os.Stdout.WriteString(char)
}
