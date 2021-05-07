package main

/*
import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type commands struct {
	name, usage string
	fs          *flag.FlagSet
	fn          func(*flag.FlagSet, []string) error
}
type commandSet map[string]command

func formatFlagHelp(f flag.Flag) string {
	return fmt.Sprintf("%s\n\t%s", f.Name, f.Usage)
}
func formatCommandHelp(cmd command) string {
	return fmt.Sprintf("%s\n\t%s", cmd.name, cmd.usage)
}
func (cmdSet commandSet) insert(name, usage string) error {
	cmd := command{name: name, usage: usage, fs: flag.NewFlagSet(name, flag.ExitOnError)}
	if _, ok := cmdSet[name]; ok {
		return errors.New(fmt.Sprintf("'%s' already exists as a command", name))
	}
	cmdSet[name] = cmd
	return nil
}
func (cmdSet commandSet) String() string {
	ret := "Flags...\n"
	flag.VisitAll(func(f *flag.Flag) {
		ret += formatFlagHelp(*f)
	})
	ret += "\n\nCommands..."
	for _, cmd := range cmdSet {
		ret += formatCommandHelp(cmd)
	}
	ret += "\n"
	return ret
}
func loadCommandSet() commandSet {
	cmdSet := commandSet{}
	cmdSet.insert("add", "command adds a new ksat")
	cmdSet["add"].fs.String("name", "", "name for new ksat (Required)")
	cmdSet["add"].fs.String("usage", "", "usage for new ksat (Required)")
	cmdSet["add"].fs.String("prompt", "", "prompt for new ksat (Required)")
	return cmdSet
}
func parseTopLevelFlags(cmdSet commandSet) ([]string, error) {
	version := flag.Bool("version", false, "version for program")
	help := flag.Bool("h", false, "help for program")
	flag.Parse()
	args := flag.Args()
	if *help && *version {
		return nil, errors.New("only one flag must be specified")
	}
	if *help && len(args) != 0 {
		return nil, errors.New("help flag only applicable when no following args followed")
	}
	if *version {
		fmt.Println("ksat version 0.0.0")
	}
	if *help {
		fmt.Println(cmdSet)
	}
	return args, nil
}
func parseCommands(cmdSet commandSet, args []string) error {
	if len(args) == 0 {
		return nil
	}
	cmd, ok := cmdSet[args[0]]
	if !ok {
		return errors.New(fmt.Sprintf("'%s' is an unknown subcommand", args[0]))
	}
	return cmd.fn(cmd.fs, args[1:])
}
func addFn(flagSet *flag.FlagSet, args []string) error {
	if len(args) == 0 {
		return errors.New("No interactivity mode yet")
	}
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	flags := map[string]*flag.Flag{}
	flagSet.VisitAll(func(flag *flag.Flag) { flags[flag.Name] = flag })
	for _, flag := range flags {
		if strings.Contains(flag.Usage, "(Required)") && flag.Value.String() == "" {
			flagSet.PrintDefaults()
			return nil
		}
	}
	return newKsat(
		flags["name"].Value.String(),
		flags["usage"].Value.String(),
		flags["prompt"].Value.String(),
	)
}
*/
