package monitor

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	filepath := "./1.txt"
	f, _ := os.Create(filepath)
	mon := NewMonitor(filepath, "", on)

	func(ok bool, err error) {
		if err != nil || !ok {
			fmt.Println(err)
			return
		}
	}(mon.FileInitialState())

	func(err error) {
		if err != nil {
			fmt.Println(err)
			return
		}
	}(f.Close())

	func() {
		for {
			_, _, err := mon.VerifyFileChange()
			if err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()
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
