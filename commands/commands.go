package commands

import (
	"github.com/ElricleNecro/TOD/checker"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/ssh"
)

// This function runs a single command on the host on argument.
func OneCommand(
	host *formatter.Host,
	user *formatter.User,
	command string,
	timeout int,
	disconnected chan<- *formatter.Host,
) (string, error) {

	// create a session object
	session := ssh.New(
		user,
		host,
	)

	// check the connection to the host
	if is, err := checker.IsConnected(host, timeout); !is || (err != nil) {

		// disconnect
		Disconnecter(
			"Can't connect to host "+host.Hostname+"\n"+
				"Reason is: "+err.Error(),
			host,
			disconnected,
		)

		// exit the loop
		return "", err
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
			"Can't create connection to host "+host.Hostname+"\n"+
				"Reason is: "+err.Error(),
			host,
			disconnected,
		)

		// exit the loop
		return "", err
	}

	// add a session to connect to host
	_, err = session.AddSession()
	if err != nil {

		// disconnect
		Disconnecter(
			"Problem when adding a session to the host!\n"+
				"Reason is: "+err.Error(),
			host,
			disconnected,
		)

		// Close the session
		session.Close()

		// exit the loop
		return "", err
	}

	// execute the command on the host
	formatter.ColoredPrintln(
		formatter.Green,
		false,
		"Execute command on", host.Hostname,
	)
	output, err2 := session.Run(command)
	if err2 != nil {

		// disconnect
		Disconnecter(
			"An error occurred during the execution of the command !\n"+
				"The command was: "+command+
				"\nand the host is: "+host.Hostname+
				"\nError information: "+err2.Error(),
			host,
			disconnected,
		)

		// Close the session
		session.Close()

		// exit the loop
		return "", err2
	}

	// Close the session
	session.Close()

	// return nil if good
	return output, nil
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
