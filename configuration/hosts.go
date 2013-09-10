package configuration

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

// A type defining the user structure in the YAML file which we need
// in argument.
type MyHosts struct {
	Port     int "port"
	Priority int "priority"
	Threads  int "threads"
}

// The map type passed to the read YAML procedure to get multiple users
// informations if we need.
type Hosts map[string]MyHosts

// A function to read an hosts file in the YAML format and returns
// a dictionary in the same format as the structured file.
func ReadHostsYAML(
	filename string,
) *Hosts {

	// Start by reading the whole file in byte
	data, _ := ioutil.ReadFile(filename)

	// Create the variable handling the type of the user file
	t := &Hosts{}

	// Now read in the YAML file the structure of the file into
	// the structured dictionary
	err := goyaml.Unmarshal(
		data,
		t,
	)

	// Check error when reading the file
	if err != nil {
		panic(err)
	}

	// return the structured file and data
	return t

}
