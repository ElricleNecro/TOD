package exec

import (
	"fmt"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/ssh"
	"log"
	"os/exec"
)

// This function loop over hosts to launch commands in concurrent mode.
func RunCommands(hosts []*formatter.Host) {

	// functions definition
	var Disconnection func(
		hosts []*formatter.Host,
		host *formatter.Host,
		commands []*formatter.Command,
	)
	var RunOnHost func(
		hosts []*formatter.Host,
		host *formatter.Host,
	)

	// number of hosts
	nhosts := len(hosts)

	// check that there is some hosts
	if nhosts == 0 {
		fmt.Println("There is no hosts given to run commands !")
	}

	// create a buffer channel to wait for end of each goroutine
	buffer := make(chan bool, nhosts)

	// variable to say if one host is being dispatched
	InDisconnection := false

	// A channel to wait for dispatching
	dispatch := make(chan bool)

	// Function to dispatch an host on other. Set variables to allow a good synchronisation
	// between go routines.
	Disconnection = func(
		hosts []*formatter.Host,
		host *formatter.Host,
		commands []*formatter.Command,
	) {

		// if one host is being disconnected, wait for the end
		if InDisconnection {
			<-dispatch
		}

		// Say that one is in disconnection
		InDisconnection = true

		// mark the host as not connected
		host.IsConnected = false

		// dispatch remaining work to other hosts
		formatter.Dispatcher(
			commands,
			hosts,
		)

		// say that no one is in disconnection
		InDisconnection = false

		// Send a signal to say that we can continue
		dispatch <- true

	}

	// This function is used to run a command on a host
	// with supplied informations.
	RunOnHost = func(
		hosts []*formatter.Host,
		host *formatter.Host,
	) {

	loop:
		// loop over commands on this hosts
		for i := 0; i < len(host.Commands); i++ {

			// create a session object
			session := ssh.New(
				host.Commands[i].User,
				host,
			)

			// Attempt a connection to the host
			err := session.Connect()

			// check the host can be called
			if err != nil {

				// display
				fmt.Println("Can't connect to host " + host.Hostname)

				// dispatch remaining work to other hosts
				Disconnection(
					hosts,
					host,
					host.Commands[i:],
				)

				// exit the loop
				break loop
			}

			// add a session to connect to host
			_, err3 := session.AddSession()
			if err3 != nil {

				// display
				fmt.Println(
					"Problem when adding a session to the host !",
				)

				// dispatch remaining work to other hosts
				Disconnection(
					hosts,
					host,
					host.Commands[i:],
				)

				// exit the loop
				break loop
			}

			// execute the command on the host
			output, err2 := session.Run(host.Commands[i].Command)
			if err2 != nil {
				fmt.Println(
					"An error occurred during the execution of the command !",
				)
				fmt.Println(
					"The command was: " + host.Commands[i].Command,
				)
				fmt.Println("Error information: " + err.Error())
			}

			// Close the session
			session.Close()

			// for now print the result of the command
			fmt.Println(output)

		}

		// say that the go routine has ended
		buffer <- true

	}

	// loop over hosts and run the command
	for _, host := range hosts {

		// in several goroutine
		go RunOnHost(hosts, host)

	}

	// Wait for the end of goroutines
	for i := 0; i < nhosts; i++ {
		<-buffer
	}

}

//Launch the task in a thread and wait for the command to terminate. Return a channel for getting the result state.
// False is send to the channel and a message is write in the Logger if their was was a problem, true elsewhere.
func LaunchTask(task *exec.Cmd, tlog *log.Logger) <-chan bool {
	c := make(chan bool)
	go func() {
		if err := task.Run(); err != nil {
			if tlog != nil {
				tlog.Println(err)
			}
			c <- false
		} else {
			c <- true
		}
	}()

	return c
}

//vim: spelllang=en
