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

type results map[string]int

const (
	resultsFile = "pomodoro.dat"
	timeout     = 25
)

var notify *notificator.Notificator

func main() {
	var graph = flag.Bool("graph", false, "show results graph")
	flag.Parse()

	var res results
	if err := load(&res); err != nil {
		exit(err)
	}

	if *graph {
		for date, count := range res {
			fmt.Printf("%s: %s %d\n", date, strings.Repeat("#", count), count)
		}
		exit(nil)
	}

	tick := time.Tick(1 * time.Minute)
	alarm := time.After(timeout * time.Minute)
	notify := notificator.New(notificator.Options{AppName: "Pomodoro"})

L:
	for {
		select {
		case <-tick:
			fmt.Printf("#")
		case <-alarm:
			res[today()]++
			message := fmt.Sprintf("%d cycles today", res[today()])
			var suggestion string
			if res[today()]%4 == 0 {
				suggestion = "take a longer break"
			} else {
				suggestion = "rest for 5 minutes"
			}

			fmt.Printf("\nPomodoro! %s, %s", message, suggestion)
			err := notify.Push(message, suggestion, "", notificator.UR_NORMAL)
			if err != nil {
				exit(err)
			}

			if err = save(res); err != nil {
				exit(err)
			}
			break L
		}
	}
}

func load(r *results) error {
	b, err := ioutil.ReadFile(resultsFile)
	if err == nil {
		if err := json.Unmarshal(b, r); err != nil {
			return fmt.Errorf("file %s does not contain a valid json", resultsFile)
		}
	}

	return nil
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

func today() string {
	return time.Now().Format("2006-01-02")
}

func exit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
