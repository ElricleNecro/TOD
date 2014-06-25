package exec

import (
	myusr "os/user"
	"testing"

	config "github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/formatter"
)

var (

	// A list of host
	hostnames = []string{
		"arbois",
		"verdejo",
		"courbu",
		"aspiran",
		"tressalier",
		"picardan",
		"roussanne",
		"doucillon",
		"molette",
		"mauzac",
		"mancin",
		"vaccarese",
		"carmenere",
		"null",
		"tockay",
		"bidule",
		"toké",
		"poulsard",
		"machin",
		"ugni",
	}
)

// To test the run of commands.
func TestRunCommands(t *testing.T) {

	var user *formatter.User

	// Read the user structure from the test file
	usr, _ := myusr.Current()
	users := config.ReadUsersYAML(usr.HomeDir + "/CONFIG/TOD/users/users.yaml")
	for myuser, fields := range *users {

		user = &formatter.User{
			Name:     myuser,
			Identity: 1,
			Key:      fields.Key,
		}
	}

	// configuration
	conf := &config.Config{}
	conf.Port = 22
	conf.Protocol = "tcp"
	conf.Timeout = 10
	conf.LogCommand = "/tmp"
	conf.CPUMax = 25.0
	conf.MemoryMax = 30.0
	conf.ExcludeLoaded = true
	conf.WorkTimer = true
	conf.WorkTime = 120
	conf.HostsMax = 5

	// Create a command which will be duplicated
	command := &formatter.Command{
		Command: "sleep $(( RANDOM % 10 )) && /bin/hostname",
		User:    user,
	}
	commands := make([]*formatter.Command, 121)
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

	}

	// display
	formatter.ColoredPrintln(formatter.Blue, false, "All data initialized !")

	// Dispatch commands on hosts
	formatter.Dispatcher(
		commands,
		hosts,
		conf.HostsMax,
		true,
	)

	// display
	formatter.ColoredPrintln(
		formatter.Blue,
		false,
		"Dispatching the commands on hosts done !",
	)

	// Run commands in concurrent
	RunCommands(hosts, len(commands), conf)

	// display
	formatter.ColoredPrintln(formatter.Blue, false, "Commands done !")

}
