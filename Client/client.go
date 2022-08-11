package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var ret = make(chan string, 10)
var isRoom = 0

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		response, _ := reader.ReadString('\n')
		response = response[:len(response)-1]
		fmt.Println(response)
		response = response[:9]
		if response == "> welcome" {
			isRoom = 1
		}
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	connection, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	go onMessage(connection)

	for {
		var msg string
		if isRoom == 1 {
			inputReader := bufio.NewReader(os.Stdin)
			input, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Println("Something went wrong")
				continue
			}
			input = input[:len(input)-2]
			ret := strings.Split(input, " ")
			if input == "/quit" {
				msg = "QUITROOM "
				isRoom = 0
			} else if ret[0] == "/private" {
				msg = "CHAT " + strings.Join(ret[1:], " ")
			} else {
				msg = "MSG " + input
			}
			msg = msg + "\n"
		} else {
			fmt.Println("Please chose the opption: ")
			fmt.Println("1. Set username ")
			fmt.Println("2. Show all rooms ")
			fmt.Println("3. Join the room ")
			fmt.Println("4. Show people online")
			fmt.Println("5. Chat with person")
			fmt.Println("6. Quit")

			inputReader := bufio.NewReader(os.Stdin)
			input, err := inputReader.ReadString('\n')
			input = input[:len(input)-2]
			if err != nil {
				fmt.Println("Something went wrong")
				continue
			}

			if input == "1" {
				fmt.Print("Please enter username: ")
				usernameReader := bufio.NewReader(os.Stdin)
				username, err := usernameReader.ReadString('\n')
				username = username[:len(username)-2]
				if err != nil {
					fmt.Println("Something went wrong")
					continue
				}
				msg = "USER " + username
			} else if input == "2" {
				msg = "ROOM "
			} else if input == "3" {
				fmt.Print("Please choose the room: ")
				roomReader := bufio.NewReader(os.Stdin)
				rooms, err := roomReader.ReadString('\n')
				rooms = rooms[:len(rooms)-2]
				if err != nil {
					fmt.Println("Something went wrong")
					continue
				}
				msg = "JOIN " + rooms
			} else if input == "4" {
				msg = "SHOW"
			} else if input == "5" {
				fmt.Println("Syntax:<username> <msg>")
				reqReader := bufio.NewReader(os.Stdin)
				req, err := reqReader.ReadString('\n')
				req = req[:len(req)-2]
				if err != nil {
					fmt.Println("Something went wrong")
					continue
				}
				msg = "CHAT " + req
			} else if input == "6" {
				msg = "QUIT "
			} else {
				fmt.Println("That option is not available, please choose again")
				continue
			}

			msg = msg + "\n"
		}
		connection.Write([]byte(msg))
		time.Sleep(1 * time.Second)
	}
}
