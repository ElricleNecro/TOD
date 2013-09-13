package configuration

import (
	"flag"
)

// The type containing the informations and parameters of the program.
type Config struct {

	// The file name for the users and associated commands
	Users string

	// The file name for the hosts on which to execute commands
	Hosts string

	// To read on stdin
	Stdin bool

	// The default protocol for hosts
	Protocol string

	// The default port on hosts
	Port int

	// The timeout in seconds for the disconnection
	Timeout int

	// To set or not the timer for displaying remaining commands
	Timer bool
}

// A function to get the data from the command line
// and store it correctly in the datastructure.
func ReadConfig() *Config {

	// Create a configuration object
	data_config := &Config{}

	// define flag to use in the command line
	flag.StringVar(
		&data_config.Users,
		"users",
		"",
		"The path to the file where users and associated commands are stored.",
	)
	flag.StringVar(
		&data_config.Hosts,
		"hosts",
		"",
		"The path to the file where hosts are stored.",
	)
	flag.BoolVar(
		&data_config.Stdin,
		"stdin",
		false,
		"The list of hosts blank separated on which to run commands.",
	)
	flag.BoolVar(
		&data_config.Timer,
		"stdin",
		false,
		"If set, a timer will be launched to display the number of"+
			" remaining commands.",
	)
	flag.StringVar(
		&data_config.Protocol,
		"protocol",
		"tcp",
		"The protocol used by default by hosts to communicate.",
	)
	flag.IntVar(
		&data_config.Port,
		"port",
		22,
		"The port used by default by hosts to listen for SSH connection.",
	)
	flag.IntVar(
		&data_config.Timeout,
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
