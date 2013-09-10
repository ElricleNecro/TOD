package configuration

import (
	"fmt"
	"testing"
)

// Test to read the users file in the YAML format.
func TestReadUsersYAML(t *testing.T) {

	// The resulting structure we must read in the example file
	tmp := MyUsers{
		Password: "potiron",
		Commands: []string{"hostname"},
	}
	tmp2 := MyUsers{
		Password: "Sloubi",
		Commands: []string{"Cépafo", "C'estpasfaux!"},
	}
	True := Users{
		"manuel":   tmp,
		"perceval": tmp2,
	}

	// Test read the file
	data := ReadUsersYAML("../examples/users/users_example.yaml")
	fmt.Println(True)
	fmt.Println(*data)
}
