package submitters

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RuCtfFlag struct {
	Msg    string `json:"msg"`
	Flag   string `json:"flag"`
	Status bool   `json:"status"`
}

func RuCTFSubmitHTTP(submitCtx context.Context) {

	flagChannel := submitCtx.Value("flagChannel").(chan string)
	gameServer := submitCtx.Value("gameServer").(string)
	token := submitCtx.Value("token").(string)
	alreadySubmitted := submitCtx.Value("alreadySubmitted").(*sync.Map)

	var flags []string
	for {
		select {
		//Read the flag
		case flag := <-flagChannel:
			flags = append(flags, flag)
		//Create json from flag
		case <-time.After(5 * time.Second):
			if flags != nil {
				flagJson, err := json.Marshal(flags)
				if err != nil {

					log.Fatalf("SUBMITTER\nError in json marshal with %s\nTrace: %s\n", gameServer, err)
				}
				req, err := http.NewRequest("PUT", gameServer, bytes.NewBuffer(flagJson))
				//Add headers
				req.Header.Set("X-Team-Token", token)
				if err != nil {
					log.Fatalf("SUBMITTER\tConnection Error HTTP:\t Server %s\n Trace:%s\n", gameServer, err)
				}
				//Send flag
				client := &http.Client{
					Timeout: time.Second * 5,
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalf("SUBMITTER\tError Send Flag:\t Server %s\nTrace: %s\n", gameServer, err)
				}
				defer resp.Body.Close()
				var flagResult []RuCtfFlag
				//Parse response

				err = json.NewDecoder(resp.Body).Decode(&flagResult)
				if err != nil {
					log.Fatalf("SUBMITTER\tError Unmarshalling Flag:\nTrace: %s\n", err)
				}
				for _, flagStatus := range flagResult {
					if flagStatus.Status {
						alreadySubmitted.Store(flagStatus.Flag, true)
					} else {
						log.Printf("SUBMITTER\tInvalid Flag:\t %s \n", flagStatus.Flag)
					}
				}
				flags = nil
			}
		}
	}

}

/*
Old type of submission
*/
func RuCTFSubmitNC(submitCtx context.Context) {

	flagChannel := submitCtx.Value("flagChannel").(chan string)
	gameServer := submitCtx.Value("gameServer").(string)
	acceptedFlag := submitCtx.Value("flagAccepted").(string)
	flagDuplicated := submitCtx.Value("flagDuplicated").(string)
	// token := submitCtx.Value("token").(string)
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
			fmt.Fprintf(connection, "%s\n", flag)
			//Read the response
			response, _ := reader.ReadString('\n')
			// Check if it was already sent
			_, result := alreadySubmitted.Load(flag)
			//If it's accepted, store it

			if strings.Contains(response, acceptedFlag) && !result {
				fmt.Println("SENDED", flag)
				Stats.IncrementSubmitted()
				alreadySubmitted.Store(flag, true)
			} else if strings.Contains(response, flagDuplicated) && !result {
				fmt.Println("DUPLICATED", flag)
				Stats.IncrementDuplicated()
				alreadySubmitted.Store(flag, true)

			} else {
				fmt.Println("ERROR", flag, "response", response)
				Stats.IncrementFailed()
				alreadySubmitted.Store(flag, true)
			}

		//After x seconds without flag, stop
		case <-time.After(10 * time.Second):
			connection.Close()
			time.Sleep(2 * time.Second)
			connection, err = net.DialTimeout("tcp", gameServer, 10*time.Second)
			if err != nil {
				log.Fatalf("SUBMITTER\tConnection Error TCP:\t Server %s\n Trace:%s\n", gameServer, err)
			}
		}
	}
}
