package main

import (
	"fmt"
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			fmt.Println(msg)
			m.msg(msg)
		}
	}
}

func (c *client) send(sender *client, msg string) {
	sender.msg(msg)
}
