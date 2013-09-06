// Package with function to check if an host is accessible, or connected.
//
// A set of functions to check and return a list of host with their
// information.
package checker

import (
	"github.com/ElricleNecro/TOD/formatter"
	"net"
	"strconv"
	"time"
)

// The function delete host from if they are not accessible
// or not connected with the provided information.
func Checker(
	hosts []*formatter.Host,
) []*formatter.Host {

	// empty slice of connected hosts
	connected := make([]*formatter.Host, len(hosts))

	// Counter for the number of host connected
	nconn := 0

	// loop over the hosts in argument and check they are connected
	for _, host := range hosts {

		// check is connected
		if res, err := IsConnected(host); err == nil && res {

			// add the host to the list
			connected[nconn] = host

			// increment counter
			nconn++
		}

	}

	// return the list
	return connected[:nconn]

}

func IsConnected(host *formatter.Host) (bool, error) {

	// create a dialer to contact the host
	dial := net.Dialer{
		Deadline: time.Now().Add(time.Duration(30) * time.Second),
	}

	// Contact the host and if no error, it is connected
	_, err := dial.Dial(
		host.Protocol,
		host.Hostname+":"+strconv.Itoa(host.Port),
	)

	// return the result
	return err == nil, err
}
