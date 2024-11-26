package main

import (
	"github.com/abnormal-security/protoconf/command"
	"github.com/abnormal-security/protoconf/compiler"
)

func main() {
	// compile does not need a GC
	// debug.SetGCPercent(-1)
	command.RunCommand("compiler", compiler.Command)
}
