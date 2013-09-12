package checker

import (
	"fmt"
	"github.com/ElricleNecro/TOD/formatter"
	"testing"
)

var (
	// A list of host
	hostnames = []string{
		"localhost",
		"127.0.0.1",
		"192.168.1.1",
		"carmenere",
		"chasselas.iap.fr",
		"null",
		"tockay",
		"tinto",
		"clairette",
	}
)

// A function to create hosts from a list of hosts.
func createList(hostnames []string) []*formatter.Host {

	// Create the list of host object
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

	return hosts

}

// A test for the IsConnected function in concurrent.
func TestIsConnectedGo(t *testing.T) {

	// Create hosts from list
	hosts := createList(hostnames)

	// loop over them and check that they are connected
	for _, host := range hosts {

		go func(myhost *formatter.Host) {
			if is, err := IsConnected(myhost); is && (err == nil) {

				formatter.ColoredPrintln(
					formatter.Green,
					false,
					"Concurrent: ",
					myhost.Hostname,
				)

			} else {

				formatter.ColoredPrintln(
					formatter.Red,
					false,
					"Concurrent: ",
					myhost.Hostname,
				)

			}
		}(host)

	}
}

// A test for the IsConnected function in sequential way.
func TestIsConnected(t *testing.T) {

	// Create hosts from list
	hosts := createList(hostnames)

	// loop over them and check that they are connected
	for _, host := range hosts {

		if is, err := IsConnected(host); is && (err == nil) {

			formatter.ColoredPrintln(
				formatter.Green,
				false,
				host.Hostname,
			)

		} else {

			formatter.ColoredPrintln(
				formatter.Red,
				false,
				host.Hostname,
			)

		}

	}
}

// A test for the checker function which runs in concurrence.
func TestChecker(t *testing.T) {

	// create the list of hosts
	hosts := createList(hostnames)

	// Now get the really connected hosts
	hosts = Checker(hosts)

	// Print their names
	for _, host := range hosts {
		fmt.Println(host.Hostname)
	}

}

// To test if the connection to an host is working.
