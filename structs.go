package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

// Data contains data for application
type Data struct {
	Ksats map[string]Ksat
}

func getData() error {
	dataFile, err := os.OpenFile("data.gob", os.O_RDONLY, 0755)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist!")
		return nil
	}

	if err != nil {
		return err
	}

	if err := gob.NewDecoder(dataFile).Decode(data); err != nil {
		return nil
	}

	dataFile.Close()
	return nil
}

// Ksat contains task information
type Ksat struct {
	Identifier string
	Prompt     string
	Label      string
}

func input(prompt string, scanner *bufio.Scanner) string {
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text()
}

func promptKsat() (ksat Ksat, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	ksat.Identifier = input("Enter command identifier: ", scanner)
	ksat.Prompt = input("Enter prompt:\n", scanner)
	ksat.Label = input("Enter input label: ", scanner)
	return ksat, nil
}

func dbCreateKsat(ksat Ksat) error {
	data.Ksats[ksat.Identifier] = ksat
	dataFile, err := os.OpenFile("data.gob", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if err := gob.NewEncoder(dataFile).Encode(*data); err != nil {
		return err
	}
	dataFile.Close()
	return nil
}

// Entry contains Ksat entry
type Entry struct {
	InputPairs [][2]string
}

func executeKsat(ksat Ksat) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(ksat.Prompt)
	value := input(ksat.Label+":", scanner)
	entry := Entry{[][2]string{
		[2]string{ksat.Label, value},
	}}
	entries := dbReadEntries(ksat.Identifier)
	entries = append(entries, entry)
	dbCreateEntry(ksat.Identifier, entries)
}

func dbReadEntries(command string) []Entry {
	entries := []Entry{}
	dataFile, err := os.OpenFile(command+".gob", os.O_RDONLY, 0755)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist!")
		return entries
	}
	if err != nil {
		return entries
	}
	if err := gob.NewDecoder(dataFile).Decode(&entries); err != nil {
		return entries
	}
	dataFile.Close()
	return entries
}

func dbCreateEntry(command string, entries []Entry) error {
	dataFile, err := os.OpenFile(command+".gob", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if err := gob.NewEncoder(dataFile).Encode(entries); err != nil {
		return err
	}
	dataFile.Close()
	return nil

}
