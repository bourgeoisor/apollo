package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

// Loop is the main loop of the application. It creates a new goroutine to fetch new events
// and forwards all of them to Apollo. The screen is redrawn after each event.
func (a *Apollo) loop() {
	go func() {
		for {
			a.events <- termbox.PollEvent()
		}
	}()
	for a.running {
		select {
		case ev := <-a.events:
			err := a.handleEvent(&ev)
			a.draw()
			if err != nil {
				a.running = false
				log.Fatal(err)
			}
		}
	}
}

// Main function of the application. Initializes termbox, creates a new Apollo,
// and calls the main loop.
func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputAlt)

	apollo := newApollo()
	apollo.draw()
	apollo.loop()
}
