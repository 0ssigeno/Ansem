package internal

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func ReadTeams(fileTeam string) []string {
	//teamChannel := make(chan string, 20)
	file, err := os.Open(fileTeam)
	if err != nil {
		log.Fatalf("TEAM:\t%s is not a valid file!\n", fileTeam)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	err = file.Close()
	if err != nil {
		log.Fatalf("TEAM:\t%s error on closing file!\n", fileTeam)
	}

	return lines
}

func MakeChan(m sync.Map, threshold int) chan string {
	teamChannel := make(chan string, 20)

	m.Range(func(key, value interface{}) bool {
		if int(*value.(*int32)) > threshold {

			log.Printf("TEAM %s exceeded threshold \n", key.(string))
		} else {
			teamChannel <- key.(string)
		}
		return true
	})
	return teamChannel
}
