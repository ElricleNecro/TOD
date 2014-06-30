package main

import (
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/dispatcher"
)

// Main function :
//	Will read and parse the program argument, read the sql databases, verify and select machines, launch the loops.
//	Need to now if wait to all process to end or if it will detach the subprocess and quit?

func main() {
	// run the program
	Run()
}

// Execute the main program.
func Run() {

	hosts_config := new(configuration.Hosts)

	// read the command line to get files names
	data_config := configuration.ReadConfig()

	// get data structure from files
	if !data_config.Stdin {
		hosts_config = configuration.ReadHostsYAML(data_config.Hosts)
	} else {
		hosts_config = configuration.ReadHostsStdin(data_config)
	}
	users_config := configuration.ReadUsersYAML(data_config.Users)

	// convert those data to dispatcher data
	hosts := configuration.HostsToDispatcher(*hosts_config)
	commands := configuration.UsersToDispatcher(*users_config)

	// create a new dispatcher
	dispatcher := dispatcher.New(data_config, hosts)

	// Dispatch commands on hosts for the first time
	hosts.Dispatcher(
		commands,
		data_config.HostsMax,
		true,
	)

	// now run commands on hosts
	dispatcher.RunCommands(len(commands))
}

//vim: spelllang=en
