/************************************************************
 * This is the client side of a simple (not secure)
 * implementation of the Diffie-hellman key exchange.
 * In this instance it is used to facilitate a command
 * line chat program that uses the key exchange to generate
 * a shared key for simple rotation cipher encryption
 * Author: Simon Watt sjwatt@gmail.com
 * Course: INFO3270
 * Instructor: Dr. Abhijit Sen
 * Filename: DHClient.go
 * Language: GO
 */
package main

import (
	"netchan" // Package for network channel-based communication
	"fmt"     // Package to allow printing to the console
	"os"      // Package to support the os.Error type
	"time"    // Package to support time.Sleep command
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

// Helper function to print ms singe program execution began
func lt() { fmt.Print("[ ", time.Nanoseconds()/1e6-begin, "]: ") }

// Factory function to create communication channel with server
func inChanFactory() (chan value, os.Error) {
	//Create and initialize the import channel
	imp, err := netchan.NewImporter("tcp", "127.0.0.1:2345")
	lt()
	fmt.Println("New incoming Connection: 127.0.0.1:2345")
	if err != nil {
		return nil, err
	}
	lt()
	fmt.Println("incoming importFactory Channel Made")
	//Make the incoming communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the import channel
	err = imp.Import("exportChannelSend", ch, netchan.Recv)
	if err != nil {
		return nil, err
	}
	return ch, err
}
func outChanFactory() (chan value, os.Error) {
	//Create and initialize the import channel
	imp, err := netchan.NewImporter("tcp", "127.0.0.1:2346")
	lt()
	fmt.Println("New outgoing connection: 127.0.0.1:2346")
	if err != nil {
		return nil, err
	}
	lt()
	fmt.Println("outgoing importFactory Channel Made")
	//Make the outgoing communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the import channel
	err = imp.Import("exportChannelRecv", ch, netchan.Send)
	if err != nil {
		return nil, err
	}
	return ch, err
}

//This thread/function spools on an incoming data channel
//Whenever data is recieved, it is decrypted and printed to the screen
func printIncoming(inChan chan value, done chan bool) {
	//Run as long as the inChan channel has not closed
	for !closed(inChan) {
		//Read data from inChan channel into new object
		inval := <-inChan
		lt()
		fmt.Println("Data recieved from server: ", inval)
		//Decrypt the incoming object and print it
		fmt.Println(decrypt(inval.s))
	}
	done <- true
}

//This thread/function spools on the stdid
//Whenever data is entered, it encrypts it into a new object
//Then sends that object on the outbound channel
func acceptOutgoing(outChan chan value, done chan bool) {
	//Create an object to read the stdIn
	input := bufio.NewReader(os.Stdin)
	//Run as long as the outChan channel is not closed
	for i := 0; !closed(outChan); i++ {
		//Read the stdin into a string
		result, _ := input.ReadString('\n')
		text := strings.TrimSpace(result)
		//Create a new outbound data object after encrypting the text
		outval := value{i, encrypt(text)}
		lt()
		fmt.Println("Sending data to outChan: ", outval)
		//Send the data object over the outChan channel
		outChan <- outval
		lt()
		fmt.Println("Data sent")
	}
	done <- true
}

//This function uses the inbound and outbound net channels to establish
//A shared secret with the other party using the Diffie-Hellman key exchange
func keyExchangeClient(outChan chan value, inChan chan value) (success bool) {

	goodConnection := false
	//receive Prime p from server
	inPrime := <-inChan
	lt()
	fmt.Println(inPrime)
	//receive base g from server
	inBase := <-inChan
	lt()
	fmt.Println(inBase)
	//receive Server's public key from server
	inServerPubKey := <-inChan
	lt()
	fmt.Println(inServerPubKey)

	//Generate secret integer b
	randMax := 100000000
	rand.Seed(time.Nanoseconds())
	secretInt := rand.Int() % randMax
	lt()
	fmt.Println("Secret Integer:", secretInt)

	//Generate Client Public key
	//g^b mod p
	PubKey := (inBase.i ^ secretInt) % (inPrime.i)
	lt()
	fmt.Println("Client Public Key:", PubKey)

	//Send Client PubKey to server
	outChan <- value{PubKey, "Client Public Key"}

	//Calculate Shared Secret
	//ClientPubKey^SecretInt mod prime
	sharedSecret = (inServerPubKey.i ^ secretInt) % inPrime.i
	lt()
	fmt.Println("Calculated Shared Secret:", sharedSecret)

	//Answer Server's Challenge
	//response = challenge^sharedSecret mod prime
	inChallenge := <-inChan
	lt()
	fmt.Println(inChallenge)
	response := (inChallenge.i ^ sharedSecret) % inPrime.i
	lt()
	fmt.Println("Response to Server:", response)
	outChan <- value{response, "Response From Client"}

	//Send encrypted challenge to Server
	//challenge formula: nonce^sharedSecret mod prime
	//Server: challenge^sharedSecret mod prime
	//Client checks result
	challenge := ((rand.Int() % randMax) ^ sharedSecret) % inPrime.i
	lt()
	fmt.Println("Challenge to Server:", challenge)
	expectedResult := (challenge ^ sharedSecret) % inPrime.i
	lt()
	fmt.Println("Expected Result:", expectedResult)
	outChan <- value{challenge, "Challenge From Client"}
	result := <-inChan
	lt()
	fmt.Println(result)
	if result.i == expectedResult {
		goodConnection = true
	}
	return goodConnection
}

//simple encryption routine
func encrypt(input string) (result string) {
	data := []int(input)
	for i := 0; i < len(data); i++ {
		data[i] += sharedSecret % 26
	}
	result = string(data)
	return result
}

//simple decryption routine
func decrypt(input string) (result string) {
	data := []int(input)
	for i := 0; i < len(data); i++ {
		data[i] -= sharedSecret % 26
	}
	result = string(data)
	return result
}

//This is the function that establishes communication with the server
func comRoutine() {
	//Use the factory to create a new network communication channel, in and out
	inChan, inerr := inChanFactory()

	if inerr != nil {
		lt()
		fmt.Println("inChanFactory error: ", inerr)
		return
	} else {
		lt()
		fmt.Println("Returned to main after inChan made")
	}
	outChan, outerr := outChanFactory()
	if outerr != nil {
		lt()
		fmt.Println("outChanFactory error: ", outerr)
		return
	} else {
		lt()
		fmt.Println("Returned to main after outChan made")
	}

	//Try to execute a keyexchange
	if keyExchangeClient(outChan, inChan) {
		lt()
		fmt.Println("Keyexchange Success")
		inDone := make(chan bool)
		outDone := make(chan bool)
		
		//Start the incoming and outgoing goroutines spooling
		go printIncoming(inChan, inDone)
		go acceptOutgoing(outChan, outDone)
		select {//Wait on a spooler quit return
		case <-inDone:
			lt()
			fmt.Println("Incoming Channel quit")
		case <-outDone:
			lt()
			fmt.Println("Outgoing Channel quit")
		}
		lt()
		fmt.Println("Killing other thread")
	} else {
		lt()
		fmt.Println("Keyexchange Failed")
	}
	time.Sleep(1e9)
}
func main() {
	//initialize program time
	begin = time.Nanoseconds() / 1e6
	for {
		// Attempt to form a com channel, and run if successful
		comRoutine()
		// If com channel quits/fails then wait 1 second and run again, forever.
		time.Sleep(2e9)
	}
}
