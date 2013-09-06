// The package containing simplified interface for the SSH
// library.
package ssh

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	//"fmt"
	"github.com/ElricleNecro/TOD/formatter"
	"strconv"
)

// A structure containing all information necessary to an SSH connection
type Session struct {

	// The configuration for the session
	Config *ssh.ClientConfig

	// The structure of the client
	Client *ssh.ClientConn

	// The structure of the session
	Session []*ssh.Session

	// The user associated
	User *formatter.User

	// The host associated
	Host *formatter.Host
}

type clientPassword string

func (p clientPassword) Password(user string) (string, error) {
	return string(p), nil
}

// A method for the construction of the configuration
// object necessary for the connection to the host.
func (s *Session) NewConfig(
	user *formatter.User,
) error {

	// Construct the configuration with password authentication
	s.Config = &ssh.ClientConfig{
		User: user.Name,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthPassword(clientPassword(user.Password)),
		},
	}

	return nil
}

// Construction of a new session between a user and a given host whose
// properties are defined in the associated object.
func New(
	user *formatter.User,
	host *formatter.Host,
) *Session {

	// create the session
	var session *Session = new(Session)

	// Set the user and host for this session
	session.User = user
	session.Host = host

	// create a new configuration
	session.NewConfig(
		user,
	)

	// return the session
	return session

}

// This function allows to connect to the host to create sessions on it
// after it.
func (s *Session) Connect() error {

	var err error

	// create a new client for dialing with the host
	s.Client, err = ssh.Dial(
		s.Host.Protocol,
		s.Host.Hostname+":"+strconv.Itoa(s.Host.Port),
		s.Config,
	)

	return err

}

// Function to add a session to the connection to the host.
// Since multiple sessions can exist for a connection, we allow
// the possibility to append a session into a list of session.
// The function returns too the created session in order to
// have an easy access to the session newly created.
// TODO: Maybe add them into a dictionary in order to allow to
// use a name for retrieving the session as in tmux, etc... just by
// typing a name.
func (s *Session) AddSession() (*ssh.Session, error) {

	// create the session
	session, err := s.Client.NewSession()

	// append the session to the list
	if err != nil {
		panic("Failed to create the session to the host !")
	} else {
		s.Session = append(s.Session, session)
	}

	// return the result
	return session, err

}

// Run the command in argument using by default the last session
// created for this connection
func (s *Session) Run(command string) (string, error) {

	// Create a bytes buffer and affect it as Stdout writer
	// to return the output of the command
	var b bytes.Buffer

	// the number of sessions minus 1 at this time
	nmax := len(s.Session) - 1

	// get the good session
	session := s.Session[nmax]

	// Affect the output to the buffer
	session.Stdout = &b

	// run the command on the host
	err := session.Run(command)

	// return the output and the error if one present
	return b.String(), err

}
