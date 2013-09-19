package load_checker

import (
	"github.com/ElricleNecro/TOD/commands"
	"github.com/ElricleNecro/TOD/formatter"
	"strconv"
	"strings"
)

// Function to get the number of physical cpu in the host.
func GetPhysicalCPU(
	host *formatter.Host,
	user *formatter.User,
	timeout int,
	disconnected chan<- *formatter.Host,
) int {

	// the command to get the number of physical CPU
	command := "cat /proc/cpuinfo | grep 'processor' | uniq | wc -l"

	// run this command on the host
	output, err := commands.OneCommand(
		host,
		user,
		command,
		timeout,
		disconnected,
	)

	// check can connect
	if err != nil {
		panic(err)
	}

	// return the number
	return ParseCPU(output)
}

// Function to parse the output of the command returning the number of processors
// available in the host.
func ParseCPU(output string) int {

	// Return the result of the conversion of the string to value
	res, _ := strconv.Atoi(output)
	return res
}

// Function to get the three values of the uptime.
func GetLoadAverage(
	host *formatter.Host,
	user *formatter.User,
	timeout int,
	disconnected chan<- *formatter.Host,
) (float64, float64, float64) {

	// the command to get the number of physical CPU
	command := "uptime"

	// run this command on the host
	output, err := commands.OneCommand(
		host,
		user,
		command,
		timeout,
		disconnected,
	)

	// check can connect
	if err != nil {
		panic(err)
	}

	// Return the result
	return ParseUptime(output)
}

// This function parses the string in output of the command uptime
// in order to extract the load average values of the host.
func ParseUptime(output string) (float64, float64, float64) {

	// split the output by space
	columns := strings.Split(output, " ")

	// return the number
	nb := len(columns)
	r1, _ := strconv.ParseFloat(columns[nb-1], 64)
	r2, _ := strconv.ParseFloat(columns[nb-2], 64)
	r3, _ := strconv.ParseFloat(columns[nb-3], 64)
	return r1, r2, r3

}
