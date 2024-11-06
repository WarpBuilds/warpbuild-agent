//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/warpbuilds/warpbuild-agent/cmd/agentd_restarter/cmd"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

const SERVICE_NAME = "warpbuild-agentd-restarter"

type myService struct{}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	status <- svc.Status{State: svc.StartPending}
	log.Println("Service is starting...")

	log.Println("Executing command as goroutine...")
	go func() {
		if err := cmd.ExecuteWithErr(); err != nil {
			log.Printf("Error in cmd.Execute: %v", err)
			status <- svc.Status{State: svc.StopPending}
		}
		log.Println("cmd.Execute() finished")
	}()

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	log.Println("Service is now running")

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				log.Printf("Interrogate received: %v", c.CurrentStatus)
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Print("Shutting service...!")
				break loop
			case svc.Pause:
				log.Print("Pausing service...!")
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				log.Print("Continuing service...!")
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
	runService(SERVICE_NAME, false)
	log.Println("warpbuild-agent service stopped.")
}
