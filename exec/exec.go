package exec

import (
	"github.com/ElricleNecro/TOD/checker"
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/ssh"
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

			// display
			formatter.ColoredPrintln(
				formatter.Blue,
				true,
				"Executing ",
				len(host.Commands),
				" commands for ", host.Hostname,
			)

			// loop over commands on this hosts
			for i := host.CommandNumber; i < len(host.Commands); i++ {

				// number of the command
				host.CommandNumber = i

				// create a session object
				session := ssh.New(
					host.Commands[i].User,
					host,
				)

				// check the connection to the host
				if is, err := checker.IsConnected(host, config.Timeout); !is || (err != nil) {

					// display
					formatter.ColoredPrintln(
						formatter.Red,
						false,
						"Can't connect to host ", host.Hostname,
					)

					// dispatch remaining work to other hosts
					select {
					case disconnected <- host:
					default:
					}

					// exit the loop
					break loop
				}

				// display that host is connected
				formatter.ColoredPrintln(
					formatter.Green,
					false,
					"Host ",
					host.Hostname,
					" seems to be online!",
				)

				// Attempt a connection to the host
				err := session.Connect()

				// check the host can be called
				if err != nil {

					// display
					formatter.ColoredPrintln(
						formatter.Red,
						false,
						"Can't connect to host ", host.Hostname,
					)

					// dispatch remaining work to other hosts
					select {
					case disconnected <- host:
					default:
					}

					// exit the loop
					break loop
				}

				// add a session to connect to host
				_, err = session.AddSession()
				if err != nil {

					// display
					formatter.ColoredPrintln(
						formatter.Red,
						false,
						"Problem when adding a session to the host !",
					)

					// dispatch remaining work to other hosts
					select {
					case disconnected <- host:
					default:
					}

					// exit the loop
					break loop
				}

				// execute the command on the host
				formatter.ColoredPrintln(
					formatter.Green,
					false,
					"Execute command on ", host.Hostname,
				)
				output, err2 := session.Run(host.Commands[i].Command)
				if err2 != nil {

					// display
					formatter.ColoredPrintln(
						formatter.Red,
						false,
						"An error occurred during the execution ",
						"of the command !",
					)
					formatter.ColoredPrintln(
						formatter.Red,
						false,
						"The command was: ", host.Commands[i].Command,
					)
					formatter.ColoredPrintln(
						formatter.Red,
						false,
						"and the host is: ", host.Hostname,
					)
					formatter.ColoredPrintln(
						formatter.Red,
						false,
						"Error information: ", err2.Error(),
					)

					// dispatch remaining work to other hosts
					select {
					case disconnected <- host:
					default:
					}

					// exit the loop
					break loop
				}

				// The command has been executed correctly, say it to other
				ender <- true

				// Close the session
				session.Close()

				// for now print the result of the command
				formatter.ColoredPrintln(
					formatter.Magenta,
					false,
					output,
				)

				// wait here for new jobs
				if i == len(host.Commands)-1 {

					// display
					formatter.ColoredPrintln(
						formatter.Magenta,
						true,
						"Waiting more jobs for ", host.Hostname,
					)

					// Now wait for new job
					<-*(host.Waiter)

					// display
					formatter.ColoredPrintln(
						formatter.Magenta,
						true,
						host.Hostname+" has more jobs !",
					)
					formatter.ColoredPrintln(
						formatter.Green,
						true,
						"Number of commands for ", host.Hostname, " :",
						len(host.Commands),
					)

				}

			}

		} else {

			// Now wait for new job
			<-*(host.Waiter)

		}

	}

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
			"Dispatch the jobs of ", host.Hostname,
			" to other connected hosts !",
		)

		// mark the host as not connected
		host.IsConnected = false

		// dispatch remaining work to other hosts
		formatter.Dispatcher(
			host.Commands[host.CommandNumber:],
			hosts,
			false,
		)

		// Set the commands to nothing for the sorter of the dispatcher
		host.Commands = make([]*formatter.Command, 0)

		// display
		formatter.ColoredPrintln(
			formatter.Green,
			false,
			"Dispatching done for ", host.Hostname, " !",
		)
	}

}

// This function sets a timer and check the remaining number of commands
// to execute.
func RemainingCommands(
	hosts []*formatter.Host,
	ncommands int,
) {

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
				counter += len(host.Commands)
			}
		}

		// display the result
		formatter.ColoredPrintln(
			formatter.Magenta,
			false,
			"Number of commands remaining :",
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
			hour, ":",
			minute, ":",
			second,
		)
	}
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
		RemainingCommands(
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
	}

}

//vim: spelllang=en
