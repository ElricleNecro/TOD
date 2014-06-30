package tools

import (
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/ElricleNecro/TOD/formatter"
)

type config interface {
	GetLogCommand() string
}

// expanduser
func Expanduser(path string) string {

	var home string

	// get the current user
	if usr, err := user.Current(); err == nil {

		// get the home directory of the current user
		home = usr.HomeDir
	} else {

		// an error occurred, fallback to the home variable
		home = os.ExpandEnv("$HOME")
	}

	// check the path in input
	if len(path) < 1 {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"The length of the path isn't sufficient!",
		)
	}

	// replace the tilde by home
	if path[:1] == "~" {
		path = strings.Replace(path, "~", home, 1)
	}
	return path
}

// Functions to write the results of command into a log file in a given
// directory specified in argument of the program with the option -log_command.
func WriteLogCommand(
	output string,
	config config,
	hostname string,
	command string,
	number int,
) {

	// Create the file name with hostname and the number of the command
	filename := config.GetLogCommand() + "/" +
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
