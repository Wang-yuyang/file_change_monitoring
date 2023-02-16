package file_change_monitoring

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	filepath := "./1.txt"
	f, _ := os.Create(filepath)
	mon := NewMonitor(filepath)

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

	for {
		if msg, ok, _ := mon.VerifyFileChange(); ok {
			log.Printf("%s [monitor] > %s \n", filepath, msg)
			init, now, _ := mon.OutFileInfo()
			fmt.Printf("%s <=> %s \n", init.FileHash, now.FileHash)
		}
		time.Sleep(5 * time.Second)
	}
}
