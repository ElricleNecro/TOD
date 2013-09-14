package exec

import (
	"github.com/ElricleNecro/TOD/checker"
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/log_command"
	"github.com/ElricleNecro/TOD/ssh"
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

				// create a session object
				session := ssh.New(
					host.Commands[i].User,
					host,
				)

				// check the connection to the host
				if is, err := checker.IsConnected(host, config.Timeout); !is || (err != nil) {

					// disconnect
					Disconnecter(
						"Can't connect to host "+host.Hostname,
						host,
						disconnected,
					)

					// exit the loop
					break loop
				}

				// display that host is connected
				formatter.ColoredPrintln(
					formatter.Green,
					false,
					"Host",
					host.Hostname,
					"seems to be online!",
				)

				// Attempt a connection to the host
				err := session.Connect()

				// check the host can be called
				if err != nil {

					// disconnect
					Disconnecter(
						"Can't connect to host "+host.Hostname,
						host,
						disconnected,
					)

					// exit the loop
					break loop
				}

				// add a session to connect to host
				_, err = session.AddSession()
				if err != nil {

					// disconnect
					Disconnecter(
						"Problem when adding a session to the host !",
						host,
						disconnected,
					)

					// exit the loop
					break loop
				}

				// execute the command on the host
				formatter.ColoredPrintln(
					formatter.Green,
					false,
					"Execute command on", host.Hostname,
				)
				output, err2 := session.Run(host.Commands[i].Command)
				if err2 != nil {

					// disconnect
					Disconnecter(
						"An error occurred during the execution of the command !\n"+
							"The command was: "+host.Commands[i].Command+
							"and the host is: "+host.Hostname+
							"\nError information: "+err2.Error(),
						host,
						disconnected,
					)

					// exit the loop
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

				// Close the session
				session.Close()

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

// Function to execute a disconnection of host with a command.
func Disconnecter(
	message string,
	host *formatter.Host,
	disconnected chan<- *formatter.Host,
) {

	// display
	formatter.ColoredPrintln(
		formatter.Red,
		false,
		message,
	)

	// dispatch remaining work to other hosts
	disconnected <- host

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
