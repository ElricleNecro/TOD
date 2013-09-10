package exec

import (
	"fmt"
	"github.com/ElricleNecro/TOD/checker"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/ssh"
	color "github.com/daviddengcn/go-colortext"
	"strconv"
)

// This function loop over hosts to launch commands in concurrent mode.
func RunCommands(hosts []*formatter.Host, ncommands int) {

	var RunOnHost func(
		hosts []*formatter.Host,
		host *formatter.Host,
	)

	// number of hosts
	nhosts := len(hosts)

	// check that there is some hosts
	if nhosts == 0 {
		color.ChangeColor(color.Red, false, color.None, false)
		fmt.Println("There is no hosts given to run commands !")
	}

	// A channel to wait for dispatching
	disconnected := make(chan *formatter.Host)

	// A channel to wait for end of program
	ender := make(chan bool, ncommands)

	// Function to dispatch an host on other. Set variables to allow a good synchronisation
	// between go routines.
	go func() {

		for {

			// display
			color.ChangeColor(color.Green, false, color.None, false)
			fmt.Println("Waiting for a disconnected host !")

			// wait for a signal from a disconnected host
			host := <-disconnected

			// display
			color.ChangeColor(color.Green, false, color.None, false)
			fmt.Println("Dispatch the jobs of " + host.Hostname +
				" to other connected hosts !")

			// mark the host as not connected
			host.IsConnected = false

			// dispatch remaining work to other hosts
			formatter.Dispatcher(
				host.Commands[host.CommandNumber:],
				hosts,
				false,
			)

			// change informations
			host.CommandNumber = 0
			host.Commands = make([]*formatter.Command, 0)

			// display
			color.ChangeColor(color.Green, false, color.None, false)
			fmt.Println("Dispatching done for " + host.Hostname + " !")
		}

	}()

	// This function is used to run a command on a host
	// with supplied informations.
	RunOnHost = func(
		hosts []*formatter.Host,
		host *formatter.Host,
	) {

	loop:
		// Do an infinite loop for waiting when ended
		for {

			// display
			color.ChangeColor(color.Magenta, true, color.Magenta, true)
			fmt.Println("Executing list of commands for " + host.Hostname)

			// loop over commands on this hosts
			for i := 0; i < len(host.Commands); i++ {

				// number of the command
				host.CommandNumber = i

				// create a session object
				session := ssh.New(
					host.Commands[i].User,
					host,
				)

				// check the host can be called
				if is, _ := checker.IsConnected(host); !is {

					// display
					color.ChangeColor(color.Red, false, color.None, false)
					fmt.Println("Can't connect to host " + host.Hostname)

					// dispatch remaining work to other hosts
					disconnected <- host

					// exit the loop
					break loop
				}

				// Attempt a connection to the host
				err := session.Connect()

				// check the host can be called
				if err != nil {

					// display
					color.ChangeColor(color.Red, false, color.None, false)
					fmt.Println("Can't connect to host " + host.Hostname)

					// dispatch remaining work to other hosts
					disconnected <- host

					// exit the loop
					break loop
				}

				// add a session to connect to host
				_, err = session.AddSession()
				if err != nil {

					// display
					color.ChangeColor(color.Red, false, color.None, false)
					fmt.Println(
						"Problem when adding a session to the host !",
					)

					// dispatch remaining work to other hosts
					disconnected <- host

					// exit the loop
					break loop
				}

				// execute the command on the host
				output, err2 := session.Run(host.Commands[i].Command)
				if err2 != nil {

					// display
					color.ChangeColor(color.Red, false, color.None, false)
					fmt.Println(
						"An error occurred during the execution " +
							"of the command !",
					)
					fmt.Println(
						"The command was: " + host.Commands[i].Command,
					)
					fmt.Println("Error information: " + err.Error())

					// exit the loop
					break loop
				}

				// The command has been executed correctly, say it to other
				ender <- true

				// Close the session
				session.Close()

				// for now print the result of the command
				color.ChangeColor(color.Magenta, false, color.None, false)
				fmt.Println(output)

				// wait here for new jobs
				if i == len(host.Commands)-1 {

					// display
					color.ChangeColor(color.Magenta, true, color.None, false)
					fmt.Println("Waiting more jobs for " + host.Hostname)

					// Now wait for new job
					<-*(host.Waiter)

					// display
					color.ChangeColor(color.Magenta, true, color.None, false)
					fmt.Println(host.Hostname + " has more jobs !")
					color.ChangeColor(color.Green, true, color.None, false)
					fmt.Println("Number of commands for " + host.Hostname + " :" +
						strconv.Itoa(len(host.Commands)))

				}

			}

		}

	}

	// loop over hosts and run the command
	for _, host := range hosts {

		fmt.Println(host.Hostname)
		fmt.Println(&(*host))

		// in several goroutine
		go RunOnHost(hosts, host)

	}

	// Wait for the end of goroutines
	for i := 0; i < ncommands; i++ {
		<-ender
	}

}

//vim: spelllang=en
