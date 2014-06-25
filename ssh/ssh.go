// The package containing simplified interface for the SSH
// library.
package ssh

import (
	"bytes"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"

	"crypto"
	"crypto/dsa"
	"crypto/rsa"
	"crypto/x509"

	"code.google.com/p/go.crypto/ssh"
	"github.com/ElricleNecro/TOD/formatter"
)

// A structure containing all information necessary to an SSH connection
type Session struct {

	// The configuration for the session
	Config *ssh.ClientConfig

	// The structure of the client
	Client *ssh.ClientConn

	// The structure of the session and their number
	Session   []*ssh.Session
	nsessions int

	// The user associated
	User *formatter.User

	// The host associated
	Host *formatter.Host
}

// keychain implements the ClientKeyring interface
type keychain struct {
	keys []interface{}
}

func (k *keychain) Key(i int) (interface{}, error) {
	if i < 0 || i >= len(k.keys) {
		return nil, nil
	}
	switch key := k.keys[i].(type) {
	case *rsa.PrivateKey:
		return &key.PublicKey, nil
	case *dsa.PrivateKey:
		return &key.PublicKey, nil
	}
	panic("unknown key type")
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	hashFunc := crypto.SHA1
	h := hashFunc.New()
	h.Write(data)
	digest := h.Sum(nil)
	switch key := k.keys[i].(type) {
	case *rsa.PrivateKey:
		return rsa.SignPKCS1v15(rand, key, hashFunc, digest)
	}
	return nil, errors.New("ssh: unknown key type")
}

func (k *keychain) loadPEM(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		return errors.New("ssh: no key found")
	}
	r, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	k.keys = append(k.keys, r)
	return nil
}

// expanduser
func expanduser(path string) string {
	usr, _ := user.Current()
	home := usr.HomeDir
	if path[:1] == "~" {
		path = strings.Replace(path, "~", home, 1)
	}
	return path
}

// A method for the construction of the configuration
// object necessary for the connection to the host.
func (s *Session) NewConfig(
	user *formatter.User,
) error {

	k := new(keychain)
	k.loadPEM(os.ExpandEnv(expanduser(user.Key)))

	// Construct the configuration with password authentication
	s.Config = &ssh.ClientConfig{
		User: user.Name,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(k),
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

	// init the number of sessions
	session.nsessions = 0

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
		net.JoinHostPort(
			s.Host.Hostname,
			strconv.Itoa(s.Host.Port),
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
func (s *Session) AddSession() (*ssh.Session, error) {

	// create the session
	session, err := s.Client.NewSession()

	// append the session to the list
	if err != nil {
		panic("Failed to create the session to the host !")
	} else {
		s.Session = append(s.Session, session)
		s.nsessions++
	}

	// return the result
	return session, err

}

// Close the last session created in the list.
func (s *Session) Close() error {

	// Close the session
	s.Session[s.nsessions-1].Close()

	// return
	return nil
}

// Run the command in argument using by default the last session
// created for this connection
func (s *Session) Run(command string) (string, error) {

	// Create a bytes buffer and affect it as Stdout writer
	// to return the output of the command
	var b bytes.Buffer

	// get the good session
	session := s.Session[s.nsessions-1]

	// Affect the output to the buffer
	session.Stdout = &b

	// run the command on the host
	err := session.Run(command)

	// return the output and the error if one present
	return b.String(), err

}
