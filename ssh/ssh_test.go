package ssh

import (
	"github.com/ElricleNecro/TOD/formatter"
	"testing"
)

func TestConnection(t *testing.T) {

	// Start by creating user and host object
	user := formatter.User{
		Name:     "duarte",
		Identity: 1,
		Password: "SYmadu10;",
	}
	host := formatter.Host{
		Hostname: "carmenere",
		Port:     22,
		Protocol: "tcp",
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
	}

	// display
	t.Log("All test for connection done !")

}
