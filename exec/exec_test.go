package exec

import (
	"fmt"
	"github.com/ElricleNecro/TOD/formatter"
	color "github.com/daviddengcn/go-colortext"
	"testing"
)

func TestRunCommands(t *testing.T) {

	// create an user
	user := &formatter.User{
		Name:     "",
		Identity: 1,
		Password: "",
	}

	// A list of host
	hostnames := []string{
		"chasselas.iap.fr",
		"vaccarese",
		"carmenere",
		"tokay",
		"null",
		//"carmenere",
		//"babel",
	}

	// Create a command which will be duplicated
	command := &formatter.Command{
		Command: "/bin/hostname",
		User:    user,
	}
	commands := make([]*formatter.Command, 20)
	for i, _ := range commands {
		commands[i] = command
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

		fmt.Println(host)
		fmt.Println(hosts[i].Waiter)
	}

	// display
	color.ChangeColor(color.Blue, false, color.None, false)
	fmt.Println("All data initialized !")

	// Dispatch commands on hosts
	formatter.Dispatcher(
		commands,
		hosts,
		true,
	)

	// display
	color.ChangeColor(color.Blue, false, color.None, false)
	fmt.Println("Dispatching the commands on hosts done !")

	// Run commands in concurrent
	RunCommands(hosts, len(commands))

	// display
	color.ChangeColor(color.Blue, false, color.None, false)
	fmt.Println("Commands done !")

}
