package main

import (
	protoconf_pb "github.com/abnormal-security/protoconf/pb/protoconf/v1"
	_ "github.com/bufbuild/protovalidate-go"
	_ "github.com/bufbuild/protovalidate-go/legacy"

	"github.com/abnormal-security/protoconf/agent"
	"github.com/abnormal-security/protoconf/command"
	"github.com/abnormal-security/protoconf/compiler"
	"github.com/abnormal-security/protoconf/devserver"
	protoconffmt "github.com/abnormal-security/protoconf/fmt"
	"github.com/abnormal-security/protoconf/inserter"
	"github.com/abnormal-security/protoconf/mod"
	"github.com/abnormal-security/protoconf/mutate"
	"github.com/abnormal-security/protoconf/server"
	"github.com/mitchellh/cli"
)

func init() {
	newCompileRequest := &protoconf_pb.ConfigMutationResponse{}
	_ = newCompileRequest
}

func main() {
	command.RunSubcommands("protoconf",
		map[string]cli.CommandFactory{
			"agent":     agent.Command,
			"compile":   compiler.Command,
			"devserver": devserver.Command,
			"fmt":       protoconffmt.Command,
			"insert":    inserter.Command,
			"mutate":    mutate.Command,
			"serve":     server.Command,
			"mod init":  mod.NewInitCommand,
			"mod sync":  mod.NewSyncCommand,
			"mod tidy":  mod.NewTidyCommand,
		},
	)
}
