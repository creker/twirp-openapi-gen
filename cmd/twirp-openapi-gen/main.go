package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/creker/twirp-openapi-gen/internal/generator"
)

type arrayFlags []string

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "Usage of twirp-openapi-gen [flags] <input files>:\n")
		flags.PrintDefaults()
	}

	protoPaths := arrayFlags{}
	servers := arrayFlags{}
	flags.Var(&protoPaths, "proto-path", "Specify the directory in which to search for imports. May be specified multiple times; directories will be searched in order.  If not given, the current working directory is used.")
	flags.Var(&servers, "servers", "Server object URL. May be specified multiple times.")
	title := flags.String("title", "open-api-v3-docs", "Document title")
	docVersion := flags.String("doc-version", "0.1", "API Document version")
	format := flags.String("format", "json", "Document format; json or yaml")
	out := flags.String("out", "./openapi-doc.json", "Output document file")
	pathPrefix := flags.String("path-prefix", "/twirp", "Twirp server path prefix")
	verbose := flags.Bool("verbose", false, "Log debug output")
	printVersion := flags.Bool("version", false, "Print version")

	if err := flags.Parse(args[1:]); err != nil {
		flags.Usage()
		return err
	}
	if *printVersion {
		fmt.Println(version())
		return nil
	}

	if flags.NArg() == 0 {
		flags.Usage()
		return errors.New("No input files specified")
	}

	opts := []generator.Option{
		generator.ProtoPaths(protoPaths),
		generator.Servers(servers),
		generator.Title(*title),
		generator.DocVersion(*docVersion),
		generator.PathPrefix(*pathPrefix),
		generator.Format(*format),
		generator.Verbose(*verbose),
	}
	gen, err := generator.NewGenerator(flags.Args(), opts...)
	if err != nil {
		return err
	}
	if err := gen.Generate(*out); err != nil {
		return err
	}
	return nil
}

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func version() string {
	bi, ok := debug.ReadBuildInfo()
	if ok && bi.Main.Version != "" {
		return bi.Main.Version
	}

	return "DEV"
}
