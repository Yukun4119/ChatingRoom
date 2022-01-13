package main

type COMMAND_TYPE int

const (
	CMD_NICK COMMAND_TYPE = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	cmd    COMMAND_TYPE
	client *client
	args   []string
}
