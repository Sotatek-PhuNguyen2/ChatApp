package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
	account  map[string]*client
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command, 1000),
		account:  make(map[string]*client),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_SHOW:
			s.show(cmd.client)
		case CMD_CHAT:
			s.chat(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_QUITROOM:
			s.quitRoom(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) *client {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
}

func (s *server) nick(c *client, args []string) {
	k := 0
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME")
		return
	}
	for username := range s.account {
		if username == args[1] {
			k = 1
		}
	}
	delete(s.account, c.nick)
	if k == 0 {
		c.nick = args[1]
		s.account[c.nick] = c
		c.msg(fmt.Sprintf("All right, I will call you %s", c.nick))
	} else {
		c.msg(fmt.Sprintf("Duplicate account: %s", args[1]))
	}
}

func (s *server) chat(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /chat <username> MSG")
		return
	}

	username := args[1]
	msg := strings.Join(args[2:], " ")
	for account := range s.account {
		if account == username {
			c.send(s.account[username], c.nick+": "+msg)
		}
	}
}

func (s *server) show(c *client) {
	var account []string
	for username := range s.account {
		if username != c.nick {
			account = append(account, username)
		}
	}
	c.msg(fmt.Sprintf("Online list: %s", strings.Join(account, ", ")))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: ROOM_NAME")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	if c.room != nil {
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
	}

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))
	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())
	delete(s.account, c.nick)
	s.quitRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
