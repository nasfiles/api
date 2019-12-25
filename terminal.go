package api

import (
	"github.com/fatih/color"
)

//LineSeparator prints a whole line of a given character
func LineSeparator(ch string, color *color.Color, width int) {
	for i := 1; i <= width; i++ {
		color.Printf(ch)
	}
}

//YesNoColored prints Yes in green and No in red
func YesNoColored(test bool) {
	if test {
		color.HiGreen("Yes")
	} else {
		color.HiRed("No")
	}
}
