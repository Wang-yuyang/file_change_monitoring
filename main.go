package main

import (
	"bufio"
	"context"
	"file_change_monitoring/monitor"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var (
		filepath string
	)
	flag.StringVar(&filepath, "file", "", "Select a file to monitor for changes.")

	flag.Parse()
	if filepath != "" {
		_, err := os.Open(filepath)
		if err != nil {
			fmt.Println("[FILE OPEN ERROR] ", err)
			os.Exit(0)
			return
		}
	}

	mon := monitor.NewMonitor(filepath, "", on)

	func(ok bool, err error) {
		if err != nil || !ok {
			fmt.Println(err)
			return
		}
	}(mon.FileInitialState())

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			input := bufio.NewReader(os.Stdin)
			line, _, _ := input.ReadLine()
			switch {
			case string(line) == "stop":
				fmt.Println("logout monitors. ")
				cancel()
				return
			case string(line) == "reset":
				func(ok bool, err error) {
					if err != nil || !ok {
						fmt.Println(err)
					}
					fmt.Println("reset monitors. ")
				}(mon.FileInitialState())
			default:
				continue
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, _, err := mon.VerifyFileChange()
				if err != nil {
					return
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	<-ctx.Done()
}

func on(msg string, level int, err error) {
	if level > 0 {
		fmt.Println(msg)
		return
	} else if level < 0 {
		fmt.Println(msg, level, err)
	} else {
		return
	}
}
