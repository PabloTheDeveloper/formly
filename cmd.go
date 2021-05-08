package main

import (
	"flag"
	"fmt"
)

type command struct {
	Ksat
	session
	*flag.FlagSet
	flagArgs []string
	entries  []entry
}

func helpMessage() error {
	fmt.Println("help message")
	return nil
}
func newCommand(args []string) (*command, error) {
	if len(args) <= 1 {
		return nil, helpMessage()
	}
	// args has len of >=1
	ksat, err := GetKsatByName(args[1])
	if err != nil {
		return nil, err
	}
	return &command{
		Ksat:     ksat,
		FlagSet:  flag.NewFlagSet(args[1], flag.ExitOnError),
		flagArgs: args[2:],
	}, nil
}
func (command *command) executeCommand() error {
	prompts, err := GetPromptsByKsatID(command.Ksat.id)
	if err != nil {
		return err
	}
	// gets them in order of creation (due to how it would be created)
	for i, prompt := range prompts {
		if err := prompt.getByID(); err != nil {
			return err
		}
		command.entries = append(command.entries, entry{promptID: prompt.id})
		command.FlagSet.StringVar(
			&command.entries[i].txt,
			prompt.flag,
			"",
			prompt.usage,
		)
	}
	if err := command.FlagSet.Parse(command.flagArgs); err != nil {
		return err
	}

	session := session{KsatID: command.Ksat.id}
	if err := session.dbInsert(); err != nil {
		return err
	}
	// relies on the order of creation in order for the prompt to be correct
	// also they are the same length. Creates an entry for entries which are not filled in
	// maybe I can make the string "" a empty value or response value
	for i := range command.entries {
		command.entries[i].sessionID = session.id
		// command.entries[i].txt = command.FlagSet.Lookup("f").Value.String()
		if err := command.entries[i].dbInsert(); err != nil {
			return err
		}
	}
	return nil
}
