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
		args := strings.Split(msg, " ") // args[0] is the command
		cmd := strings.TrimSpace(args[0])
		switch cmd {
		case "/nick":
			c.commands <- command{
				cmd:    CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				cmd:    CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				cmd:    CMD_ROOMS,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				cmd:    CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				cmd:    CMD_QUIT,
				client: c,
				args:   args,
			}
		case "?":
			c.commands <- command{
				cmd:    CMD_HELP,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg2client(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
