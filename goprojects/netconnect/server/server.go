package main

import (
	"netchan"	// Package for network channel-based communication
	"fmt"		// Package to allow printing to the console
	"os"		// Package to support the os.Error type
	"bufio"
	"strings"
	"time"
)

// Struct type for Data communication
type value struct {
	i int
	s string
}

var begin int64

func lt() {fmt.Print("[ ",time.Nanoseconds() / 1e6 - begin,"]: ")}

// Factory function to create listening communication channel that will wait for a connection from the client
func outChanFactory() (chan value, os.Error){
	//Create and initialize the export channel
	exp, err := netchan.NewExporter("tcp",":2345")
	lt();fmt.Println("exp.Addr().String(): ", exp.Addr().String())
	if err != nil {
		return nil, err
	}
	lt();fmt.Println("outgoing exportFactory Channel Made")
	//Make the communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the export channel
	err = exp.Export("exportChannelSend",ch,netchan.Send,new(value))
	if err != nil {
		return nil, err
	}
	return ch, err
}

func inChanFactory() (chan value, os.Error){
	//Create and initialize the export channel
	exp, err := netchan.NewExporter("tcp",":2346")
	lt();fmt.Println("exp.Addr().String(): ", exp.Addr().String())
	if err != nil {
		return nil, err
	}
	lt();fmt.Println("incoming exportFactory Channel Made")
	//Make the communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the export channel
	err = exp.Export("exportChannelRecv",ch,netchan.Recv,new(value))
	if err != nil {
		return nil, err
	}
	return ch, err
}
func printIncoming(inChan chan value, quit chan bool) {
	for ;  !closed(inChan) ; {
		inval := <- inChan
		lt();fmt.Println("Data recieved from server: ",inval)
	}
	quit <- true
}

func acceptOutgoing(outChan chan value, quit chan bool) {
	input := bufio.NewReader(os.Stdin)
	for i:= 0 ; !closed(outChan); i ++ {
		result, _ := input.ReadString('\n')
		if closed(outChan) { quit <- true }
		text := strings.TrimSpace(result)  
		outval := value{i,text}
		lt();fmt.Println("Sending data to outChan: ",outval)
		outChan <- outval
		lt();fmt.Println("Data sent")
	}
	quit <- true
}
func main() {
	//initialize program time
	begin = time.Nanoseconds() / 1e6
	//Create a new [value] variable
	//Use the factory to create a new network communication channel
	outChan, outerr := outChanFactory()
	if outerr != nil {
		lt();fmt.Println("outChanFactory error: ",outerr)
	}
	lt();fmt.Println("Returned to main after outChan made")
	inChan, inerr := inChanFactory()
	if inerr != nil {
		lt();fmt.Println("inChanFactory error: ",inerr)
	}
	lt();fmt.Println("Returned to main after inChan made")
	
	//Write all data recieved from stardard input to the network channel
	inQuit := make(chan bool)
	outQuit := make(chan bool)
	go printIncoming(inChan,inQuit)
	go acceptOutgoing(outChan, outQuit)
	switch {
		case <- inQuit:lt();fmt.Println("Incoming Channel quit")
		case <- outQuit:lt();fmt.Println("Outgoing Channel quit")
	}
	os.Exit(1)
	
}
