package configuration

import (
	"fmt"
	"testing"
)

// Test to read the users file in the YAML format.
func TestReadHostsYAML(t *testing.T) {

	// Test read the file
	data := ReadHostsYAML("../examples/hosts/hosts_examples.yaml")
	fmt.Println(*data)
}
