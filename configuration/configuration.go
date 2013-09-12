package configuration

import (
	"flag"
)

// The type containing the informations and parameters of the program.
type Config struct {

	// The file name for the users and associated commands
	Users *string

	// The file name for the hosts on which to execute commands
	Hosts *string

	// To read on stdin
	Stdin *bool

	// The default protocol for hosts
	Protocol *string

	// The default port on hosts
	Port *int

	// The timeout in seconds for the disconnection
	Timeout *int
}

// A function to get the data from the command line
// and store it correctly in the datastructure.
func ReadConfig() *Config {

	// Create a configuration object
	data_config := &Config{}

	// define flag to use in the command line
	data_config.Users = flag.String(
		"users",
		"",
		"The path to the file where users and associated commands are stored.",
	)
	data_config.Hosts = flag.String(
		"hosts",
		"",
		"The path to the file where hosts are stored.",
	)
	data_config.Stdin = flag.Bool(
		"stdin",
		false,
		"The list of hosts blank separated on which to run commands.",
	)
	data_config.Protocol = flag.String(
		"protocol",
		"tcp",
		"The protocol used by default by hosts to communicate.",
	)
	data_config.Port = flag.Int(
		"port",
		22,
		"The port used by default by hosts to listen for SSH connection.",
	)
	data_config.Timeout = flag.Int(
		"timeout",
		10,
		"The default time out in second to wait before to say that the host"+
			" is disconnected.",
	)

	// parse the command line
	flag.Parse()

	// return the data structure
	return data_config

}
