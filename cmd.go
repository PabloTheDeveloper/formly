package main

import (
	"errors"
	"flag"
	"fmt"
)

func execute() error {
	args, err := topLevelFlags()
	if err != nil {
		return err
	}
	err = topLevelCommands(args)
	return err
}
func topLevelFlags() ([]string, error) {
	version := flag.Bool("version", false, "version for program")
	if *version == true {
		fmt.Println("ksat version 0.0.0")
		return nil, nil
	}
	flag.Parse()
	return flag.Args(), nil
}
func topLevelCommands(rawArgs []string) error {
	if len(rawArgs) == 0 {
		return nil
	}
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	switch rawArgs[0] {
	case "new":
		// 'new' flags
		name := newCmd.String("name", "", "name for new ksat (Required)")
		desc := newCmd.String("desc", "", "description for new ksat (Required)")
		// no flags
		if len(rawArgs) == 1 {
			newCmd.PrintDefaults()
			break
		}
		if err := newCmd.Parse(rawArgs[1:]); err != nil {
			break
		}
		if *name == "" || *desc == "" {
			newCmd.PrintDefaults()
			break
		}
		if ksatId, err := getKsatIdByName(*name); err != nil {
			return err
		} else if ksatId == -1 {
			_, err := newKsat(*name, *desc)
			return err
		} else {
			return errors.New("ksat name already exists")
		}
	default:
		return errors.New("unknown subcommand")
	}
	return nil
}
