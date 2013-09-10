package runner

import (
	"flag"
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/exec"
	"github.com/ElricleNecro/TOD/formatter"
)

type Config struct {

	// The file name for the users and associated commands
	Users *string

	// The file name for the hosts on which to execute commands
	Hosts *string
}

// Execute the main program.
func Run() {

	// read the command line to get files names
	data_config := ReadConfig()

	// get data structure from files
	hosts_config := configuration.ReadHostsYAML(*data_config.Hosts)
	users_config := configuration.ReadUsersYAML(*data_config.Users)

	// convert those data to dispatcher data
	hosts := formatter.HostsToDispatcher(*hosts_config)
	commands := formatter.UsersToDispatcher(*users_config)

	// Dispatch commands on hosts for the first time
	formatter.Dispatcher(
		commands,
		hosts,
		true,
	)

	// now run commands on hosts
	exec.RunCommands(
		hosts,
		len(commands),
	)
}

// A function to get the data from the command line
// and store it correctly in the datastructure.
func ReadConfig() *Config {

	// Create a configuration object
	data_config := &Config{}

	// define flag to use in the command line
	data_config.Users = flag.String(
		"users",
		"",
		"The path to the file where users and associated commands are stored.",
	)
	data_config.Hosts = flag.String(
		"hosts",
		"",
		"The path to the file where hosts are stored.",
	)

	// parse the command line
	flag.Parse()

	// return the data structure
	return data_config

}
