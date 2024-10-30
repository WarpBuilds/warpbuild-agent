//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/warpbuilds/warpbuild-agent/cmd/agentd/cmd"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

const SERVICE_NAME = "warpbuild-agentd"

type myService struct{}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	status <- svc.Status{State: svc.StartPending}
	log.Println("Service is starting...")

	log.Println("Executing command as goroutine...")
	go cmd.Execute()

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	log.Println("Service is now running")

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Print("Shutting service...!")
				break loop
			case svc.Pause:
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				log.Printf("Unexpected service control request #%d", c)
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	log.Println("Service is stopping")
	return false, 1
}

func runService(name string, isDebug bool) {
	if isDebug {
		err := debug.Run(name, &myService{})
		if err != nil {
			log.Fatalln("Error running service in debug mode.")
		}
	} else {
		err := svc.Run(name, &myService{})
		if err != nil {
			log.Fatalln("Error running service in Service Control mode.")
		}
	}
}

func main() {

	f, err := os.OpenFile("C:/warpbuild-agentd-debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening file: %v", err))
	}
	defer f.Close()

	// Set log output to the file
	log.SetOutput(f)

	log.Println("Starting warpbuild-agent service...")
	runService(SERVICE_NAME, true)
	log.Println("warpbuild-agent service stopped.")
}
