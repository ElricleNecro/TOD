package ssh

import (
	"fmt"
	"testing"
)

var (
	myuser    = "duarte"
	myhost    = "carmenere"
	port      = 22
	protocol  = "tcp"
	id        = 1
	key       = "~/.ssh/id_rsa"
	mycommand = "/usr/bin/whoami"
)

type User struct {
	Username   string
	PrivateKey string
}

func (user *User) GetUsername() string {
	return user.Username
}

func (user *User) GetPrivateKey() string {
	return user.PrivateKey
}

type Host struct {
	Hostname string
	Port     int
	Protocol string
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

func TestConnection(t *testing.T) {

	// Start by creating user and host object
	user := User{
		Username:   myuser,
		PrivateKey: key,
	}
	host := Host{
		Hostname: myhost,
		Port:     port,
		Protocol: protocol,
	}
	t.Log(host, user)

	// display
	t.Log("Creation of user and host done!")

	// Now create a session object
	session, err := New(&user)
	if err != nil {
		t.Errorf("Can't create the session to the host specified in " +
			"the test!\nReason is " + err.Error())
		t.FailNow()
	}

	// display
	t.Log("Creation of the session done!")

	// Attempt a connection to the host
	err = session.Connect(&host)
	if err != nil {
		t.Errorf("Can't connect to the host specified in the test!\n" +
			err.Error())
		t.FailNow()
	}

	// display
	t.Log("All test for connection done!")

}

func TestAddSession(t *testing.T) {

	// Start by creating user and host object
	user := User{
		Username:   myuser,
		PrivateKey: key,
	}
	host := Host{
		Hostname: myhost,
		Port:     port,
		Protocol: protocol,
	}

	// display
	t.Log("Creation of user and host done!")

	// Now create a session object
	session, err := New(&user)
	if err != nil {
		t.Errorf("Can't create the session to the host specified in " +
			"the test!\nReason is " + err.Error())
		t.FailNow()
	}

	// display
	t.Log("Creation of the session done!")

	// Attempt a connection to the host
	err = session.Connect(&host)
	if err != nil {
		t.Errorf("Can't connect to the host specified in the test!\n" +
			err.Error())
		t.FailNow()
	}

	// display
	t.Log("Connection done!")

	// Run a command for test
	err = session.AddSession()

	// test
	if err != nil {
		t.Errorf("Can't add a session to the connected host!\n" + err.Error())
		t.FailNow()
	}

	// display
	t.Log("All done to add sessions!")
}

func TestRun(t *testing.T) {

	// Start by creating user and host object
	user := User{
		Username:   myuser,
		PrivateKey: key,
	}
	host := Host{
		Hostname: myhost,
		Port:     port,
		Protocol: protocol,
	}

	// display
	t.Log("Creation of user and host done!")

	// Now create a session object
	session, err := New(&user)
	if err != nil {
		t.Errorf("Can't create the session to the host specified in " +
			"the test!\nReason is " + err.Error())
		t.FailNow()
	}

	// display
	t.Log("Creation of the session done!")

	// Attempt a connection to the host
	err = session.Connect(&host)
	if err != nil {
		t.Errorf("Can't connect to the host specified in the test!\n" +
			err.Error(),
		)
		t.FailNow()
	}

	// display
	t.Log("Connection done!")

	// Run a command for test
	err = session.AddSession()

	// test
	if err != nil {
		t.Errorf("Can't add a session to the connected host!\n" + err.Error())
		t.FailNow()
	}

	// display
	t.Log("Adding a session is done!")

	// Now run a command to see the result
	output, err3 := session.Run(mycommand)
	if err3 != nil {
		t.Errorf("Can't run a simple command on the host!\n")
		t.Errorf("The error is: " + err3.Error())
		t.FailNow()
	}

	// close the session we have opened
	session.Close()

	// display the result of the command
	// TODO: add a assertion for the expected result after running this command
	fmt.Println(output)
}
