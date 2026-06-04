package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/manifoldco/promptui"
)

type ServerMessage struct {
	Type    string   `json:"type"`
	Title   string   `json:"title,omitempty"`
	Items   []string `json:"items,omitempty"`
	Board   string   `json:"board,omitempty"`
	Message string   `json:"message,omitempty"`
}

type ClientMessage struct {
	Index int    `json:"index,omitempty"`
	Input string `json:"input,omitempty"`
}

func clearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(conn)

	for {
		var msg ServerMessage
		if err := decoder.Decode(&msg); err != nil {
			fmt.Println("Disconnected")
			return
		}

		switch msg.Type {
		case "menu":
			prompt := promptui.Select{
				Label: msg.Title,
				Items: msg.Items,
			}
			index, _, err := prompt.Run()
			if err != nil {
				return
			}
			if err := encoder.Encode(ClientMessage{Index: index}); err != nil {
				return
			}

		case "prompt":
			prompt := promptui.Prompt{
				Label: msg.Title,
			}
			input, err := prompt.Run()
			if err != nil {
				return
			}
			if err := encoder.Encode(ClientMessage{Input: input}); err != nil {
				return
			}

		case "board":
			clearScreen()
			fmt.Println(msg.Board)

		case "info":
			fmt.Println()
			fmt.Println(msg.Message)
			fmt.Println()

		case "bye":
			fmt.Println()
			fmt.Println(msg.Message)
			return

		default:
			fmt.Printf("Неизвестный тип сообщения: %s\n", msg.Type)
		}
	}
}
