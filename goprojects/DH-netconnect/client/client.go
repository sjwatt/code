package main

import (
	"netchan"	// Package for network channel-based communication
	"fmt"		// Package to allow printing to the console
	"os"		// Package to support the os.Error type
	"time"		// Package to support time.Sleep command
	"strings"
	"bufio"
	"rand"
)

// Struct type for Data communication
type value struct {
	i int
	s string
}

var begin int64
var sharedSecret int
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
	err = imp.Import("exportChannelSend",ch,netchan.Recv)
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
	err = imp.Import("exportChannelRecv",ch,netchan.Send)
	if err != nil {
		return nil, err
	}
	return ch, err	
}

func printIncoming(inChan chan value, done chan bool) {
	for ; !closed(inChan) ; {
		inval := <- inChan
		lt();fmt.Println("Data recieved from server: ",inval)
		fmt.Println(decrypt(inval.s))
	}
	done <- true
}

func acceptOutgoing(outChan chan value, done chan bool) {
	input := bufio.NewReader(os.Stdin)
	for i:= 0 ; !closed(outChan); i ++ {
		result, _ := input.ReadString('\n')
		text := strings.TrimSpace(result)
		outval := value{i,encrypt(text)}
		lt();fmt.Println("Sending data to outChan: ",outval)
		outChan <- outval
		lt();fmt.Println("Data sent")
	}
	done <- true
}

func keyExchangeClient(outChan chan value, inChan chan value) (success bool) {

	goodConnection := false
	//receive Prime p from server
	inPrime := <- inChan
	lt();fmt.Println(inPrime)
	//receive base g from server
	inBase := <- inChan
	lt();fmt.Println(inBase)
	//receive Server's public key from server
	inServerPubKey := <- inChan
	lt();fmt.Println(inServerPubKey)
	
	//Generate secret integer b
	randMax := 100000000
	rand.Seed(time.Nanoseconds())
	secretInt := rand.Int() % randMax
	lt();fmt.Println("Secret Integer:",secretInt)

	//Generate Client Public key
	//g^b mod p
	PubKey := (inBase.i^secretInt)%(inPrime.i)
	lt();fmt.Println("Client Public Key:",PubKey)
	
	//Send Client PubKey to server
	outChan <- value{PubKey,"Client Public Key"}
	
	//Calculate Shared Secret
	//ClientPubKey^SecretInt mod prime
	sharedSecret = (inServerPubKey.i ^ secretInt) % inPrime.i
	lt();fmt.Println("Calculated Shared Secret:",sharedSecret)
	
	
	
	//Answer Server's Challenge
	//response = challenge^sharedSecret mod prime
	inChallenge := <- inChan
	lt();fmt.Println(inChallenge)
	response := (inChallenge.i ^ sharedSecret) % inPrime.i
	lt();fmt.Println("Response to Server:",response)
	outChan <- value{response,"Response From Client"}
	
	//Send encrypted challenge to Server
	//challenge formula: nonce^sharedSecret mod prime
	//Server: challenge^sharedSecret mod prime
	//Client checks result
	challenge := ((rand.Int()%randMax)^sharedSecret) % inPrime.i
	lt();fmt.Println("Challenge to Server:",challenge)
	expectedResult := (challenge^sharedSecret) % inPrime.i
	lt();fmt.Println("Expected Result:",expectedResult)
	outChan <- value{challenge,"Challenge From Client"}
	result := <- inChan
	lt();fmt.Println(result)
	if result.i == expectedResult {
		goodConnection = true
	}
	return goodConnection
}

//simple encryption routine
func encrypt(input string) (result string){
	data := []int(input)
	for i := 0; i < len(data); i ++ {
		data[i] += sharedSecret % 26
	}
	result = string(data)
	return result
}

//simple decryption routine
func decrypt(input string) (result string){
	data := []int(input)
	for i := 0; i < len(data); i ++ {
		data[i] -= sharedSecret % 26
	}
	result = string(data)
	return result
}
func comRoutine() {
	//Use the factory to create a new network communication channel, in and out
	inChan, inerr := inChanFactory()
	
	if inerr != nil {
		lt();fmt.Println("inChanFactory error: ",inerr)
		return
	}else{
		lt();fmt.Println("Returned to main after inChan made")
	}
	outChan, outerr := outChanFactory()
	if outerr != nil {
		lt();fmt.Println("outChanFactory error: ",outerr)
		return
	}else{
		lt();fmt.Println("Returned to main after outChan made")
	}

	if keyExchangeClient(outChan, inChan) {
		lt();fmt.Println("Keyexchange Success")
		inDone := make(chan bool)
		outDone := make(chan bool)
		go printIncoming(inChan,inDone)
		go acceptOutgoing(outChan,outDone)
		select {
			case <- inDone:lt();fmt.Println("Incoming Channel quit")
			case <- outDone:lt();fmt.Println("Outgoing Channel quit")
		}
		lt();fmt.Println("Killing other thread")
	} else {
		lt();fmt.Println("Keyexchange Failed")
	}
	time.Sleep(1e9)
}
func main() {
	//initialize program time
	begin = time.Nanoseconds() / 1e6
	for {
		// Attempt to form a com channel, and run if siccessful
		comRoutine()
		// If com channel quits/fails then wait 1 second and run again, forever.
		time.Sleep(2e9)
	}
}
