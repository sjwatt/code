/************************************************************
 * This is the server side of a simple (not secure)
 * implementation of the Diffie-hellman key exchange.
 * In this instance it is used to facilitate a command
 * line chat program that uses the key exchange to generate
 * a shared key for simple rotation cipher encryption
 * Author: Simon Watt sjwatt@gmail.com
 * Course: INFO3270
 * Instructor: Dr. Abhijit Sen
 * Filename: DHServer.go
 * Language: GO
 */

package main

import (
	"netchan" // Package for network channel-based communication
	"fmt"     // Package to allow printing to the console
	"os"      // Package to support the os.Error type
	"bufio"
	"strings"
	"time"
	"rand"
)

// Struct type for Data communication
type value struct {
	i int
	s string
}

var begin int64
var sharedSecret int
//Macro helper to print program time
func lt() { fmt.Print("[ ", time.Nanoseconds()/1e6-begin, "]: ") }

// Factory function to create listening communication channel that will wait for a connection from the client
func outChanFactory() (chan value, os.Error) {
	//Create and initialize the export channel
	exp, err := netchan.NewExporter("tcp", ":2345")
	lt()
	fmt.Println("exp.Addr().String(): ", exp.Addr().String())
	if err != nil {
		return nil, err
	}
	lt()
	fmt.Println("outgoing exportFactory Channel Made")
	//Make the communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the export channel
	err = exp.Export("exportChannelSend", ch, netchan.Send)
	if err != nil {
		return nil, err
	}
	return ch, err
}

func inChanFactory() (chan value, os.Error) {
	//Create and initialize the export channel
	exp, err := netchan.NewExporter("tcp", ":2346")
	lt()
	fmt.Println("exp.Addr().String(): ", exp.Addr().String())
	if err != nil {
		return nil, err
	}
	lt()
	fmt.Println("incoming exportFactory Channel Made")
	//Make the communication channel for this program
	ch := make(chan value)
	//Link the ch channel to the export channel
	err = exp.Export("exportChannelRecv", ch, netchan.Recv)
	if err != nil {
		return nil, err
	}
	return ch, err
}

//This thread/function spools on an incoming data channel
//Whenever data is recieved, it is decrypted and printed to the screen
func printIncoming(inChan chan value, quit chan bool) {
	//Run as long as the inChan channel has not closed
	for !closed(inChan) {
		//Read data from inChan channel into new object
		inval := <-inChan
		lt()
		fmt.Println("Data recieved from Client: ", inval)
		fmt.Println(decrypt(inval.s))
	}
	quit <- true
}

//This thread/function spools on the stdid
//Whenever data is entered, it encrypts it into a new object
//Then sends that object on the outbound channel
func acceptOutgoing(outChan chan value, quit chan bool) {
	//Create an object to read the stdin
	input := bufio.NewReader(os.Stdin)
	//Run as long as the outChan channel is not closed
	for i := 0; !closed(outChan); i++ {
		//Read the stdin into a string
		result, _ := input.ReadString('\n')
		text := strings.TrimSpace(result)
		//Create a new outbound data object after encrypting thee text
		outval := value{i, encrypt(text)}
		lt()
		fmt.Println("Sending data to outChan: ", outval)
		//Send the data object over the outChan channel
		outChan <- outval
		lt()
		fmt.Println("Data sent")
	}
	quit <- true
}

//This function uses the inbound and outbound net channels to establish
//A shared secret with the other party using the Diffie-Hellman key exchange
func keyExchangeServer(outChan chan value, inChan chan value) (success bool) {

	goodConnection := false
	//Generate random prime p (very dumb algorithm)
	lt()
	fmt.Println("Generating Prime Number p")
	rand.Seed(time.Nanoseconds())
	randMax := rand.Int()
	randMax = randMax % 100000000
	//computedPrime := 0
	running := true
	prime := randMax
	for ; prime > 0; prime-- {
		running = false
		for j := 2; j < randMax/2; j++ {
			if prime%j == 0 {
				running = true
				break
			}
		}
		if running == false {
			break
		}
	}
	lt()
	fmt.Println("Prime p found:", prime)

	//generate base g which is coprime to p
	//every number except 1 is coprime to p since p is prime
	base := randMax - (rand.Int() % randMax)
	lt()
	fmt.Println("Base g:", base)

	//Send p and g to the client
	outChan <- value{prime, "prime"}
	outChan <- value{base, "base"}

	//Choose secret integer a
	secretInt := rand.Int() % randMax
	lt()
	fmt.Println("Secret Integer:", secretInt)

	//Calculate G^a mod p (Public Key)
	PubKey := (base ^ secretInt) % (prime)
	lt()
	fmt.Println("Server Public Key:", PubKey)

	//Send Public Key to Client
	outChan <- value{PubKey, "Server Public Key"}

	//receive Client PubKey
	inClientPubKey := <-inChan
	lt()
	fmt.Println(inClientPubKey)

	//Calculate Shared Secret
	//ClientPubKey^SecretInt mod prime
	sharedSecret = (inClientPubKey.i ^ secretInt) % prime
	lt()
	fmt.Println("Calculated Shared Secret:", sharedSecret)

	//Send encrypted challenge to client
	//challenge formula: nonce^sharedSecret mod prime
	//client: challenge^sharedSecret mod prime
	//server checks result
	challenge := ((rand.Int() % randMax) ^ sharedSecret) % prime
	lt()
	fmt.Println("Challenge to Client:", challenge)
	expectedResult := (challenge ^ sharedSecret) % prime
	lt()
	fmt.Println("Expected Result:", expectedResult)
	outChan <- value{challenge, "Challenge From Server"}
	result := <-inChan
	lt()
	fmt.Println(result)
	if result.i == expectedResult {
		goodConnection = true
	}

	//Answer Client's Challenge
	//response = challenge^sharedSecret mod prime
	inChallenge := <-inChan
	lt()
	fmt.Println(inChallenge)
	response := (inChallenge.i ^ sharedSecret) % prime
	lt()
	fmt.Println("Response to Client:", response)
	outChan <- value{response, "Response From Server"}

	return goodConnection
}

//This is the function that establishes communication with the client
func comRoutine() {

	//Create a new [value] variable
	//Use the factory to create a new network communication channel
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
	
	if keyExchangeServer(outChan, inChan) {
		lt()
		fmt.Println("Keyexchange Success")
		inQuit := make(chan bool)
		outQuit := make(chan bool)
		go printIncoming(inChan, inQuit)
		go acceptOutgoing(outChan, outQuit)
		select {
		case <-inQuit:
			lt()
			fmt.Println("Incoming Channel quit")
		case <-outQuit:
			lt()
			fmt.Println("Outgoing Channel quit")
		}
		lt()
		fmt.Println("Killing other thread")
	} else {
		lt()
		fmt.Println("Key Exchange Failed")
		
	}
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
func main() {
	//initialize program time
	begin = time.Nanoseconds() / 1e6

	// Attempt to form a com channel, and run if successful
	comRoutine()
	// If com channel quits/fails then wait 1 second and run again, forever.

}
