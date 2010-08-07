package main

import (
	"netchan"	// Package for network channel-based communication
	"fmt"		// Package to allow printing to the console
	"os"		// Package to support the os.Error type
	"time"		// Package to support time.Sleep command
	"strings"
	"bufio"
)

// Struct type for Data communication
type value struct {
	i int
	s string
}

var begin int64

func lt() {fmt.Print("[ ",time.Nanoseconds() / 1e6 - begin,"]: ")}

// Factory function to create communication channel with server
func inChanFactory() (chan value, os.Error){
	//Create and initialize the import channel
	imp, err := netchan.NewImporter("tcp","127.0.0.1:2345")
	lt();fmt.Println("New incoming Connection: 127.0.0.1:2345")
	if err != nil {
		return nil, err
	}
	lt();fmt.Println("incoming importFactory Channel Made")
	//Make the incoming communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the import channel
	err = imp.Import("exportChannelSend",ch,netchan.Recv,new(value))
	if err != nil {
		return nil, err
	}
	return ch, err
}
func outChanFactory() (chan value, os.Error){
	//Create and initialize the import channel
	imp, err := netchan.NewImporter("tcp","127.0.0.1:2346")
	lt();fmt.Println("New outgoing connection: 127.0.0.1:2346")
	if err != nil {
		return nil, err
	}
	lt();fmt.Println("outgoing importFactory Channel Made")
	//Make the outgoing communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the import channel
	err = imp.Import("exportChannelRecv",ch,netchan.Send,new(value))
	if err != nil {
		return nil, err
	}
	return ch, err	
}
func printIncoming(inChan chan value, done chan bool) {
	for ; !closed(inChan) ; {
		inval := <- inChan
		lt();fmt.Println("Data recieved from server: ",inval)
	}
	done <- true
}

func acceptOutgoing(outChan chan value, done chan bool) {
	input := bufio.NewReader(os.Stdin)
	for i:= 0 ; !closed(outChan); i ++ {
		result, _ := input.ReadString('\n')
		text := strings.TrimSpace(result)
		outval := value{i,text}
		lt();fmt.Println("Sending data to outChan: ",outval)
		outChan <- outval
		lt();fmt.Println("Data sent")
	}
	done <- true
}
func comRoutine() {
	//Use the factory to create a new network communication channel, in and out
	inChan, inerr := inChanFactory()
	outChan, outerr := outChanFactory()
	if inerr != nil {
		lt();fmt.Println("inChanFactory error: ",inerr)
		return
	}else{
		lt();fmt.Println("Returned to main after inChan made")
	}
	if outerr != nil {
		lt();fmt.Println("outChanFactory error: ",outerr)
		return
	}else{
		lt();fmt.Println("Returned to main after outChan made")
	}
	
	inDone := make(chan bool)
	outDone := make(chan bool)
	go printIncoming(inChan,inDone)
	go acceptOutgoing(outChan,outDone)
	select {
		case <- inDone:lt();fmt.Println("Incoming Channel quit")
		case <- outDone:lt();fmt.Println("Outgoing Channel quit")
	}
	lt();fmt.Println("Killing other thread")	
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
