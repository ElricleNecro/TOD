package log_command

import (
	"github.com/ElricleNecro/TOD/configuration"
	"io/ioutil"
	"strconv"
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
	filename := hostname + strconv.Itoa(number) + ".log"

	// content of the file
	content := command + "\n" + output

	// open the file
	err := ioutil.WriteFile(
		filename,
		[]byte(content),
		0666,
	)

	// error
	if err != nil {
		panic(err)
	}
}
