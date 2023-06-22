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
}

func newServer() *server {
	return &server{
		make(map[string]*room),
		make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.cmd {
		case CMD_NICK:
			s.nick(cmd.client, len(cmd.args), cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, len(cmd.args), cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.sendMsg(cmd.client, len(cmd.args), cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_HELP:
			s.help(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new Client : %s", conn.RemoteAddr().String())
	c := &client{
		conn:     conn,
		nick:     "Anonymous",
		commands: s.commands,
	}
	c.msg2client(cheatSheet)
	c.readInput()
}

func (s *server) nick(c *client, argc int, args []string) {
	if argc < 2 {
		c.msg2client("Error: nick name is required. Usage : /nick <Name>")
		return
	}
	c.nick = args[1]
	c.msg2client(fmt.Sprintf("Welcome! %s", c.nick))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg2client(fmt.Sprintf("Rooms available: %s", strings.Join(rooms, ", ")))
	c.msg2client("You can also use /join to create a new room")
}

func (s *server) join(c *client, argc int, args []string) {
	if argc < 2 {
		c.msg2client("Error: room name is required. Usage : /join <Room Name>")
		return
	}
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		// create a new room
		r = &room{
			roomName,
			make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))
	c.msg2client(fmt.Sprintf("Welcome to room %s", roomName))
}

func (s *server) sendMsg(c *client, argc int, args []string) {
	if argc < 2 {
		c.msg2client("Error! Usage: /msg <Message>")
		return
	}
	if c.room == nil {
		c.msg2client("You have not entered a room! PLS choose a room to start chatting!")
		return
	}
	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("[In server] % has left the room", c.conn.RemoteAddr().String())
	if c.room != nil {
		curRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		curRoom.broadcast(c, fmt.Sprintf("[In client] %s has left the room", c.nick))
	}
	c.msg2client("Bye~\n")
	c.conn.Close()
}

func (s *server) help(c *client) {
	c.msg2client(cheatSheet)
}
