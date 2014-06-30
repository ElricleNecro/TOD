package host

import (
	"strconv"
	"testing"
)

// To test the result of parsing the output of the uptime command
// in order to get the load average of an host.
func TestParseUptime(t *testing.T) {

	// Some possible outputs of the uptime command
	outputs := []string{
		"15:51:31 up 6 days, 17:44,  4 users,  load average: 0.00, 0.01, 0.05",
		" 15:51:51 up 1 day, 19:31, 19 users,  load average: 0,28, 0,25, 0,23",
	}

	// Expected result
	trues := []float64{
		0.00,
		0.01,
		0.05,
		0.28,
		0.25,
		0.23,
	}

	// loop over outputs and parse them
	for i, output := range outputs {

		// parse the command
		r1, r5, r15 := ParseUptime(output)

		// Compare the results
		if r1 != trues[3*i] {
			t.Errorf("Error parsing the load for last minute !")
		}
		if r5 != trues[3*i+1] {
			t.Errorf("Error parsing the load for last five minutes !")
		}
		if r15 != trues[3*i+2] {
			t.Errorf("Error parsing the load for last fifteen minutes !")
		}
	}
}

// To test the result of the parsing of the number of processor.
func TestParseCPU(t *testing.T) {

	// write some commands
	commands := []string{
		"  3\n",
		"4",
		"\t8\n",
		"\t9 ",
	}

	// true values
	trues := []int{
		3,
		4,
		8,
		9,
	}

	// loop over commands and check the result
	for i, command := range commands {

		// parse the command
		res := ParseCPU(command)

		// check the result is good
		if res != trues[i] {

			t.Errorf(
				"The parser of CPU doesn't work !\nValue " +
					strconv.Itoa(res) + " for value " + strconv.Itoa(trues[i]),
			)
		}
	}
}
