package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const suspendForSeconds = 30

// waitsForSignal is waiting for given signal. Once signal is received, function sends back information to
// the done channel
func waitsForSignal(signalsChannel chan os.Signal, doneChannel chan bool) {
	for {
		sig := <-signalsChannel
		fmt.Printf("\nwaitsForSignal: I have received signal: %v\n", sig)
		fmt.Println("waitsForSignal: Notifying executor!")
		doneChannel <- true
	}
}

// executor waits for message in done channel, once received it is going to invoke other stuff
func executor(doneChannel chan bool) {
	i := 1
	for {
		fmt.Println("\nexecutor: Awaiting signal...")
		<-doneChannel
		fmt.Printf("executor: Signal received for the %v time\n", i)
		i++
	}
}

func main() {

	signalsChannel := make(chan os.Signal)
	doneChannel := make(chan bool)

	// We are going to support signals below
	mySignals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGILL,
		syscall.SIGTERM,
	}

	// Let us wait for previously defined signals, once signal is received we will populate signalsChannel
	signal.Notify(signalsChannel, mySignals...)

	go waitsForSignal(signalsChannel, doneChannel)
	go executor(doneChannel)

	// Main thread will last for suspendForSeconds seconds
	time.Sleep(suspendForSeconds * time.Second)
	fmt.Println("Timeout exceeded")
}
