package configuration

import (
	"io/ioutil"

	"github.com/ElricleNecro/TOD/commands"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/user"

	"launchpad.net/goyaml"
)

// A type defining the user structure in the YAML file which we need
// in argument.
type MyUsers struct {
	Key      string   "private_key"
	Commands []string "commands"
}

// The map type passed to the read YAML procedure to get multiple users
// informations if we need.
type Users map[string]MyUsers

// A function to read an users file in the YAML format and returns
// a dictionary in the same format as the structured file.
func ReadUsersYAML(
	filename string,
) *Users {

	// Start by reading the whole file in byte
	data, _ := ioutil.ReadFile(filename)

	// Create the variable handling the type of the user file
	t := &Users{}

	// Now read in the YAML file the structure of the file into
	// the structured dictionary
	err := goyaml.Unmarshal(
		data,
		t,
	)

	// Check error when reading the file
	if err != nil {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"The file "+filename+" can't be read for accessing"+
				"the YAML structure!\n"+
				"Reason is: "+err.Error(),
		)
		return nil
	}

	// return the structured file and data
	return t

}

// This function converts users structure from the configuration file into
// the structure used by the dispatcher.
func UsersToDispatcher(users Users) []commands.Command {

	// init the command slice
	cmds := make([]commands.Command, 0)

	// loop over users
	for username, fields := range users {

		// create an user structure
		user := &user.User{
			Username:   username,
			PrivateKey: fields.Key,
		}

		// loop over commands
		for _, command := range fields.Commands {
			cmds = append(
				cmds,
				commands.Command{
					Command: command,
					User:    user,
				},
			)
		}
	}

	// returns the list of commands
	return cmds
}
