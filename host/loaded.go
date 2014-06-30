package host

import (
	"strconv"
	"strings"

	"github.com/ElricleNecro/TOD/commands"
)

// interface for the parameters of the configuration
type config interface {
	GetCPUMax() float64
	GetMemoryMax() float64
	GetTimer() int
}

// This function checks if the host has a too high charge and if it
// cans accept an other charge.
func (host *Host) IsTooLoaded(
	user user,
	config config,
) (bool, error) {

	// Get the number of processors of the host
	nprocs, err := host.GetPhysicalCPU(user)
	if err != nil {
		return true, err
	}

	// get the list of users
	users, err := host.GetUsers(user)
	if err != nil {
		return true, err
	}

	// get the total CPU and memory used on the host
	CPU, memory, err := host.GetTotalCPUMemory(
		user,
		users,
	)
	if err != nil {
		return true, err
	}

	// return if the host is too loaded
	return CPU/float64(nprocs) >= config.GetCPUMax() ||
		memory >= config.GetMemoryMax(), nil
}

// Get the CPU and memory used by an user on a host.
func (host *Host) GetCPUMemory(
	user user,
	query_user user,
) ([2]float64, error) {

	// the command to get the cpu and memory of an user
	command := &commands.Command{
		Command: "top -b -n 1 -u " + strings.TrimSpace(
			query_user.GetUsername(),
		) +
			" | awk 'NR>7 { cpu += $9 ; mem += $10 } END { print cpu, mem; }'",
		User: user,
	}

	// run this command on the host
	output, err := host.OneCommand(command)

	// check can connect
	if err != nil {
		return [2]float64{0., 0.}, err
	}

	// return the slice of CPU and memory.
	return parseCPUMemory(output), err

}

// Function to parse the output of the cpu and memory.
func parseCPUMemory(output string) [2]float64 {

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
func (host *Host) GetUsers(
	user user,
) ([]string, error) {

	// the command to execute in order to get the list of users
	command := &commands.Command{
		Command: "ps haeo user | sort -u",
		User:    user,
	}

	// run this command on the host
	output, err := host.OneCommand(command)

	// check can connect
	if err != nil {
		return []string{}, err
	}

	// return the number
	return parseUsers(output), err

}

// This function parses the output of users to transform it
// on a slice of users.
func parseUsers(output string) []string {

	// split the string by returns
	return strings.Fields(output)

}

// Function to get the number of physical cpu in the host.
func (host *Host) GetPhysicalCPU(
	user user,
) (int, error) {

	// the command to get the number of physical CPU
	command := &commands.Command{
		Command: "cat /proc/cpuinfo | grep 'processor' | uniq | wc -l",
		User:    user,
	}

	// run this command on the host
	output, err := host.OneCommand(command)

	// check can connect
	if err != nil {
		return 0, err
	}

	// return the number
	return parseCPU(output), err
}

// Function to parse the output of the command returning the number of processors
// available in the host.
func parseCPU(output string) int {

	// Return the result of the conversion of the string to value
	res, _ := strconv.Atoi(strings.TrimSpace(output))
	return res
}

// define a structure for the users
type myUser struct {
	Username string
}

func (user *myUser) GetPrivateKey() string {
	return ""
}

func (user *myUser) GetUsername() string {
	return user.Username
}

// Function to get the total CPU and memory used.
func (host *Host) GetTotalCPUMemory(
	user user,
	users []string,
) (float64, float64, error) {

	// loop over users
	cpu := 0.0
	mem := 0.0
	for _, myuser := range users {

		usr := &myUser{
			Username: myuser,
		}
		// get the cpu and memory
		stats, err := host.GetCPUMemory(user, usr)
		if err != nil {
			return 0.0, 0.0, err
		}
		cpu += stats[0]
		mem += stats[1]

	}

	// return the total
	return cpu, mem, nil

}
