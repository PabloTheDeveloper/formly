package main

import (
	"fmt"
)

type command struct {
	ksat
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
	cmd := &command{}
	cmd.ksat.name = args[1]
	if err := cmd.ksat.validateName(); err != nil {
		return nil, err
	}
	if err := cmd.ksat.getByName(); err != nil {
		return nil, err
	}
	return cmd, nil
}
