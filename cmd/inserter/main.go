package main

import (
	"github.com/abnormal-security/protoconf/command"
	"github.com/abnormal-security/protoconf/inserter"
)

func main() {
	command.RunCommand("inserter", inserter.Command)
}
