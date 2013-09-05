package exec

import (
	"log"
	"os/exec"
)

//Launch the task in a thread and wait for the command to terminate. Return a channel for getting the result state.
// False is send to the channel and a message is write in the Logger if their was was a problem, true elsewhere.
func LaunchTask(task *exec.Cmd, tlog *log.Logger) <-chan bool {
	c := make(chan bool)
	go func() {
		if err := task.Run(); err != nil {
			if tlog != nil {
				tlog.Println(err)
			}
			c <- false
		} else {
			c <- true
		}
	}()

	return c
}

//vim: spelllang=en
