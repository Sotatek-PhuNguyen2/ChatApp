package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		switch cmd {
		case "USER":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "CHAT":
			c.commands <- command{
				id:     CMD_CHAT,
				client: c,
				args:   args,
			}
		case "SHOW":
			c.commands <- command{
				id:     CMD_SHOW,
				client: c,
			}
		case "JOIN":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "ROOM":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "MSG":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "QUIT":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		case "QUITROOM":
			c.commands <- command{
				id:     CMD_QUITROOM,
				client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
