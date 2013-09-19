package load_checker

import (
	"github.com/ElricleNecro/TOD/commands"
	"github.com/ElricleNecro/TOD/formatter"
	"strconv"
	"strings"
)

// This function checks if the host has a too high charge and if it
// cans accept an other charge.
func IsTooLoaded(
	host *formatter.Host,
	user *formatter.User,
	timeout int,
	cpu_max float64,
	memory_max float64,
	disconnected chan<- *formatter.Host,
) bool {

	// Get the number of processors of the host
	nprocs := GetPhysicalCPU(
		host,
		user,
		timeout,
		disconnected,
	)

	// get the list of users
	users := GetUsers(
		host,
		user,
		timeout,
		disconnected,
	)

	// get the total CPU and memory used on the host
	CPU, memory := GetTotalCPUMemory(
		host,
		user,
		timeout,
		disconnected,
		users,
	)

	// return if the host is too loaded
	return CPU/float64(nprocs) >= cpu_max || memory >= memory_max
}

// Function to get the total CPU and memory used.
func GetTotalCPUMemory(
	host *formatter.Host,
	user *formatter.User,
	timeout int,
	disconnected chan<- *formatter.Host,
	users []string,
) (float64, float64) {

	// loop over users
	cpu := 0.0
	mem := 0.0
	for _, myuser := range users {

		// get the cpu and memory
		stats := GetCPUMemory(
			host,
			user,
			timeout,
			disconnected,
			myuser,
		)
		cpu += stats[0]
		mem += stats[1]

	}

	// return the total
	return cpu, mem

}

// Get the CPU and memory used by an user on a host.
func GetCPUMemory(
	host *formatter.Host,
	user *formatter.User,
	timeout int,
	disconnected chan<- *formatter.Host,
	username string,
) [2]float64 {

	// the command to get the cpu and memory of an user
	command := "top -b -n 1 -u " + strings.TrimSpace(username) +
		" | awk 'NR>7 { cpu += $9 ; mem += $10 } END { print cpu, mem; }'"

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

	// return the slice of CPU and memory.
	return ParseCPUMemory(output)

}

// Function to parse the output of the cpu and memory.
func ParseCPUMemory(output string) [2]float64 {

	// split by space
	split := strings.Fields(output)

	// convert the values and return in a slice
	cpu, _ := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
	mem, _ := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
	return [2]float64{
		cpu,
		mem,
	}

}

// A function to get the list of users connected to an host.
func GetUsers(
	host *formatter.Host,
	user *formatter.User,
	timeout int,
	disconnected chan<- *formatter.Host,
) []string {

	// the command to execute in order to get the list of users
	command := "who | cut -f 1 -d' ' | uniq"

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
	return ParseUsers(output)

}

// This function parses the output of users to transform it
// on a slice of users.
func ParseUsers(output string) []string {

	// split the string by returns
	return strings.Split(output, "\n")

}

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
	res, _ := strconv.Atoi(strings.TrimSpace(output))
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

	// format the strings for the different platforms results
	// Sometimes its written 0.1 and other 0,1 .
	nb := len(columns)
	c1 := strings.Replace(strings.TrimSpace(columns[nb-3]), ",", ".", 1)
	c5 := strings.Replace(strings.TrimSpace(columns[nb-2]), ",", ".", 1)
	c15 := strings.Replace(strings.TrimSpace(columns[nb-1]), ",", ".", 1)

	// return the number
	r1, _ := strconv.ParseFloat(c1[:len(c1)-1], 64)
	r2, _ := strconv.ParseFloat(c5[:len(c5)-1], 64)
	r3, _ := strconv.ParseFloat(c15, 64)
	return r1, r2, r3

}