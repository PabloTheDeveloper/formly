package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
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
	case "read":
		return "", nil
	case "create":
		ksat, err := promptKsat()
		if err != nil {
			return "Error in prompting ksat creation", err
		}
		if err := dbCreateKsat(ksat); err != nil {
			return "Error in db creation of ksat", err
		}
	case "update":
		return "", nil
	case "delete":
		return "", nil
	default:
		// fetch commands and check if valid
		return "", errors.New("Command not found")
	}
	return "", nil
}

// Ksat contains task information
type Ksat struct {
	identifier string
	prompt     string
	label      string
}

func promptKsat() (ksat Ksat, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	ksat.identifier = input("Enter command identifier: ", scanner)
	ksat.prompt = input("Enter prompt:\n", scanner)
	ksat.label = input("Enter input label: ", scanner)
	return ksat, nil
}

func input(prompt string, scanner *bufio.Scanner) string {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text()
}

func dbCreateKsat(ksat Ksat) error {
	// TODO
	return nil
}
