package compiler

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	compilerlib "github.com/abnormal-security/protoconf/compiler/lib"
	"github.com/abnormal-security/protoconf/consts"
	protoconf_pb "github.com/abnormal-security/protoconf/pb/protoconf/v1"
	"github.com/mitchellh/cli"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type cliCommand struct{}

type cliConfig struct {
	repl                    bool
	verboseLogging          bool
	processTemplates        bool
	cpuprofile              string
	memprofile              string
	compilerAddress         string
	src                     string
	includes                string
	prefix                  string
	compiledConfigExtension string
}

func newFlagSet() (*flag.FlagSet, *cliConfig) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: [OPTION]... protoconf_root [config]...")
		flags.PrintDefaults()
	}

	compilerAddress := os.Getenv(`PROTOCONF_COMPILER_ADDR`)
	config := &cliConfig{}
	flags.BoolVar(&config.repl, "repl", false, "Interactive REPL mode")
	flags.BoolVar(&config.verboseLogging, "V", false, "Verbose logging")
	flags.BoolVar(&config.processTemplates, "process-templates", false, "Process template files")
	flags.StringVar(&config.cpuprofile, "cpuprofile", "", "Write cpu profiling info to this file")
	flags.StringVar(&config.memprofile, "memprofile", "", "Write memory profiling info to this file")
	flags.StringVar(&config.compilerAddress, "compiler-address", compilerAddress, "if set, the command will issue a gRPC request to the compiler service at the given address instead of running the compiler locally. The compiler service must be running.")
	flags.StringVar(&config.src, "src", "src/", "The root of protoconf src.")
	flags.StringVar(&config.includes, "includes", "", "files to include in the compilation")
	flags.StringVar(&config.prefix, "prefix", "", "added prefix for proto import path")
	flags.StringVar(&config.compiledConfigExtension, "extension", ".materialized_JSON", "extension for compiled config files")
	return flags, config
}

func (c *cliCommand) Run(args []string) int {
	flags, config := newFlagSet()
	flags.Parse(args)

	if flags.NArg() < 1 {
		flags.Usage()
		return 1
	}
	if config.cpuprofile != "" {
		f, err := os.Create(config.cpuprofile)
		if err != nil {
			slog.Error("Could not create CPU profile:", "error", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			slog.Error("Could not start CPU profile:", "error", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}
	consts.SrcPath = config.src
	consts.Includes = config.includes
	consts.Prefix = config.prefix
	consts.CompiledConfigExtension = config.compiledConfigExtension

	println("abnormal added args: src = ", consts.SrcPath)
	println("abnormal added args: includes = ", consts.Includes)
	println("abnormal added args: prefix = ", consts.Prefix)
	println("abnormal added args: compiled-config-extension = ", consts.CompiledConfigExtension)

	protoconfRoot := strings.TrimSpace(flags.Args()[0])
	var configs []string

	if flags.NArg() == 1 {
		var err error
		configs, err = GetAllConfigs(protoconfRoot)
		if err != nil {
			slog.Error("Error getting all configs", "config", protoconfRoot, "error", err)
			return 1
		}
	} else {
		configs = flags.Args()[1:]
	}

	if config.compilerAddress != "" {
		return runRemote(config, configs)
	}
	return runLocally(protoconfRoot, config, configs)
}

func runRemote(config *cliConfig, configs []string) int {
	conn, err := grpc.Dial(config.compilerAddress, grpc.WithInsecure())
	if err != nil {
		slog.Error("error connecting to server", "error", err)
	}
	client := protoconf_pb.NewProtoconfCompileClient(conn)
	stream, err := client.CompileFiles(context.Background(), &protoconf_pb.CompileRequest{Files: configs})
	if err != nil {
		slog.Error("error compiling files", "error", err)
		return 1
	}
	ret := 0
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return ret
		}
		if err != nil {
			slog.Error("error receiving response", "error", err)
			return 1
		}
		if resp != nil {
			if len(resp.Errors) > 0 {
				slog.Error("Error compiling config", "path", resp.Path, "errors", resp.Errors)
				ret = 1
				continue
			}
			if config.verboseLogging {
				slog.Error("Compiled", "path", resp.Path, "result", resp.Result)
			}
		}
	}

}

func runLocally(protoconfRoot string, config *cliConfig, configs []string) int {
	compiler := compilerlib.NewCompiler(protoconfRoot, config.verboseLogging)
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	if config.repl {
		REPL(compiler)
		return 0
	}

	if config.processTemplates {
		if err := findTemplateFilesAndProccess(); err != nil {
			log.Fatal(err)
		}
	}
	g, _ := errgroup.WithContext(context.Background())

	for _, config := range configs {
		filename := strings.TrimSpace(config)
		g.Go(func() error {
			err := compiler.CompileFile(filename)
			if err != nil {
				ui.Error(fmt.Sprintf("Error compiling config %s:\n    %s", filename, err))
			}
			return err
		})
	}
	err := g.Wait()
	if err != nil {
		// log.Println(err)
		return 1
	}

	if config.memprofile != "" {
		f, err := os.Create(config.memprofile)
		if err != nil {
			log.Fatal("Could not create memory profile:", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("Could not start memory profile:", err)
		}
	}
	return 0
}

func (c *cliCommand) Help() string {
	var b bytes.Buffer
	b.WriteString(c.Synopsis())
	b.WriteString("\n")
	flags, _ := newFlagSet()
	flags.SetOutput(&b)
	flags.Usage()
	return b.String()
}

func (c *cliCommand) Synopsis() string {
	return "Compile configs"
}

// Command is a cli.CommandFactory
func Command() (cli.Command, error) {
	return &cliCommand{}, nil
}

func GetAllConfigs(protoconfRoot string) ([]string, error) {
	srcDir, err := filepath.Abs(filepath.Join(protoconfRoot, consts.SrcPath))
	if err != nil {
		return nil, err
	}

	var configs []string
	err = filepath.Walk(srcDir, func(path string, f os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext == consts.ConfigExtension || ext == consts.MultiConfigExtension {
			configs = append(configs, strings.TrimPrefix(path, srcDir))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return configs, nil
}

func REPL(c *compilerlib.Compiler) {
	fmt.Printf("Protoconf %s\n", consts.Version)

	loader := c.GetLoader()
	thread := &starlark.Thread{
		Load: loader.Load,
	}

	repl.REPL(thread, loader.Modules)
}
