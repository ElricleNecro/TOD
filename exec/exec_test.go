package exec

import (
	"bytes"
	"fmt"
	"github.com/ElricleNecro/TOD/formatter"
	"os/exec"
	"strings"
	"testing"
)

func TestRunCommands(t *testing.T) {

	// create an user
	user := &formatter.User{
		Name:     "duarte",
		Identity: 1,
		Password: "SYmadu10;",
	}

	// A list of host
	hostnames := []string{
		"chasselas.iap.fr",
		//"vaccarese",
		//"carmenere",
		//"babel",
	}

	// Create a command which will be duplicated
	command := &formatter.Command{
		Command: "/bin/hostname",
		User:    user,
	}
	commands := make([]*formatter.Command, 5)
	for i, _ := range commands {
		commands[i] = command
	}

	// Create the list of commands and hosts
	hosts := make([]*formatter.Host, len(hostnames))
	for i, host := range hostnames {

		// Create the host object in the slice
		hosts[i] = &formatter.Host{
			Hostname:    host,
			Port:        22,
			Protocol:    "tcp",
			IsConnected: true,
		}
	}

	// display
	fmt.Println("All data initialized !")

	// Dispatch commands on hosts
	formatter.Dispatcher(
		commands,
		hosts,
	)

	// display
	fmt.Println("Dispatching the commands on hosts done !")

	// Run commands in concurrent
	RunCommands(hosts)

	// display
	fmt.Println("Commands done !")

}

func TestLaunchTask(t *testing.T) {
	command := exec.Command("tr", "a-z", "A-Z")
	command.Stdin = strings.NewReader("Un petit test")
	var out bytes.Buffer
	command.Stdout = &out

	res := LaunchTask(command, nil)
	if <-res {
		t.Log(out.String())
	}

	if out.String() != "UN PETIT TEST" {
		t.Error("Loupé, la commande a donnée : ", out.String(), " au lieu de ", "UN PETIT TEST")
	}
}

// A simple test function for the range taking into account the new length
// of a slice or not.
func TestSlices(t *testing.T) {

	// A slice for test
	test := []int{1, 2, 3, 4, 5}

	// loop over the slice
	for i, bidule := range test {

		if i == 2 {
			test = append(test, 6, 7, 8, 9)
		}

		fmt.Println(bidule)
	}

	// with a simple loop
	for i := 0; i < len(test); i++ {

		if i == 2 {
			test = append(test, 6, 7, 8, 9)
		}

		fmt.Println(test[i])
	}

}
