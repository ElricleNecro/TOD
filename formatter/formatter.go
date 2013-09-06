package formatter

import (
	"fmt"
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

	// Store here if the host is connected or not
	IsConnected bool
}

// This function dispatches the commands on some hosts
func Dispatcher(
	commands []*Command,
	hosts []*Host,
) {

	// the pointer to the host structure
	var host int

	// Determine the number of hosts available in theory
	nhosts := CountConnectedHosts(hosts)

	// check there is at least one host connected
	if nhosts == 0 {
		fmt.Println("There is no hosts available to do the job !")
	}

	// The same for the number of commands to execute
	ncomm := len(commands)

	// Compute the number of commands per hosts to execute
	NCH := math.Ceil(float64(ncomm) / float64(nhosts))

	// init by pointing the current host to the first one
	host = -1

	// loop over commands and affect them to hosts
	for i, command := range commands {

		// check if we need to pass to an other host with our repartition
		if math.Mod(float64(i), NCH) == 0 {

			host++
		}

		// if the host isn't connected
		if !hosts[host].IsConnected {

			// pass to the next host
			host++
		}

		// append to the list of commands
		hosts[host].Commands = append(hosts[host].Commands, command)
	}

}

// This function counts the number of connected hosts given in
// a slice of those objects.
func CountConnectedHosts(hosts []*Host) int {

	var counter int

	// init the counter
	counter = 0

	// loop over host and increment a counter when connected
	for _, host := range hosts {

		// check the field for connected
		if host.IsConnected {
			counter++
		}
	}

	// return the number of connected machines
	return counter
}
