package formatter

import (
	"fmt"

	color "github.com/daviddengcn/go-colortext"
)

// A wrapping type for colors
type Color color.Color

// The types of colors
const (
	None    = color.None
	Black   = color.Black
	Red     = color.Red
	Green   = color.Green
	Yellow  = color.Yellow
	Blue    = color.Blue
	Magenta = color.Magenta
	Cyan    = color.Cyan
	White   = color.White
)

// A simple function to print message with colors with line return.
func ColoredPrintln(thecolor color.Color, bold bool, values ...interface{}) {

	// change the color of the terminal
	color.ChangeColor(
		thecolor,
		bold,
		None,
		false,
	)

	// print
	fmt.Println(values...)

	// reset the color
	color.ResetColor()
}

// A simple function to print message with colors.
func ColoredPrint(thecolor color.Color, bold bool, values ...interface{}) {

	// change the color of the terminal
	color.ChangeColor(
		thecolor,
		bold,
		None,
		false,
	)

	// print
	fmt.Print(values...)

	// reset the color
	color.ResetColor()
}
