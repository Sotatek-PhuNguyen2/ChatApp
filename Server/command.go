package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_CHAT
	CMD_SHOW
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_QUITROOM
)

type command struct {
	id     commandID
	client *client
	args   []string
}
package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_CHAT
	CMD_SHOW
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_QUITROOM
)

type command struct {
	id     commandID
	client *client
	args   []string
}
