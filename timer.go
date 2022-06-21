package main

import (
	"time"
)

func setInterval(someFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {

			select {
			case <-ticker.C:
				if async {
					// This won't block
					go someFunc()
				} else {
					// This will block
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				// return
			}

		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear

}

// func main() {

// A counter for the number of times we print

// We call set interval to print Hello World forever
// every 1 second
// clear :=setInterval(func() {
// 	fmt.Println("Hello World")
// 	printed++
// }, 1000, false)

// If we wanted to we had a long running task (i.e. network call)
// we could pass in true as the last argument to run the function
// as a goroutine

// Some artificial work here to wait till we've printed
// 5 times
// for {
// 	if printed == 1 {
// Stop the ticket, ending the interval go routine
// 		clear <- true
// return
// 	}
// }
// }
