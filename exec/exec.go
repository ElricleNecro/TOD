package exec

import (
	"github.com/ElricleNecro/TOD/commands"
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/load_checker"
	"github.com/ElricleNecro/TOD/log_command"
	"strconv"
	"time"
)

// This function is used to run a command on a host
// with supplied informations.
func RunOnHost(
	hosts []*formatter.Host,
	host *formatter.Host,
	config *configuration.Config,
	disconnected chan<- *formatter.Host,
	ender chan<- bool,
) {

loop:
	// Do an infinite loop for waiting when ended
	for {

		// check the size of commands to execute before
		if len(host.Commands) != 0 {

			// loop over commands on this hosts
			for i := host.CommandNumber; i < len(host.Commands); i++ {

				// display
				formatter.ColoredPrintln(
					formatter.Blue,
					true,
					"Executing command", i, "for", host.Hostname,
				)

				// number of the command
				host.CommandNumber = i + 1

				// check if we want to exclude too loaded hosts
				if config.ExcludeLoaded {
					load_checker.CheckLoaded(
						host,
						host.Commands[i].User,
						config.Timeout,
						config.CPUMax,
						config.MemoryMax,
						disconnected,
					)
				}

				// Execute the command on the specified host
				output, err := commands.OneCommand(
					host,
					host.Commands[i].User,
					host.Commands[i].Command,
					config.Timeout,
					disconnected,
				)

				// check the command as executed correctly, else exit loop
				if err != nil {
					break loop
				}

				// The command has been executed correctly, say it to other
				ender <- true

				// Write the log of the command
				log_command.WriteLogCommand(
					output,
					config,
					host.Hostname,
					host.Commands[i].Command,
					i,
				)

				// for now print the result of the command
				if !config.NoResults {
					formatter.ColoredPrintln(
						formatter.Magenta,
						false,
						output,
					)
				}

				// wait here for new jobs
				if i == len(host.Commands)-1 {

					//Wait for other hosts
					Waiter(host)

				}

			}

		} else {

			// Now wait for new job
			Waiter(host)

		}

	}

}

// Function which executes commands when a host has to wait for other hosts.
func Waiter(host *formatter.Host) {

	// display
	formatter.ColoredPrintln(
		formatter.Magenta,
		true,
		"Waiting more jobs for", host.Hostname,
	)

	// Now wait for new job
	<-*(host.Waiter)

	// display
	formatter.ColoredPrintln(
		formatter.Magenta,
		true,
		host.Hostname, "has more jobs !",
	)
	formatter.ColoredPrintln(
		formatter.Green,
		true,
		"Number of commands for", host.Hostname, ":",
		len(host.Commands),
	)

}

// Function to dispatch an host on other. Set variables to allow a good synchronisation
// between go routines.
func Disconnection(
	hosts []*formatter.Host,
	disconnected <-chan *formatter.Host,
) {

	for {

		// display
		formatter.ColoredPrintln(
			formatter.Green,
			false,
			"Waiting for a disconnected host !",
		)

		// wait for a signal from a disconnected host
		host := <-disconnected

		// display
		formatter.ColoredPrintln(
			formatter.Green,
			false,
			"Dispatch the jobs of", host.Hostname,
			"to other connected hosts !",
		)

		// mark the host as not connected
		host.IsConnected = false

		// dispatch remaining work to other hosts
		formatter.Dispatcher(
			host.Commands[host.CommandNumber-1:],
			hosts,
			false,
		)

		// Set the commands to nothing for the sorter of the dispatcher
		//host.Commands = make([]*formatter.Command, 0)

		// display
		formatter.ColoredPrintln(
			formatter.Green,
			false,
			"Dispatching done for", host.Hostname, "!",
		)
	}

}

// This function sets a timer and check the remaining number of commands
// to execute.
func RemainingCommands(
	hosts []*formatter.Host,
	ncommands int,
) {

	// display
	formatter.ColoredPrintln(
		formatter.White,
		true,
		"Timer for commands is set !",
	)

	// set a timer
	timer := time.NewTimer(time.Duration(180) * time.Second)

	// start an infinite loop
	for {

		// Wait for the timer
		<-timer.C

		// counter
		counter := 0

		// loop over hosts
		for _, host := range hosts {

			// check that it is connected
			if host.IsConnected {

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

// A function which displays hosts and the number of commands they have executed.
func DisplayHostsCommands(hosts []*formatter.Host) {

	// counter
	counter := 0

	// loop over hosts
	for _, host := range hosts {

		// display hostname
		formatter.ColoredPrint(
			formatter.Magenta,
			false,
			host.Hostname, ": ",
		)

		// display the number of command executed with different
		// colors in case of a disconnected host
		if host.IsConnected {
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

// This function loop over hosts to launch commands in concurrent mode.
func RunCommands(
	hosts []*formatter.Host,
	ncommands int,
	config *configuration.Config,
) {

	// number of hosts
	nhosts := len(hosts)

	// check that there is some hosts
	if nhosts == 0 {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"There is no hosts given to run commands !",
		)
	}

	// A channel to wait for dispatching
	disconnected := make(chan *formatter.Host)

	// A channel to wait for end of program
	ender := make(chan bool, ncommands)

	// run the routine to manage disconnections
	go Disconnection(
		hosts,
		disconnected,
	)

	// if we use the timer, run the go routine
	if config.Timer {
		go RemainingCommands(
			hosts,
			ncommands,
		)
	}

	// loop over hosts and run the command
	for i, _ := range hosts {

		// in several goroutine
		go RunOnHost(
			hosts,
			hosts[i],
			config,
			disconnected,
			ender,
		)

	}

	// Wait for the end of goroutines
	for i := 0; i < ncommands; i++ {
		<-ender
		formatter.ColoredPrintln(
			formatter.White,
			true,
			"Number of remaining commands:", ncommands-1-i,
		)
	}

	// display the summary of commands on hosts
	DisplayHostsCommands(hosts)

}

//vim: spelllang=en
