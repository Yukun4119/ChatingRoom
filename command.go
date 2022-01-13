package main

type COMMAND_TYPE int

const cheatSheet = "Here is cheat sheet:\n /nick <nick name>: pick a nick name\n /join <room name>: join a room\n /rooms: list all the rooms\n /msg <message>: send message in the room\n /quit: log out\n ?: get help\n"

const (
	CMD_NICK COMMAND_TYPE = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_HELP
)

type command struct {
	cmd    COMMAND_TYPE
	client *client
	args   []string
}
