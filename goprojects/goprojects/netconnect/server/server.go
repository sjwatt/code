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
	close bool
}

var begin int64

func lt() {fmt.Print("[ ",time.Nanoseconds() / 1e6 - begin,"]: ")}

// Factory function to create listening communication channel that will wait for a connection from the client
func outChanFactory() (chan value, os.Error){
	//Create and initialize the import channel
	exp, err := netchan.NewExporter("tcp",":2345")
	lt();fmt.Println("exp.Addr().String(): ", exp.Addr().String())
	if err != nil {
		return nil, err
	}
	lt();fmt.Println("exportFactory Channel Made")
	//Make the communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the export channel
	err = exp.Export("exportChannel",ch,netchan.Send,new(value))
	if err != nil {
		return nil, err
	}
	return ch, err
}

func main() {
	//initialize program time
	begin = time.Nanoseconds() / 1e6
	//Create a new [value] variable
	testval := value{}
	input := bufio.NewReader(os.Stdin)
	//Use the factory to create a new network communication channel
	outChan, err := outChanFactory()
	if err != nil {
		lt();fmt.Println("outChanFactory error: ",err)
	}
	lt();fmt.Println("Returned to main after outChan made")
	//Write all data recieved from stardard input to the network channel
	for i:= 0 ; !testval.close && !closed(outChan); i ++ {
		result, _ := input.ReadString('\n')
		text := strings.TrimSpace(result)  
		testval = value{i,text,(text == "quit")}
		lt();fmt.Println("Sending data to outChan: ",testval)
		outChan <- testval
		lt();fmt.Println("Data sent")
	}
	lt();fmt.Println("Sending 2 Terminate signals to client")
	testval = value{0,"Goodbye!",true}
	outChan <- testval
	outChan <- testval
	lt();fmt.Println("Terminate signals sent - Exiting")
	
}
