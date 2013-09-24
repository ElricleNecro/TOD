package timers

import (
	"github.com/ElricleNecro/TOD/formatter"
	"strconv"
	"time"
)

// The function timer which displays results all three minutes
func WorkingTimer(hosts []*formatter.Host, step int) {

	// set a timer
	timer := time.NewTicker(time.Duration(step) * time.Second)

	// infinite loop
	for {

		// wait for timer
		<-timer.C

		// get the current time
		hour, minute, second := time.Now().Clock()
		formatter.ColoredPrintln(
			formatter.Blue,
			false,
			"At time",
			strconv.Itoa(hour)+":"+
				strconv.Itoa(minute)+":"+strconv.Itoa(second),
		)

		// display working hosts
		DisplayWorkingHosts(hosts)

	}

}

// A routine to check the hosts which are executing commands at this time
func DisplayWorkingHosts(hosts []*formatter.Host) {

	// loop over hosts
	for _, host := range hosts {

		// if the host is connected and is marked as executing command
		if host.IsConnected && host.IsWorking {

			// Display the name of the host and the command being executed
			formatter.ColoredPrintln(
				formatter.Blue,
				false,
				host.Hostname, ":", host.Commands[host.CommandNumber-1].Command,
			)
			// Display the number of remaining commands
			formatter.ColoredPrintln(
				formatter.Blue,
				false,
				"Command",
				strconv.Itoa(host.CommandNumber)+"/"+
					strconv.Itoa(len(host.Commands)),
			)

		}
	}
}
