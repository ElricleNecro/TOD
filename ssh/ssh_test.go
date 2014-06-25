package ssh

import (
	"fmt"
	"testing"

	"github.com/ElricleNecro/TOD/formatter"
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

func TestConnection(t *testing.T) {

	// Start by creating user and host object
	user := formatter.User{
		Name:     myuser,
		Identity: id,
		Key:      key,
	}
	host := formatter.Host{
		Hostname: myhost,
		Port:     port,
		Protocol: protocol,
	}
	t.Log(host)

	// display
	t.Log("Creation of user and host done !")

	// Now create a session object
	session := New(
		&user,
		&host,
	)

	// display
	t.Log("Creation of the session done !")

	// Attempt a connection to the host
	err := session.Connect()
	if err != nil {
		t.Errorf("Can't connect to the host specified in the test !")
		t.FailNow()
	}

	// display
	t.Log("All test for connection done !")

}

func TestAddSession(t *testing.T) {

	// Start by creating user and host object
	user := formatter.User{
		Name:     myuser,
		Identity: id,
		Key:      key,
	}
	host := formatter.Host{
		Hostname: myhost,
		Port:     port,
		Protocol: protocol,
	}

	// display
	t.Log("Creation of user and host done !")

	// Now create a session object
	session := New(
		&user,
		&host,
	)

	// display
	t.Log("Creation of the session done !")

	// Attempt a connection to the host
	err := session.Connect()
	if err != nil {
		t.Errorf("Can't connect to the host specified in the test !")
		t.FailNow()
	}

	// display
	t.Log("Connection done !")

	// Run a command for test
	_, err2 := session.AddSession()

	// test
	if err2 != nil {
		t.Errorf("Can't add a session to the connected host !")
		t.FailNow()
	}

	// display
	t.Log("All done to add sessions !")

}

func TestRun(t *testing.T) {

	// Start by creating user and host object
	user := formatter.User{
		Name:     myuser,
		Identity: id,
		Key:      key,
	}
	host := formatter.Host{
		Hostname: myhost,
		Port:     port,
		Protocol: protocol,
	}

	// display
	t.Log("Creation of user and host done !")

	// Now create a session object
	session := New(
		&user,
		&host,
	)

	// display
	t.Log("Creation of the session done !")

	// Attempt a connection to the host
	err := session.Connect()
	if err != nil {
		t.Errorf("Can't connect to the host specified in the test !")
		t.FailNow()
	}

	// display
	t.Log("Connection done !")

	// Run a command for test
	_, err2 := session.AddSession()

	// test
	if err2 != nil {
		t.Errorf("Can't add a session to the connected host !")
		t.FailNow()
	}

	// display
	t.Log("Adding a session is done !")

	// Now run a command to see the result
	output, err3 := session.Run(mycommand)
	if err3 != nil {
		t.Errorf("Can't run a simple command on the host !")
		t.Errorf("The error is: " + err3.Error())
		t.FailNow()
	}

	// close the session we have opened
	session.Close()

	// display the result of the command
	// TODO: add a assertion for the expected result after running this command
	fmt.Println(output)

}
