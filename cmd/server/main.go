package main

import (
	"github.com/abnormal-security/protoconf/command"
	"github.com/abnormal-security/protoconf/server"
)

func main() {
	command.RunCommand("server", server.Command)
}
