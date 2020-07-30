package terminal

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

// Size returns the terminal characters per line limit
func Size() int {
	// get terminal size
	var width int = 120

	x, err := terminal.Width()
	if err != nil {
		log.Println("Couldn't get the terminal width")
	}
	width = int(x)

	return width
}

// LineSeparator prints a whole line of a given character
func LineSeparator(ch string, color *color.Color, width int) {
	for i := 1; i <= width; i++ {
		color.Printf(ch)
	}
	fmt.Println()
}

// YesNoColored prints Yes in green and No in red
func YesNoColored(test bool) {
	if test {
		color.HiGreen("Yes")
	} else {
		color.HiRed("No")
	}
}
