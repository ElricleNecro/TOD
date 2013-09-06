package formatter

import (
	"math"
)

// The structure containing information for the connection as a given
// user to a host.
type User struct {

	// The user name
	Name string

	// An identity to use as reference for the user if necessary
	Identity int

	// The password to the host if necessary
	Password string
}

// The structure used to pass commands to host. To one host and the
// corresponding user, we put here the list of commands to execute.
type Command struct {

	// The command to execute
	Command string

	// The user which needs to execute the command
	User *User
}

// This structure contains the list of commands to execute on the
// corresponding host. So we allow to easily dispatch those commands
// on the host.
type Host struct {

	// the host
	Hostname string

	// The port on which to connect
	Port int

	// The protocol to use in the host for the connection
	Protocol string

	// The list of commands to execute on the host
	Commands []*Command
}

// This function dispatches the commands on some hosts
func Dispatcher(
	commands []*Command,
	hosts []*Host,
) {

	// the pointer to the host structure
	var host int

	// Determine the number of hosts available in theory
	nhosts := len(hosts)

	// The same for the number of commands to execute
	ncomm := len(commands)

	// Compute the number of commands per hosts to execute
	NCH := math.Ceil(float64(ncomm) / float64(nhosts))

	// init by pointing the current host to the first one
	host = -1

	// loop over commands and affect them to hosts
	for i, command := range commands {

		// check if we need to pass to an other host with our repartition
		if math.Mod(float64(i), float64(NCH)) == 0 {
			host++
		}

		// append to the list of commands
		hosts[host].Commands = append(hosts[host].Commands, command)
	}

}
