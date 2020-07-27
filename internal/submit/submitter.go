package submit

import (
	"../exploit"
	"../promStats"
	"bufio"
	//"bytes"
	"context"
	//"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func httpSub(submitCtx context.Context) {

	//flagChannel := submitCtx.Value("flagChannel").(chan flag)
	//gameServer := submitCtx.Value("gameServer").(string)
	//token := submitCtx.Value("token").(string)
	//alreadySubmitted := submitCtx.Value("alreadySubmitted").(*sync.Map)
	//
	//var flags []string
	//for {
	//	select {
	//	//Read the flag
	//	case flag := <-flagChannel:
	//		flags = append(flags, flag)
	//	//Create json from flag
	//	case <-time.After(5 * time.Second):
	//		if flags != nil {
	//			flagJson, err := json.Marshal(flags)
	//			if err != nil {
	//
	//				log.Fatalf("SUBMITTER\nError in json marshal with %s\nTrace: %s\n", gameServer, err)
	//			}
	//			req, err := http.NewRequest("PUT", gameServer, bytes.NewBuffer(flagJson))
	//			//Add headers
	//			req.Header.Set("X-Team-Token", token)
	//			if err != nil {
	//				log.Fatalf("SUBMITTER\tConnection Error HTTP:\t Server %s\n Trace:%s\n", gameServer, err)
	//			}
	//			//Send flag
	//			client := &http.Client{
	//				Timeout: time.Second * 5,
	//			}
	//			resp, err := client.Do(req)
	//			if err != nil {
	//				log.Fatalf("SUBMITTER\tError Send Flag:\t Server %s\nTrace: %s\n", gameServer, err)
	//			}
	//			defer resp.Body.Close()
	//			var flagResult []RuCtfFlag
	//			//Parse response
	//
	//			err = json.NewDecoder(resp.Body).Decode(&flagResult)
	//			if err != nil {
	//				log.Fatalf("SUBMITTER\tError Unmarshalling Flag:\nTrace: %s\n", err)
	//			}
	//			for _, flagStatus := range flagResult {
	//				if flagStatus.Status {
	//					alreadySubmitted.Store(flagStatus.Flag, true)
	//				} else {
	//					log.Printf("SUBMITTER\tInvalid Flag:\t %s \n", flagStatus.Flag)
	//				}
	//			}
	//			flags = nil
	//		}
	//	}
	//}

}

func ncSub(submitCtx context.Context) {

	flagChannel := submitCtx.Value("flagChannel").(chan exploit.Flag)
	gameServer := submitCtx.Value("gameServer").(string)
	acceptedFlag := submitCtx.Value("flagAccepted").(string)
	flagDuplicated := submitCtx.Value("flagDuplicated").(string)
	alreadySubmitted := submitCtx.Value("alreadySubmitted").(*sync.Map)

	//Create the tcp connection
	connection, err := net.DialTimeout("tcp", gameServer, 100*time.Second)
	if err != nil {
		log.Fatalf("SUBMITTER\tConnection Error TCP:\t Server %s\n Trace:%s\n", gameServer, err)
	}
	for {
		//Buffered reader
		reader := bufio.NewReader(connection)
		select {
		//Read the flag
		case flag := <-flagChannel:
			//Send the flag
			flagValue := flag.Flag
			_, _ = fmt.Fprintf(connection, "%s\n", flagValue)
			//Read the response
			response, _ := reader.ReadString('\n')
			// Check if it was already sent
			_, ok := alreadySubmitted.Load(flagValue)
			//If it's accepted, store it

			if strings.Contains(response, acceptedFlag) && !ok {
				fmt.Println("SENT", flagValue)
				promStats.Stats.IncrementSubmitted()

				alreadySubmitted.Store(flagValue, true)
			} else if strings.Contains(response, flagDuplicated) && !ok {
				fmt.Println("DUPLICATED", flag)
				promStats.Stats.IncrementDuplicated()
				alreadySubmitted.Store(flag, true)

			} else {
				fmt.Println("ERROR", flagValue, "response", response)
				promStats.Stats.IncrementFailed()
				alreadySubmitted.Store(flagValue, true)
			}

		//After x seconds without flag, stop
		case <-time.After(10 * time.Second):
			log.Printf("SUBMITTER\tRestarting connection with flag master")
			_ = connection.Close()
			time.Sleep(5 * time.Second)
			connection, err = net.DialTimeout("tcp", gameServer, 10*time.Second)
			if err != nil {
				log.Fatalf("SUBMITTER\tConnection Error TCP:\t Server %s\n Trace:%s\n", gameServer, err)
			}
		}
	}
}

func StartSubmitter(submitterCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	toSubmit := submitterCtx.Value("submit").(chan exploit.Flag)
	workers := submitterCtx.Value("workers").(int)
	var submitted sync.Map

	//Create channel to pass filtered flags
	flagChannel := make(chan exploit.Flag, workers*5)
	submitterCtx = context.WithValue(submitterCtx, "flagChannel", flagChannel)
	submitterCtx = context.WithValue(submitterCtx, "alreadySubmitted", &submitted)

	//Start the submitter
	switch subType := submitterCtx.Value("subType").(string); subType {
	case "nc":
		go ncSub(submitterCtx)
	case "http":
		go httpSub(submitterCtx)
	default:
		log.Fatalf("SUBMITTER:\n Submission type %s doesn't exist!\n", subType)
	}
	//Check if the flags are already submitted
	for flag := range toSubmit {
		//The regex is checked via exploiter
		if _, result := submitted.Load(flag.Flag); result {
			// already submitted
			continue
		} else {
			// new flag
			flagChannel <- flag
		}
	}

}
