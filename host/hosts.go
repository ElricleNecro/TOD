package host

import (
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/ElricleNecro/TOD/commands"
	"github.com/ElricleNecro/TOD/formatter"
)

type Hosts struct {
	Hosts []*Host
}

// This function dispatches the commands on some hosts
func (hosts *Hosts) Dispatcher(
	commands []commands.Command,
	nhosts_max int,
	first bool,
) {

	// the pointer to the host structure
	var host int

	// Determine the number of hosts available in theory
	nhosts := hosts.CountConnectedHosts()

	// check there is at least one host connected
	if nhosts == 0 {
		formatter.ColoredPrintln(
			formatter.Red,
			true,
			"There is no hosts available to do the job !",
		)

		// Display not done commands
		displayRemainingCommands(commands)

		// Exit the program
		os.Exit(1)
	}

	// store here the list of selected hosts
	myhost := make([]int, 0)

	// create the sorter
	sorter := &hostSorter{
		Hosts: hosts,
	}

	// Now sort by charge in commands
	sort.Sort(sorter)

	// loop over commands and affect them to hosts
	host = -1
	for _, command := range commands {

		// pass to another host
		host = (host + 1) % nhosts

		// if the host isn't connected
		for !hosts.Hosts[host].Connected {
			// pass to the next host
			host = (host + 1) % nhosts
		}

		// count the number of workers
		nworkers := hosts.CountWorkers()

		// if the maximal number of hosts is get, affect to only working
		// hosts
		if nworkers >= nhosts_max && nhosts_max > 0 {
			for !hosts.Hosts[host].IsWorker() {
				host = (host + 1) % nhosts
			}
		}

		// add the host to list
		myhost = append(myhost, host)

		// append to the list of commands
		hosts.Hosts[host].Commands = append(
			hosts.Hosts[host].Commands,
			command,
		)
	}

	// message to say that the host has more jobs
	if !first {
		// loop over hosts which have more jobs
		for _, host := range myhost {
			// send a non blocking signal
			formatter.ColoredPrintln(
				formatter.None,
				true,
				"Send more jobs signal to", hosts.Hosts[host].Hostname,
				"\nwith new length of", len(hosts.Hosts[host].Commands),
			)
			select {
			case *(hosts.Hosts[host].Wait) <- 1:
			default:
			}
		}
	}
}

// This function counts the number of connected hosts given in
// a slice of those objects.
func (hosts *Hosts) CountConnectedHosts() int {

	// init the counter
	counter := 0

	// loop over host and increment a counter when connected
	for _, host := range hosts.Hosts {

		// check the field for connected
		if host.Connected {
			counter++
		}
	}

	// return the number of connected machines
	return counter
}

// A function which displays hosts and the number of commands they have executed.
func (hosts *Hosts) DisplayHostsCommands() {

	// counter
	counter := 0

	// loop over hosts
	for _, host := range hosts.Hosts {

		// display hostname
		formatter.ColoredPrint(
			formatter.Magenta,
			false,
			host.Hostname, ": ",
		)

		// display the number of command executed with different
		// colors in case of a disconnected host
		if host.Connected {
			formatter.ColoredPrintln(
				formatter.Green,
				false,
				host.CommandNumber, "/", len(host.Commands),
			)
			counter += host.CommandNumber
		} else {
			formatter.ColoredPrintln(
				formatter.Red,
				false,
				host.CommandNumber-1, "/", len(host.Commands),
			)
			counter += host.CommandNumber - 1
		}
	}

	// display the total number to check coherence
	formatter.ColoredPrintln(
		formatter.Magenta,
		false,
		"Total of commands:", counter,
	)

}

// This function computes the number of hosts which have to execute
// commands. It's the number of hosts which are connected and which
// have a given number of command to execute.
func (hosts *Hosts) CountWorkers() int {

	// init the counter
	counter := 0

	// loop over hosts
	for _, host := range hosts.Hosts {
		// increment if the host is a worker as defined before.
		if host.IsWorker() {
			counter++
		}
	}

	// return the counter
	return counter
}

// This function counts the number of remaining commands
func (hosts *Hosts) CountCommands() int {

	// counter
	counter := 0

	// loop over hosts
	for _, host := range hosts.Hosts {
		// don't use not connected host
		if host.Connected {
			// sum commands
			counter += len(host.Commands)
		}
	}

	// return the counter
	return counter
}

// A function to return the remaining commands if the number of hosts
// available is zero
func displayRemainingCommands(commands []commands.Command) {

	// Display message
	formatter.ColoredPrint(
		formatter.Magenta,
		false,
		"The list of not runned commands is:\n",
	)

	// loop over commands
	for _, command := range commands {

		// Display the command
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			command.Command,
		)
	}
}

// A type used to sort the hosts by their charge in term of number of commands
// to execute. The goal is to allow less loaded hosts to run the dispatched
// commands first, so that a bigger number of commands is executed in a given
// laps of time.
type hostSorter struct {
	Hosts *Hosts
}

// Return the length of the array to sort
func (h *hostSorter) Len() int {
	return len(h.Hosts.Hosts)
}

// Swap two hosts
func (h *hostSorter) Swap(i, j int) {
	h.Hosts.Hosts[i], h.Hosts.Hosts[j] = h.Hosts.Hosts[j], h.Hosts.Hosts[i]
}

// The function to compare two hosts
func (h *hostSorter) Less(i, j int) bool {

	// First host length
	ni := len(h.Hosts.Hosts[i].Commands)
	if !h.Hosts.Hosts[i].Connected {
		ni = MaxInt
	}

	// Same with second host
	nj := len(h.Hosts.Hosts[j].Commands)
	if !h.Hosts.Hosts[j].Connected {
		nj = MaxInt
	}
	return ni < nj
}

// This function sets a timer and check the remaining number of commands
// to execute.
func (hosts *Hosts) RemainingCommands(
	config config,
	ncommands int,
) {

	// display
	formatter.ColoredPrintln(
		formatter.White,
		true,
		"Timer for commands is set !",
	)

	// set a timer
	timer := time.NewTicker(time.Duration(config.GetTimer()) * time.Second)

	// start an infinite loop
	for {

		// Wait for the timer
		<-timer.C

		// counter
		counter := 0

		// loop over hosts
		for _, host := range hosts.Hosts {
			// check that it is connected
			if host.Connected {
				// increment counter
				counter += (host.CommandNumber - 1)
			}
		}

		// display the result
		formatter.ColoredPrintln(
			formatter.Magenta,
			false,
			"Number of commands executed :",
			counter,
			"/",
			ncommands,
		)

		// get the current time
		hour, minute, second := time.Now().Clock()
		formatter.ColoredPrintln(
			formatter.Magenta,
			false,
			"at time ",
			strconv.Itoa(hour)+":"+
				strconv.Itoa(minute)+":"+strconv.Itoa(second),
		)
	}
}
