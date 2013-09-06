package checker

import (
	"fmt"
	"github.com/ElricleNecro/TOD/formatter"
	"testing"
)

func TestChecker(t *testing.T) {

	// A list of host
	hostnames := [5]string{
		"localhost",
		"127.0.0.1",
		"192.168.1.1",
		"carmenere",
		"chasselas.iap.fr",
	}

	// Create the list of host object
	hosts := make([]*formatter.Host, len(hostnames))
	for i, host := range hostnames {

		// Create the host object in the slice
		hosts[i] = &formatter.Host{
			Hostname: host,
			Port:     22,
			Protocol: "tcp",
		}
	}

	// Now get the really connected hosts
	hosts = Checker(hosts)

	// Print their names
	for _, host := range hosts {
		fmt.Println(host.Hostname)
	}

}
