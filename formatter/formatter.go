package formatter

import (
	"fmt"
	"github.com/ElricleNecro/TOD/configuration"
	color "github.com/daviddengcn/go-colortext"
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

	// The command number being executed
	CommandNumber int

	// Store here if the host is connected or not
	IsConnected bool

	// Channel on which to wait for new job
	Waiter *(chan int)
}

// This function dispatches the commands on some hosts
func Dispatcher(
	commands []*Command,
	hosts []*Host,
	first bool,
) {

	// the pointer to the host structure
	var host int

	// Determine the number of hosts available in theory
	nhosts := CountConnectedHosts(hosts)

	// check there is at least one host connected
	if nhosts == 0 {
		color.ChangeColor(color.Red, true, color.None, false)
		fmt.Println("There is no hosts available to do the job !")
	}

	// The same for the number of commands to execute
	ncomm := len(commands)

	// Compute the number of commands per hosts to execute
	NCH := math.Ceil(float64(ncomm) / float64(nhosts))

	// init by pointing the current host to the first one
	host = -1

	// store here the list of selected hosts
	myhost := make([]int, 0)

	// loop over commands and affect them to hosts
	for i, command := range commands {

		// check if we need to pass to an other host with our repartition
		if math.Mod(float64(i), NCH) == 0 {

			host++

			// if the host isn't connected
			for !hosts[host].IsConnected {

				// pass to the next host
				host++
			}

			// add the host to list
			myhost = append(myhost, host)
		}

		// append to the list of commands
		hosts[host].Commands = append(hosts[host].Commands, command)
	}

	// message to say that the host has more jobs
	if !first {
		color.ChangeColor(color.None, true, color.None, false)
		for _, host := range myhost {
			fmt.Println("Send more jobs signal to " + hosts[host].Hostname)
			fmt.Println(hosts[host].Waiter)
			*(hosts[host].Waiter) <- 1
		}
	}

}

// This function counts the number of connected hosts given in
// a slice of those objects.
func CountConnectedHosts(hosts []*Host) int {

	// init the counter
	counter := 0

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

// This function counts the number of remaining commands
func CountCommands(hosts []*Host) int {

	// counter
	counter := 0

	// loop over hosts
	for _, host := range hosts {

		// don't use not connected host
		if host.IsConnected {

			// sum commands
			counter += len(host.Commands)
		}
	}

	// return the counter
	return counter
}

// This function converts the dictionary of hosts to the structure needed
// by the dispatcher.
func HostsToDispatcher(hosts configuration.Hosts) []*Host {

	// init the hosts for output
	myhosts := make([]*Host, 0)

	// loop over elements in the map
	for hostname, fields := range hosts {

		// new channel
		channel := make(chan int)

		// create a host object
		host := &Host{
			Hostname:    hostname,
			Port:        fields.Port,
			Protocol:    fields.Protocol,
			IsConnected: true,
			Waiter:      &channel,
		}

		// append to hosts
		myhosts = append(myhosts, host)

	}

	return myhosts

}

// This function converts users structure from the configuration file into
// the structure used by the dispatcher.
func UsersToDispatcher(users configuration.Users) []*Command {

	// init the command slice
	commands := make([]*Command, 0)

	// loop over users
	for username, fields := range users {

		// create an user structure
		user := &User{
			Name:     username,
			Password: fields.Password,
			Identity: 0,
		}

		// loop over commands
		for _, command := range fields.Commands {

			commands = append(
				commands,
				&Command{
					Command: command,
					User:    user,
				},
			)

		}
	}

	// returns the list of commands
	return commands

}
