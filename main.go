package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	prompt "github.com/c-bata/go-prompt"
	"github.com/gorilla/websocket"
)

func connect(addr string, userID string, deviceID string) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s?id=%s&deviceid=%s", addr, userID, deviceID), nil)
	return c, err
}

func executor(t string) {
	if t == "bash" {
		cmd := exec.Command("bash")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	return
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "bash"},
	}
}

func main() {
	{
		p := prompt.New(
			executor,
			completer,
		)
		p.Run()
	}

	addr := flag.String("addr", "localhost:8080", "server ws address")
	userID := flag.String("userId", "1", "userid")
	deviceID := flag.String("deviceId", "1", "device id")

	flag.Parse()
	log.SetFlags(0)

	conn, err := connect(*addr, *userID, *deviceID)
	if err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	go handleIncomingEvent(conn)

	for {
		select {
		case <-interrupt:
			log.Println("Exiting ...")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}

func handleIncomingEvent(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("error reading: ", err)
			continue
		}
		log.Println(message)
	}
}
