package dispatcher

import (
	myusr "os/user"
	"testing"

	"github.com/ElricleNecro/TOD/commands"
	"github.com/ElricleNecro/TOD/configuration"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/host"
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

	// Read the user structure from the test file
	usr, _ := myusr.Current()
	users := configuration.ReadUsersYAML(usr.HomeDir + "/CONFIG/TOD/users/users.yaml")

	// configuration
	conf := &configuration.Config{}
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
	conf.Stdin = true

	// read command from the example configuration
	cmds := configuration.UsersToDispatcher(*users)

	// replicate the command in the example
	commands := make([]commands.Command, 121)
	for i := range commands {
		commands[i] = cmds[0]
	}

	// Create the list of commands and hosts
	hsts := new(host.Hosts)
	hosts := make([]*host.Host, len(hostnames))
	for i, hst := range hostnames {
		// Create the host object in the slice
		hosts[i] = &host.Host{
			Hostname: hst,
		}
	}
	hsts.Hosts = hosts

	// display
	formatter.ColoredPrintln(
		formatter.Blue,
		false,
		"All data initialized !",
	)

	// Create dispatcher
	dispatcher := New(conf, hsts)

	// Dispatch commands on hosts
	hsts.Dispatcher(
		commands,
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
	dispatcher.RunCommands(len(commands))

	// display
	formatter.ColoredPrintln(
		formatter.Blue,
		false,
		"Commands done !",
	)
}
