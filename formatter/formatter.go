package formatter

//import (
//"strings"
//)

// The structure containing information for the connection as a given
// user to a host.
type User struct {

	// The user name
	Name string

	// An identity to use as reference for the user if necessary
	Identity int

	// The password to the host if necessary
	Password string
}

// This structure contains the list of commands to execute on the
// corresponding host. So we allow to easily dispatch those commands
// on the host.
type Host struct {

	// the host
	Hostname string

	// The port on which to connect
	Port int

	// The protocol to use in the host for the connection
	Protocol string
}
