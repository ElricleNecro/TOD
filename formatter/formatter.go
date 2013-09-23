package formatter

import (
	"fmt"
	"github.com/ElricleNecro/TOD/configuration"
	color "github.com/daviddengcn/go-colortext"
	"math"
	"sort"
)

// A wrapping type for colors
type Color color.Color

// The types of colors
const (
	None    = color.None
	Black   = color.Black
	Red     = color.Red
	Green   = color.Green
	Yellow  = color.Yellow
	Blue    = color.Blue
	Magenta = color.Magenta
	Cyan    = color.Cyan
	White   = color.White
)

// the maximal value of the integer 32
const (
	MaxInt = 1<<31 - 1
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

// A type used to sort the hosts by their charge in term of number of commands
// to execute. The goal is to allow less loaded hosts to run the dispatched
// commands first, so that a bigger number of commands is executed in a given
// laps of time.
type hostSorter struct {
	Hosts []*Host
	by    func(h1, h2 *Host) bool
}

// Return the length of the array to sort
func (h *hostSorter) Len() int {
	return len(h.Hosts)
}

// Swap two hosts
func (h *hostSorter) Swap(i, j int) {
	h.Hosts[i], h.Hosts[j] = h.Hosts[j], h.Hosts[i]
}

// The function to compare two hosts
func (h *hostSorter) Less(i, j int) bool {

	// First host length
	ni := len(h.Hosts[i].Commands)
	if !h.Hosts[i].IsConnected {
		ni = MaxInt
	}

	// Same with second host
	nj := len(h.Hosts[j].Commands)
	if !h.Hosts[j].IsConnected {
		nj = MaxInt
	}
	return ni < nj
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
		ColoredPrintln(
			Red,
			true,
			"There is no hosts available to do the job !",
		)
	}

	// The same for the number of commands to execute
	ncomm := len(commands)

	// Compute the number of commands per hosts to execute
	NCH := math.Ceil(float64(ncomm) / float64(nhosts))

	// init by pointing the current host to the first one
	host = -1

	// store here the list of selected hosts
	myhost := make([]int, 0)

	// create the sorter
	sorter := &hostSorter{
		Hosts: hosts,
	}

	// Now sort by charge in commands
	sort.Sort(sorter)

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

		// loop over hosts which have more jobs
		for _, host := range myhost {

			// send a non blocking signal
			ColoredPrintln(
				None,
				true,
				"Send more jobs signal to", hosts[host].Hostname,
				"\nwith new length of", len(hosts[host].Commands),
			)
			select {
			case *(hosts[host].Waiter) <- 1:
			default:
			}
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

// A simple function to print message with colors with line return.
func ColoredPrintln(thecolor color.Color, bold bool, values ...interface{}) {

	// change the color of the terminal
	color.ChangeColor(
		thecolor,
		bold,
		None,
		false,
	)

	// print
	fmt.Println(values...)

	// reset the color
	color.ResetColor()
}

// A simple function to print message with colors.
func ColoredPrint(thecolor color.Color, bold bool, values ...interface{}) {

	// change the color of the terminal
	color.ChangeColor(
		thecolor,
		bold,
		None,
		false,
	)

	// print
	fmt.Print(values...)

	// reset the color
	color.ResetColor()
}
