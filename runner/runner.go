package runner

import (
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/exec"
	"github.com/ElricleNecro/TOD/formatter"
)

// Execute the main program.
func Run() {

	// define the variable
	var hosts_config *configuration.Hosts

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

	// Dispatch commands on hosts for the first time
	formatter.Dispatcher(
		commands,
		hosts,
		data_config.HostsMax,
		true,
	)

	// now run commands on hosts
	exec.RunCommands(
		hosts,
		len(commands),
		data_config,
	)
}
