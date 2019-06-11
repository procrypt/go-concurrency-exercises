//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	signalStream := make(chan os.Signal)
	signal.Notify(signalStream, os.Interrupt)

	// Create a process
	proc := MockProcess{}
	count := 0
	go func() {
		for {
			select {
			case <-signalStream:
				count++
				if count == 2 {
					fmt.Println("Forcefully Exiting")
					os.Exit(1)
				}
				go func() {
					proc.Stop()
				}()
			}
		}
	}()

	// Run the process (blocking)
	proc.Run()
}
