package internal

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"
)

func GetTeamAsChan(fileTeam string) chan string {
	teamChannel := make(chan string, 20)
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
	file.Close()
	go func() {
		for _, line := range lines {
			teamChannel <- line
		}
	}()
	return teamChannel
}
