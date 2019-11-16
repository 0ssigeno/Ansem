package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func StartSubmitter(gameServer string, toSubmit <-chan string) {

	//Define submission method
	submitFunction := submitNC

	//Create a map to verify flags
	submitted := make(map[string]bool)
	//Create channel to pass filtered flags
	flagChannel := make(chan string, 10)
	//Create the channel to communicate with the map handler
	mapWrite := make(chan string)
	mapRead := make(chan string)
	mapGet := make(chan bool)

	//Start the handler of the map
	go func() {
		for {
			select {
			case write := <-mapWrite:
				submitted[write] = true
			case read := <-mapRead:
				_, found := submitted[read]
				mapGet <- found
			}
		}
	}()

	//Start the submitter
	go submitFunction(gameServer, flagChannel, mapWrite)

	//Check if the flags are already submitted
	for flag := range toSubmit {
		mapRead <- flag
		present := <-mapGet
		if present {
			continue
		} else {
			flagChannel <- flag
		}
	}

}

func submitNC(gameServer string, flagChannel <-chan string, handler chan<- string) {

	//Create the tcp connection
	connection, err := net.DialTimeout("tcp", gameServer, 10*time.Second)
	if err != nil {
		log.Fatalf("Error in connection with %s\n", gameServer)
	}
	for {
		//Buffered reader
		reader := bufio.NewReader(connection)
		select {
		//Read the flag
		case flag := <-flagChannel:
			//Send the flag
			fmt.Fprintf(connection, "%s\n", flag)
			//Read the response
			response, _ := reader.ReadString('\n')
			//If it's accepted, store it
			if strings.Contains(response, "ACCEPTED") {
				handler <- flag

			}
			//After x seconds without flag, stop
		case <-time.After(10 * time.Second):
			connection.Close()
			fmt.Print("Chiudo\n")
			time.Sleep(10 * time.Second)
			connection, err = net.DialTimeout("tcp", gameServer, 10*time.Second)
			if err != nil {
				log.Fatalf("Error in connection with %s\n", gameServer)
			}
		}

	}

}
