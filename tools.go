package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func ValidateIPv4(host string) bool {
	// validate that supplied host ip fits the IPv4 address range

	if host == "localhost" {
		return true
	}

	parts := strings.Split(host, ".")
	if len(parts) < 4 {
		return false
	}

	for _, n := range parts {
		if i, err := strconv.Atoi(n); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}
	return true
}

func GetTime() string {
	// get system time (utc) (hh:mm:ss format)

	t := time.Now()
	formatted := fmt.Sprintf("%02d:%02d:%02d",
		t.Hour(),
		t.Minute(),
		t.Second())

	return formatted
}

func ExitGracefully() {
	// exit program with exit code 0

	fmt.Println(" \n \b")
	os.Exit(0)
}

func InputInCommands(in string, coms []string) bool {
	// check if 'in' string is contained in the 'coms' list

	if len(coms) != 0 {
		for _, com := range coms {
			if com == in {
				return true
			}
		}
	} else {
		return true
	}

	return false
}

// currently unused
func ExecuteCommand(shell string, args []string) string {
	// use custom shell and param args to execute command

	cmd := exec.Command(shell, args...)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return "None"
	}

	return string(stdout)
}

// replaced by inbuilt interface *gocui.View Clear()
func ClearConsole() {
	// use ANSI code to clear console screen

	fmt.Printf("\x1bc") // clear from bottom ^ up
}

// replaced by interface std input
func Input(text string, req []string) string {
	// loop get user input with required input check

	var input string
	fmt.Printf("\n    [%s?%s] %s", Yellow, Reset, text)
	fmt.Scanln(&input)

	for !InputInCommands(input, req) {
		fmt.Printf("\n    [%s?%s] %s", Yellow, Reset, text)
		fmt.Scanln(&input)
	}

	return input
}

// replaced as tcp socket monitor is now internal
func ReplaceMonitorHost() {
	// replace host for external tcp socket monitor

	input, err := os.ReadFile("../monitor/monitor.go")
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "var host = ") {
			new_host := fmt.Sprintf(`var host = "%s"`, ip)
			lines[i] = new_host
		}
	}
	output := strings.Join(lines, "\n")

	err = os.WriteFile("../monitor/monitor.go", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
