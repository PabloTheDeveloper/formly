package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type command struct {
	usage string
	fs    *flag.FlagSet
	fn    func(*flag.FlagSet, []string) error
}
type commandSet map[string]command

func (cmdSet commandSet) insert(fn func(*flag.FlagSet, []string) error, name, usage string) error {
	cmd := command{fmt.Sprintf("  %s\n\t%s", name, usage), flag.NewFlagSet(name, flag.ExitOnError), fn}
	if _, ok := cmdSet[name]; ok {
		return errors.New(fmt.Sprintf("'%s' already exists as a command", name))
	}
	cmdSet[name] = cmd
	return nil
}
func (cmdSet commandSet) printHelp() {
	fmt.Println("\nFlags...")
	flag.PrintDefaults()
	fmt.Println("\nCommands...")
	for _, cmd := range cmdSet {
		fmt.Println(cmd.usage)
	}
	fmt.Println()
}
func execute() error {
	version := flag.Bool("version", false, "version for program")
	help := flag.Bool("h", false, "help for program")
	flag.Parse()
	if *version == true {
		fmt.Println("ksat version 0.0.0")
		return nil
	}
	cmdSet := commandSet{}
	// add flagset
	cmdSet.insert(addFn, "add", "command adds a new ksat")
	cmdSet["add"].fs.String("name", "", "name for new ksat (Required)")
	cmdSet["add"].fs.String("usage", "", "usage for new ksat (Required)")
	cmdSet["add"].fs.String("prompt", "", "prompt for new ksat (Required)")
	cmdSet["add"].fs.String("sched", "", "schedule for new ksat (Required)")

	args := flag.Args()
	if len(args) == 0 {
		if *help {
			cmdSet.printHelp()
		}
		return nil
	}
	// execute command
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
		flags["sched"].Value.String(),
		flags["prompt"].Value.String(),
	)
}
