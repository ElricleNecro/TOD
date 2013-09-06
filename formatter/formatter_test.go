package formatter

import (
	"fmt"
	"testing"
)

func TestDispatcher(t *testing.T) {

	// A list of host
	hostnames := [5]string{
		"localhost",
		"127.0.0.1",
		"192.168.1.1",
		"carmenere",
		"chasselas.iap.fr",
	}

	// A list of fake commands
	mycommands := [11]string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
	}

	// Create the list of commands and hosts
	hosts := make([]*Host, len(hostnames))
	for i, host := range hostnames {

		// Create the host object in the slice
		hosts[i] = &Host{
			Hostname: host,
			Port:     22,
			Protocol: "tcp",
		}
	}
	commands := make([]*Command, len(mycommands))
	for i, command := range mycommands {

		// Create a part of the Command object
		commands[i] = &Command{
			Command: command,
		}
	}

	// display
	t.Log("Setting data for test done!")

	// Dispatch commands to hosts
	Dispatcher(
		commands,
		hosts,
	)

	// display
	t.Log("Dispatching commands on hosts done!")

	// Print the host and the associated commands
	for _, host := range hosts {

		// loop over commands
		for _, command := range host.Commands {

			// display the host and the associated command
			fmt.Println(host.Hostname + "\t" + command.Command)
		}
	}
}
