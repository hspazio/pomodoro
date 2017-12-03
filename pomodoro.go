package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/0xAX/notificator"
)

// NOTE:
// pmd > starts timer. Also can take parameters to override default timeout
// ctrl-c to stop  and capture signal for clean shutdown
type results map[string]int

const (
	resultsFile = "pomodoro.dat"
	timeout     = 2
)

var notify *notificator.Notificator

func main() {
	var graph = flag.Bool("graph", false, "show results graph")
	flag.Parse()

	res, err := load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *graph {
		for date, count := range res {
			fmt.Printf("%s: %s %d\n", date, strings.Repeat("#", count), count)
		}
		os.Exit(0)
	}

	tick := time.Tick(1 * time.Second)
	alarm := time.After(timeout * time.Second)

L:
	for {
		select {
		case <-tick:
			fmt.Printf("#")
		case <-alarm:
			res[today()]++
			message := fmt.Sprintf("%d cycles today", res[today()])
			fmt.Printf("\nPomodoro! %s", message)

			notify = notificator.New(notificator.Options{
				AppName: "Pomodoro",
			})

			err := notify.Push(message, "", "", notificator.UR_NORMAL)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err = save(res); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			break L
		}
	}
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func load() (results, error) {
	results := make(map[string]int)
	b, err := ioutil.ReadFile(resultsFile)
	if err == nil {
		if err := json.Unmarshal(b, &results); err != nil {
			return nil, fmt.Errorf("file %s does not contain a valid json", resultsFile)
		}
	}

	return results, nil
}

func save(r results) error {
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("could not encode json %v", err)
	}

	if err := ioutil.WriteFile(resultsFile, b, 0644); err != nil {
		return fmt.Errorf("could not write to file %s", resultsFile)
	}

	return nil
}
