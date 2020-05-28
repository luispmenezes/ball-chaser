# Ball Chaser

A GO library to parse Rocket League match replays (**only headers for now**).  

## Getting Started

```
go mod init github.com/my/repo 
go get github.com/luispmenezes/ball-chaser
```

## Example

```go
package main

import (
	"fmt"
	"github.com/luispmenezes/ball-chaser/pkg/ballchaser"
	"log"
	"os"
)

func main() {
	f, err := os.Open("/path/rep-file.replay")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	parser := ballchaser.NewParser(f)

	fmt.Printf("Parsing Replay: %s\n", parser.Header.Label)
	team1Score, err := parser.Header.GetTeam1Score()

	if err != nil {
		log.Fatal(err)
	}

	team2Score, err := parser.Header.GetTeam2Score()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Score %d - %d \n", team1Score, team2Score)

	playerStats, err := parser.Header.GetPlayerStatistics()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Player Statistics:")

	for _, stat := range playerStats {
		fmt.Printf("Name:%20s Team:%d Score:%4d G:%2d A:%2d Sa:%2d Sh:%2d\n", stat.Name, stat.Team, stat.Score,
			stat.Goals, stat.Assists, stat.Saves, stat.Shots)
	}
}
```