package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

var data *Data

func main() {
	data = &Data{map[string]Ksat{}}
	err := getData()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd, err := getCommand(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg, err := execCommand(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msg)
}

func getCommand(rawArgs []string) (string, error) {
	switch len(rawArgs) {
	case 0:
		return "", errors.New("No arguments passed in 'os.Args'")
	case 1:
		return "", nil
	case 2:
		return rawArgs[1], nil
	default:
		return "", errors.New("Extra arguments passed in")
	}
}

func execCommand(cmd string) (string, error) {
	switch cmd {
	case "", "-h", "--help":
		return "there are four commands: read, create, update, delete", nil
	case "commands":
		fmt.Println("Commands:")
		for k := range data.Ksats {
			fmt.Println(k)
		}
		return "", nil
	case "read":
		fmt.Println("Commands:")
		for k := range data.Ksats {
			fmt.Println(k)
		}
		scanner := bufio.NewScanner(os.Stdin)
		cmd := input("enter command:", scanner)
		entries := dbReadEntries(cmd)
		if len(entries) == 0 {
			return "no ksat entries or ksat command", errors.New("no file")
		}

		for _, entry := range entries {
			for _, pair := range entry.InputPairs {
				fmt.Println(pair[0] + ":")
				fmt.Println(pair[1] + "\n")
			}
		}

	case "create":
		ksat, err := promptKsat()
		if err != nil {
			return "Error in prompting ksat creation", err
		}
		if err := dbCreateKsat(ksat); err != nil {
			return "Error in db creation of ksat", err
		}
	default:
		// fetch commands and check if valid
		ksat, ok := data.Ksats[cmd]
		if !ok {
			return "", errors.New("Command not found")
		}
		executeKsat(ksat)
	}
	return "", nil
}
