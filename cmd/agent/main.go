package main

import (
	"github.com/abnormal-security/protoconf/agent"
	"github.com/abnormal-security/protoconf/command"
)

func main() {
	command.RunCommand("agent", agent.Command)
}
