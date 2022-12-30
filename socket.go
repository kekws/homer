package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	listener net.Listener
	addr     string
}

type Socket struct {
	protocol string
	host     string
	port     string
}

type Connection struct {
	id       int
	conn     net.Conn
	username string
	ip       string
	port     int
	state    bool
}

type Log struct {
	neutral  int
	incoming int
	outgoing int
	logs     []string
}

var connections = []Connection{}

var socketLog = Log{0, 0, 0, []string{}}

func StartServer(host string, port string) {
	// set up server net socket using custom Socket struct

	sock := Socket{
		protocol: "tcp",
		host:     host,
		port:     port,
	}

	data := fmt.Sprintf("%s:%s", sock.host, sock.port)
	dstream, err := net.Listen("tcp", data)
	if err != nil {
		fmt.Println("\n\n[!] Can't host server! (maybe the host IP is already in use?)")
		ExitGracefully()
	}

	TCPServer := Server{
		listener: dstream,
		addr:     data,
		// monitor:  ConnectMonitorServer(sock.host),
	}
	AcceptConnections(TCPServer.listener)
}

func AcceptConnections(server net.Listener) {
	// continually accept connections from server

	for {
		con, _ := server.Accept()

		socketLog.neutral += 1
		RecieveRaw(con) // clear packet headers
		socketLog.incoming += 1
		go HandleConnections(con) // connection go-routine
	}
}

func HandleConnections(con net.Conn) {
	// handle single connection -> thread

	defer con.Close()

	getConnId := func() int {
		return len(connections) + 1
	}

	if _, ok := con.RemoteAddr().(*net.TCPAddr); ok {
		newConnection := Connection{
			id:    getConnId(),
			conn:  con,
			port:  8080,
			state: true,
		}

		newConnection.username = RecieveRaw(con)
		newConnection.ip = newConnection.Recieve()

		connections = append(connections, newConnection)
		socketLog.neutral += 1
		AddLog("-- Accepted", newConnection.username)

		for {
			data := newConnection.Recieve()
			if !newConnection.state {
				RemoveConnection(newConnection)
				newConnection.state = false
				newConnection.conn.Close()
				AddLog("-- Removed", newConnection.username)
				socketLog.neutral += 1
				break
			} else {
				fmt.Println(data)
			}
		}
	}
}

func (connection *Connection) Send(data string) {
	// simply send param data (+newline) to received client (net.Conn)

	conn := connection.conn
	fmt.Fprintf(conn, data+"\n")
	socketLog.outgoing += 1
}

func (connection *Connection) Recieve() string {
	// wait for incoming string packet and return sanitized data

	conn := connection.conn
	data, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		connection.state = false
		return "None"
	}
	socketLog.incoming += 1
	AddLog("<-", connection.username)

	return strings.TrimSuffix(string(data), "\n")
}

func RecieveRaw(con net.Conn) string {
	// wait for incoming string packet and return data from unknown client

	data, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return "None"
	}
	socketLog.incoming += 1
	AddLog("<-", "????")

	return strings.TrimSuffix(string(data), "\n")
}

func AddLog(text string, user string) {
	// add text (in, out, nil, accepted, removed, etc.)
	// and connection username to socket log list

	logText := fmt.Sprintf("%s %s%s%s", text, session.theme, user, Reset)
	socketLog.logs = append(socketLog.logs, logText)
}

// replaced by id method
func RemoveConnection(toRemove Connection) {
	// iterate through connections and remove the item that matches the supplied connection

	for i, connection := range connections {
		if toRemove.conn == connection.conn && toRemove.state {
			connections = append(connections[:i], connections[i+1:]...)
			session.machineCount = len(connections)
		}
	}
}
