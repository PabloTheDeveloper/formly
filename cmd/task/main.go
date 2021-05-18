package main

import (
	"bufio"
	"errors"
	stdFlag "flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pablothedeveloper/ksat"
)

type command struct {
	formID           int64
	fs               *stdFlag.FlagSet
	usage            string
	flagArgSeperator string
	rawArgs          []string
	flags            []flag
}
type flag struct {
	labelID     int64
	repeatable  bool
	name, usage string
	arg         string
}

var errNoArgNewCommand = errors.New("No arguments passed into newCommand")

// newCommand returns a command with populated fields unless an error occurs along the way
func newCommand(env *ksat.Env, rawArgs ...string) (*command, error) {
	// new command
	if len(rawArgs) == 0 {
		return nil, errNoArgNewCommand
	}
	cmd := &command{
		fs:               stdFlag.NewFlagSet(rawArgs[0], stdFlag.ExitOnError),
		flagArgSeperator: ",/",
		rawArgs:          rawArgs[1:],
	}

	// load command
	form, err := env.FormModel.GetByName(cmd.fs.Name())
	if err != nil {
		return nil, err
	}
	if form == (ksat.Form{}) {
		return nil, fmt.Errorf("form '%s' does not exist", cmd.fs.Name())
	}
	cmd.formID = form.GetID()
	cmd.usage = form.GetUsage()

	// load flags of command
	labels, err := env.LabelModel.GetLabels(form.GetID())
	if err != nil {
		return nil, err
	}
	for i, label := range labels {
		cmd.flags = append(
			cmd.flags,
			flag{
				labelID:    label.GetID(),
				repeatable: label.GetRepeatable(),
				name:       label.GetName(),
				usage:      label.GetUsage(),
				arg:        "",
			},
		)
		cmd.fs.StringVar(
			&cmd.flags[i].arg,
			cmd.flags[i].name,
			"",
			cmd.flags[i].usage,
		)
	}
	return cmd, nil
}

var errRepeatableFlagSeperator = errors.New("flag contains a seperator while not being repeatable")

// parse gets flags from command and insert into the command struct. activates interactivity mode when no flags are passed in.
func (cmd *command) parse() error {
	if err := cmd.fs.Parse(cmd.rawArgs); err != nil {
		return err
	}
	for _, flag := range cmd.flags {
		if !flag.repeatable && strings.Contains(flag.arg, cmd.flagArgSeperator) {
			return errRepeatableFlagSeperator
		}
	}
	if cmd.fs.NFlag() != 0 {
		return nil
	}
	for i, flag := range cmd.flags {
		s := bufio.NewScanner(os.Stdin)
		inputs := []string{}
		prompt := flag.name + ":\n"
		for fmt.Print(prompt); s.Scan(); fmt.Print(prompt) {
			txt := s.Text()
			if strings.Contains(txt, cmd.flagArgSeperator) {
				fmt.Printf("input cannot contain '%s'\n", cmd.flagArgSeperator)
				continue
			}
			if txt == "" {
				break
			}
			inputs = append(inputs, txt)
			if !flag.repeatable {
				break
			}
		}
		cmd.flags[i].arg = strings.Join(inputs, cmd.flagArgSeperator)
		if err := s.Err(); err != io.EOF {
			return err
		}
	}
	return nil
}

// execute submits the form and creates a submission and all entries.
func (cmd *command) execute(env *ksat.Env) error {
	submission, err := env.SubmissionModel.Create(cmd.formID)
	if err != nil {
		return err
	}
	for _, flag := range cmd.flags {
		if flag.arg == "" {
			continue
		}
		fmt.Println(strings.Split(flag.arg, cmd.flagArgSeperator))
		for _, arg := range strings.Split(flag.arg, cmd.flagArgSeperator) {
			entry, err := env.EntryModel.Create(submission.GetID(), flag.labelID, arg)
			if err != nil {
				return err
			}
			fmt.Println("entry:", entry)
		}
	}
	return nil
}
func main() {
	env, err := ksat.NewLocalSqLiteEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()

	// setup help flag for command
	stdFlag.CommandLine.Usage = func() {
		forms, err := env.FormModel.GetAll()
		if err != nil {
			log.Fatal("main:", err)
		}
		if len(forms) == 0 {
			fmt.Println("main: no subcommands found...")
			return
		}
		for _, item := range forms {
			fmt.Printf("%s\n	-%s\n", item.GetName(), item.GetUsage())
		}
	}
	stdFlag.Parse()

	// process subcommand
	cmd, err := newCommand(env, os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.parse(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.execute(env); err != nil {
		log.Fatal(err)
	}
}
