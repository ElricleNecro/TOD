package host

import (
	"strconv"
	"time"

	"github.com/ElricleNecro/TOD/formatter"
)

// The function timer which displays results all three minutes
func (hosts *Hosts) WorkingTimer(step int) {

	// set a timer
	timer := time.NewTicker(time.Duration(step) * time.Second)

	// infinite loop
	for {

		// wait for timer
		<-timer.C

		// get the current time

		formatter.ColoredPrintln(
			formatter.Magenta,
			false,
			"#################################################",
		)
		formatter.ColoredPrintln(
			formatter.Blue,
			false,
			"At time ",
			time.Now().Format("15:04:05"),
		)

		// display working hosts
		hosts.DisplayWorkingHosts()
	}
}

// A routine to check the hosts which are executing commands at this time
func (hosts *Hosts) DisplayWorkingHosts() {

	// loop over hosts
	for _, host := range hosts.Hosts {
		// if the host is connected and is marked as executing command
		if host.Connected && host.IsWorking {
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
