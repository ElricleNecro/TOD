package host

import (
	"net"
	"strconv"
	"time"

	"github.com/ElricleNecro/TOD/commands"
	"github.com/ElricleNecro/TOD/formatter"
	"github.com/ElricleNecro/TOD/ssh"
)

// the maximal value of the integer 32
const (
	MaxInt = 1<<31 - 1
)

// interface for users
type user interface {
	GetPrivateKey() string
	GetUsername() string
}

type Host struct {

	// The list of commands to execute on the host
	Commands []commands.Command

	// the timeout for the connection to this hosts
	Timeout int

	// the name of the host
	Hostname string

	// the protocol for the communication
	Protocol string

	// The port for the connection to the host
	Port int

	// The command number being executed
	CommandNumber int

	// Store here if the host is connected or not
	Connected bool

	// Store if the host is working
	IsWorking bool

	// Channel on which to wait for new job
	Wait *(chan int)

	// the session for the conenction
	session *ssh.Session

	// Channel to pass information for disconnection
	disconnected *(chan<- *Host)
}

func (host *Host) GetProtocol() string {
	return host.Protocol
}

func (host *Host) GetPort() int {
	return host.Port
}

func (host *Host) GetHostname() string {
	return host.Hostname
}

// Function which check that an host is connected or not by checking errors
// when attempting to connect to it and by setting a timer for the connection
// timeout if nothing is responding.
func (host *Host) IsConnected() (bool, error) {

	// create a dialer
	dial := net.Dialer{
		Deadline:  time.Now().Add(time.Duration(host.Timeout) * time.Second),
		Timeout:   time.Duration(host.Timeout) * time.Second,
		LocalAddr: nil,
	}

	// Contact the host and if no error, it is connected
	_, err := dial.Dial(
		host.Protocol,
		net.JoinHostPort(host.Hostname, strconv.Itoa(host.Port)),
	)

	return err == nil, err
}

// This function returns true if the host in argument is considered as
// worker.
func (host *Host) IsWorker() bool {
	return host.Connected && len(host.Commands) > 0
}

// Create a session for the host
func (host *Host) CreateSession(user user) (string, error) {

	var err error

	// create a session object
	host.session, err = ssh.New(user)
	if err != nil {
		return "The session for the connection can't be created!\n" +
			"Reason is: " + err.Error(), err
	}

	// check the connection to the host
	if is, err := host.IsConnected(); !is || (err != nil) {
		// exit the loop
		return "Can't connect to host " + host.Hostname + "\n" +
			"Reason is: " + err.Error(), err
	}

	// Attempt a connection to the host
	err = host.session.Connect(host)

	// check the host can be called
	if err != nil {
		// exit the loop
		return "Can't create connection to host " + host.Hostname + "\n" +
			"Reason is: " + err.Error(), err
	}

	// add a session to connect to host
	_, err = host.session.AddSession()
	if err != nil {
		// Close the session
		host.session.Close()

		// exit the loop
		return "Problem when adding a session to the host!\n" +
			"Reason is: " + err.Error(), err
	}

	return "", nil
}

// This function runs a single command on the host on argument.
func (host *Host) OneCommand(command *commands.Command) (string, error) {

	// create a new session for the host
	message, err := host.CreateSession(command.User)
	if err != nil {
		return message, err
	}
	defer host.session.Close()

	// execute the command on the host
	output, err := host.session.Run(command.Command)
	if err != nil {
		// exit the loop
		return "An error occurred during the execution of the command !\n" +
			"The command was: " + command.Command +
			"\nand the host is: " + host.Hostname +
			"\nError information: " + err.Error(), err
	}

	// return nil if good
	return output, nil
}

// Function which executes commands when a host has to wait for other hosts.
func (host *Host) Waiter() {

	// display
	formatter.ColoredPrintln(
		formatter.Magenta,
		true,
		"Waiting more jobs for", host.Hostname,
	)

	// say it is not working
	host.IsWorking = false

	// Now wait for new job
	<-*(host.Wait)

	// display
	formatter.ColoredPrintln(
		formatter.Magenta,
		true,
		host.Hostname, "has more jobs !",
	)
	formatter.ColoredPrintln(
		formatter.Green,
		true,
		"Number of commands for", host.Hostname, ":",
		len(host.Commands),
	)

	// indicate that it is working
	host.IsWorking = true

}
