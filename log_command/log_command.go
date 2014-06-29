package log_command

import (
	"os"
	"strconv"

	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/formatter"
)

// Functions to write the results of command into a log file in a given
// directory specified in argument of the program with the option -log_command.
func WriteLogCommand(
	output string,
	config *configuration.Config,
	hostname string,
	command string,
	number int,
) {

	// Create the file name with hostname and the number of the command
	filename := config.LogCommand + "/" +
		hostname + strconv.Itoa(number) + ".log"

	// content of the file
	content := command + "\n" + output

	// open the file
	f, err := os.Create(filename)

	// defer the close
	defer func() {
		if err := f.Close(); err != nil {
			formatter.ColoredPrintln(
				formatter.Red,
				false,
				"The file "+filename+" can't be closed!\n"+
					"Reason is: "+err.Error(),
			)
		}
	}()

	// error
	if err != nil {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"The file "+filename+" can't be open for logging!\n"+
				"Reason is: "+err.Error(),
		)
	}

	// write the content in the file
	f.WriteString(content)

}
