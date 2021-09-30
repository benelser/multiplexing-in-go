package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	// Create our tx || rx channel which will be handling a *Person{}
	personChannel := make(chan *Person)

	// pass our person channel in order to communicate back to main go routine which is multplexed below
	go GetDataFromChildRoutineBackToParent(personChannel)

	wireUpSender := make(chan *Person)
	wireUpReceiver := make(chan *Person)
	go ReceiveDataFromParentRoutineAndEchoBack(wireUpSender, wireUpReceiver)
	
	// Multiplexing pattern
	fmt.Println("Going into infinite loop to receive our data from our go routines")

	var timer *time.Timer
	timer = time.NewTimer(5 * time.Second)
	for {
		select {
		// Syntax is receiving the data off of channel person into var me which is of type Person{}
		case me := <-personChannel:
			fmt.Printf("From our person channel we received the following name: %v\n", me.Name)
			return
                // Syntax is receiving the data off of channel person into var echoPerson which is of type Person{}
		case echoPerson := <-wireUpReceiver:
			fmt.Printf("Echo case fired!!!!!. Name provided is: %s\n", echoPerson.Name)
			continue
                // Case to execute our scheduled job //// whatever it may be here we are daisy-chaining channels follow along in output
		case <-timer.C:
			fmt.Println("Our 5 second timer expired. Sending signal to our Sender channel (wireUpSender )")
			wireUpSender <- &Person{Name: "Our echo dude"}
			timer = time.NewTimer(5 * time.Second)
			continue
		default:
		}
	}

}

type Person struct {
	Name string
}

func GetDataFromChildRoutineBackToParent(communicationsChannel chan<- *Person) {
	// BLOCKS and waits for a key to be pressed. Once key is pressed we are into our go routine
	os.Stdin.Read(make([]byte, 1))
	me := &Person{
		Name: "Benjamin Elser",
	}
	// send our pointer to Person back up to parent go routine
	communicationsChannel <- me
}

func ReceiveDataFromParentRoutineAndEchoBack(rxCommunicationsChannel chan *Person, txCommunicationsChannel chan<- *Person) {
	// rx on one channel and tx to another
  for {
		select {
		case receivedData := <-rxCommunicationsChannel:
			txCommunicationsChannel <- receivedData

		}
	}
}
