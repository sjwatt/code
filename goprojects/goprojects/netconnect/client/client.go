package main

import (
	"netchan"	// Package for network channel-based communication
	"fmt"		// Package to allow printing to the console
	"os"		// Package to support the os.Error type
	"time"		// Package to support time.Sleep command
)

// Struct type for Data communication
type value struct {
	i int
	s string
	close bool
}

var begin int64

func lt() {fmt.Print("[ ",time.Nanoseconds() / 1e6 - begin,"]: ")}

// Factory function to create communication channel with server
func inChanFactory() (chan value, os.Error){
	//Create and initialize the import channel
	imp, err := netchan.NewImporter("tcp","127.0.0.1:2345")
	lt();fmt.Println("New Connection: 127.0.0.1:2345")
	if err != nil {
		return nil, err
	}
	lt();fmt.Println("importFactory Channel Made")
	//Make the communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the import channel
	err = imp.Import("exportChannel",ch,netchan.Recv,new(value))
	if err != nil {
		return nil, err
	}
	return ch, err
}

func comRoutine() {
	//Create a new [value] variable
	testval := value{0,"",false}
	//Use the factory to create a new network communication channel
	inChan, err := inChanFactory()
	if err != nil {
		lt();fmt.Println("inChanFactory error: ",err)
		testval = value{0,"",true}
	}else{
		lt();fmt.Println("Returned to main after inChan made")
	}
	for ; !testval.close && !closed(inChan); {
		//Read the data from the network communication channel and print it to the screen
		testval = <- inChan
		lt();fmt.Println("Data recieved from server: ",testval)
	}
}
func main() {
	//initialize program time
	begin = time.Nanoseconds() / 1e6
	for {
		// Attempt to form a com channel, and run if siccessful
		comRoutine()
		// If com channel quits/fails then wait 1 second and run again, forever.
		time.Sleep(1e9)
	}
}
