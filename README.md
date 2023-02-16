# File Change Monitoring

This is a 'file change' monitoring package.

**[ Tips: Continuously updating ... ]**

The package `monitor` basically implements the monitoring of changes to a single file. After creating a file listener in general, it initializes (loads) the initial state of the file; creates a cycle time, continuously refreshes and obtains information and md5 values of the file through the file listener, checks the difference between the initial state of the created file listener task and the current state (MD5/ModTime/Mode), and calls the incoming `on()` to output the information.

```shell
go run main.go -h
```

```shell
go run main.go -file ./1.txt
```

When the program is running, it receives action commands via standard inputs.

e.g:

- enter `stop` : close the file listener and exit the program.
- enter `reset`: will reset the accident status of the file


Example：

> This is the program code in the `main.go` file, which is a sample code that implements the above requirements

``` go
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
```