package exec

import (
	"github.com/ElricleNecro/TOD/formatter"
	"testing"
)

var (

	// User
	user = &formatter.User{
		Name:     "perceval",
		Identity: 1,
		Password: "sloubi",
	}

	// A list of host
	hostnames = []string{
		"chasselas.iap.fr",
		"vaccarese",
		"carmenere",
		"tressalier",
		"null",
		"tockay",
	}
)

// To test the run of commands.
func TestRunCommands(t *testing.T) {

	// Create a command which will be duplicated
	command := &formatter.Command{
		Command: "/bin/hostname",
		User:    user,
	}
	command2 := &formatter.Command{
		Command: "whoami",
		User:    user,
	}
	commands := make([]*formatter.Command, 25)
	for i, _ := range commands {
		commands[i] = command
	}
	for i := 12; i < 16; i++ {
		commands[i] = command2
	}

	// Create the list of commands and hosts
	hosts := make([]*formatter.Host, len(hostnames))
	for i, host := range hostnames {

		// new channel
		channel := make(chan int)

		// Create the host object in the slice
		hosts[i] = &formatter.Host{
			Hostname:    host,
			Port:        22,
			Protocol:    "tcp",
			IsConnected: true,
			Waiter:      &channel,
		}

	}

	// display
	formatter.ColoredPrintln(formatter.Blue, false, "All data initialized !")

	// Dispatch commands on hosts
	formatter.Dispatcher(
		commands,
		hosts,
		true,
	)

	// display
	formatter.ColoredPrintln(
		formatter.Blue,
		false,
		"Dispatching the commands on hosts done !",
	)

	// Run commands in concurrent
	RunCommands(hosts, len(commands))

	// display
	formatter.ColoredPrintln(formatter.Blue, false, "Commands done !")

}
