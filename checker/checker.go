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
	conn := make(chan bool, len(hosts))
	for i, _ := range hosts {

		go func(myhost *formatter.Host) {
			res, err := IsConnected(myhost, 10)
			// check is connected
			if err == nil && res {

				// add the host to the list
				connected[nconn] = myhost

				// increment counter
				nconn++
			}

			conn <- true
		}(hosts[i])

	}

	for i := 0; i < len(hosts); i++ {
		<-conn
	}

	// return the list
	return connected[:nconn]

}

// Function which check that an host is connected or not by checking errors
// when attempting to connect to it and by setting a timer for the connection
// timeout if nothing is responding.
func IsConnected(
	host *formatter.Host,
	timeout int,
) (bool, error) {

	// create a dialer
	dial := net.Dialer{
		Deadline:  time.Now().Add(time.Duration(timeout) * time.Second),
		Timeout:   time.Duration(timeout) * time.Second,
		LocalAddr: nil,
	}

	// Contact the host and if no error, it is connected
	_, err := dial.Dial(
		host.Protocol,
		net.JoinHostPort(host.Hostname, strconv.Itoa(host.Port)),
	)

	return err == nil, err
}
