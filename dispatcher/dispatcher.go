package dispatcher

import (
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/host"
	"github.com/ElricleNecro/TOD/tools"
)

type Dispatcher struct {

	// the list of hosts
	Hosts *host.Hosts

	// the configuration
	Config *configuration.Config

	// for the connection
	ender chan bool

	// for the waiting
	wait chan int

	// chan for the disconnection
	disconnected chan *host.Host
}

// to make an instance of dispatcher
func New(
	config *configuration.Config,
	hosts *host.Hosts,
) *Dispatcher {
	// create a dispatcher
	dispatcher := new(Dispatcher)
	dispatcher.Config = config
	dispatcher.Hosts = hosts

	// A channel to wait for dispatching
	dispatcher.disconnected = make(chan *host.Host)

	// a channel to wait for hosts
	dispatcher.wait = make(chan int)

	// set hosts from the CLI
	dispatcher.HostsFromCLI()

	return dispatcher
}

// This function is used to run a command on a host
// with supplied informations.
func (d *Dispatcher) RunOnHost(
	host *host.Host,
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
				if d.Config.ExcludeLoaded {
					if is, err := host.IsTooLoaded(
						host.Commands[i].User,
						d.Config,
					); is && err == nil {
						d.Disconnect(
							host,
							"The host "+host.Hostname+" is too loaded!",
						)
						break loop
					} else if err != nil {
						d.Disconnect(
							host,
							"Problem occurred when checking loading for "+
								host.Hostname+"!\n"+
								"Reason is: "+err.Error(),
						)
						break loop
					}
				}

				// say the host is working
				host.IsWorking = true

				// Execute the command on the specified host
				output, err := host.OneCommand(&host.Commands[i])
				if err != nil {
					d.Disconnect(host, output)
					break loop
				}

				// The command has been executed correctly, say it to other
				d.ender <- true

				// Write the log of the command
				tools.WriteLogCommand(
					output,
					d.Config,
					host.Hostname,
					host.Commands[i].Command,
					i,
				)

				// for now print the result of the command
				if !d.Config.NoResults {
					formatter.ColoredPrintln(
						formatter.Magenta,
						false,
						output,
					)
				}

				// wait here for new jobs
				if i == len(host.Commands)-1 {
					//Wait for other hosts
					host.Waiter()
				}
			}
		} else {
			// Now wait for new job
			host.Waiter()
		}
	}
}

// Function to dispatch an host on other. Set variables to allow a good synchronisation
// between go routines.
func (d *Dispatcher) Disconnection() {
	for {
		// display
		formatter.ColoredPrintln(
			formatter.Green,
			false,
			"Waiting for a disconnected host !",
		)

		// wait for a signal from a disconnected host
		host := <-d.disconnected

		// display
		formatter.ColoredPrintln(
			formatter.Green,
			false,
			"Dispatch the jobs of", host.Hostname,
			"to other connected hosts !",
		)

		// mark the host as not connected
		host.Connected = false

		// and not working
		host.IsWorking = false

		// dispatch remaining work to other hosts
		d.Hosts.Dispatcher(
			host.Commands[host.CommandNumber-1:],
			d.Config.HostsMax,
			false,
		)

		// display
		formatter.ColoredPrintln(
			formatter.Green,
			false,
			"Dispatching done for", host.Hostname, "!",
		)
	}
}

// This function loop over hosts to launch commands in concurrent mode.
func (d *Dispatcher) RunCommands(
	ncommands int,
) {

	// number of hosts
	nhosts := len(d.Hosts.Hosts)

	// check that there is some hosts
	if nhosts == 0 {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"There is no hosts given to run commands !",
		)
	}

	// A channel to wait for end of program
	d.ender = make(chan bool, ncommands)

	// run the routine to manage disconnections
	go d.Disconnection()

	// if we use the timer, run the go routine
	if d.Config.Timer {
		go d.Hosts.RemainingCommands(
			d.Config,
			ncommands,
		)
	}

	// if we use the timer of working hosts, run the go routine
	if d.Config.WorkTimer {
		go d.Hosts.WorkingTimer(
			d.Config.WorkTime,
		)
	}

	// loop over hosts and run the command
	for _, host := range d.Hosts.Hosts {
		// in several goroutine
		go d.RunOnHost(host)
	}

	// Wait for the end of goroutines
	for i := 0; i < ncommands; i++ {
		<-d.ender
		formatter.ColoredPrintln(
			formatter.White,
			true,
			"Number of remaining commands:", ncommands-1-i,
		)
	}

	// display the summary of commands on hosts
	d.Hosts.DisplayHostsCommands()

}

// Function to execute a disconnection of host with a command.
func (d *Dispatcher) Disconnect(
	host *host.Host,
	message string,
) {

	// display
	formatter.ColoredPrintln(
		formatter.Red,
		false,
		message,
	)

	// dispatch remaining work to other hosts
	d.disconnected <- host
}

// set some properties of hosts with the configuration from command line
func (d *Dispatcher) HostsFromCLI() {

	// loop over hosts
	for _, host := range d.Hosts.Hosts {
		// make a channel for waiting new jobs
		channel := make(chan int)
		// do it only if we use the stdin to set hosts
		if d.Config.Stdin {
			host.Timeout = d.Config.Timeout
			host.Port = d.Config.Port
			host.Protocol = d.Config.Protocol
		}
		host.Connected = true
		host.IsWorking = false
		host.Wait = &channel
	}
}
