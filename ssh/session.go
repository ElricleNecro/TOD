package ssh

import (
	"bytes"
	"io"
	"net"
	"os"
	"strconv"

	"code.google.com/p/go.crypto/ssh"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/tools"
)

// A structure containing all information necessary to an SSH connection
type Session struct {

	// The configuration for the session
	Config *ssh.ClientConfig

	// The structure of the client
	Client *ssh.Client

	// The structure of the session and their number
	Session *ssh.Session
}

type user interface {
	GetPrivateKey() string
	GetUsername() string
}

type host interface {
	GetProtocol() string
	GetPort() int
	GetHostname() string
}

// Construction of a new session between a user and a given host whose
// properties are defined in the associated object.
func New(user user) (*Session, error) {

	// create the session
	var session *Session = new(Session)

	// create a new configuration
	err := session.NewConfig(user)
	if err != nil {
		return nil, err
	}

	// return the session
	return session, nil
}

// load a private key
func loadPEM(file string) ([]byte, error) {

	// open the file
	f, err := os.Open(file)
	defer func() {
		err := f.Close()
		if err != nil {
			formatter.ColoredPrintln(
				formatter.Red,
				false,
				"The file can't be closed for the private key!\n",
				"Reason is: ", err.Error(),
			)
		}
	}()

	// check errors when opening
	if err != nil {
		return nil, err
	}

	// read data
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	// parse private keys
	return buf.Bytes(), nil
}

// A method for the construction of the configuration
// object necessary for the connection to the host.
func (s *Session) NewConfig(user user) error {

	// get the content of the private key file
	key, err := loadPEM(os.ExpandEnv(tools.Expanduser(user.GetPrivateKey())))
	if err != nil {
		return err
	}

	// parse the key
	parsed, err := ssh.ParseRawPrivateKey(key)
	if err != nil {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"Can't parse the private key!\n",
			"Reason is: ", err.Error(),
		)

	}

	// convert into signer
	signer, err := ssh.NewSignerFromKey(parsed)
	if err != nil {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"Can't create signer from private key!\n",
			"Reason is: ", err.Error(),
		)

	}

	// Construct the configuration with password authentication
	s.Config = &ssh.ClientConfig{
		User: user.GetUsername(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return nil
}

// This function allows to connect to the host to create sessions on it
// after it.
func (s *Session) Connect(host host) error {

	var err error

	// create a new client for dialing with the host
	s.Client, err = ssh.Dial(
		host.GetProtocol(),
		net.JoinHostPort(
			host.GetHostname(),
			strconv.Itoa(host.GetPort()),
		),
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
func (s *Session) AddSession() error {

	// create the session
	session, err := s.Client.NewSession()

	// append the session to the list
	if err != nil {
		formatter.ColoredPrintln(
			formatter.Red,
			false,
			"Failed to create the session to the host!\n"+
				"Reason is: "+err.Error(),
		)
		return nil
	} else {
		s.Session = session
	}

	// return the result
	return err
}

// Close the last session created in the list.
func (s *Session) Close() error {
	// Close the session
	if s.Session != nil {
		err := s.Session.Close()
		if err != nil {
			return err
		}
	}
	if s.Client != nil {
		err := s.Client.Close()
		return err
	}
	return nil
}

// Run the command in argument using by default the last session
// created for this connection
func (s *Session) Run(command string) (string, error) {

	// Create a bytes buffer and affect it as Stdout writer
	// to return the output of the command
	var b bytes.Buffer

	// Affect the output to the buffer
	s.Session.Stdout = &b

	// run the command on the host
	err := s.Session.Run(command)

	// return the output and the error if one present
	return b.String(), err

}
