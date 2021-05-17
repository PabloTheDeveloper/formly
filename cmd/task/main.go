package main

import (
	"bufio"
	"errors"
	stdFlag "flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pablothedeveloper/ksat"
)

type command struct {
	formID  int64
	usage   string
	fs      *stdFlag.FlagSet
	rawArgs []string
	flags   []flag
}
type flag struct {
	labelID            int64
	name, usage, value string
}

var errNoArgNewCommand = errors.New("No arguments passed into newCommand")

// newCommand returns a command with populated fields unless an error occurs along the way
func newCommand(env *ksat.Env, rawArgs ...string) (*command, error) {
	// new command
	if len(rawArgs) == 0 {
		return nil, errNoArgNewCommand
	}
	cmd := &command{fs: stdFlag.NewFlagSet(rawArgs[0], stdFlag.ExitOnError), rawArgs: rawArgs[1:]}

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
		cmd.flags = append(cmd.flags, flag{labelID: label.GetID(), name: label.GetName(), usage: label.GetUsage(), value: ""})
		cmd.fs.StringVar(
			&cmd.flags[i].value,
			cmd.flags[i].name,
			"",
			cmd.flags[i].usage,
		)
	}
	return cmd, nil
}

var errUnitializedCommand = errors.New("command was created via 'newCommand'")

func (cmd *command) process(env *ksat.Env) error {
	// check command was initialized correctly
	if cmd.formID == 0 || cmd.fs == nil || cmd.flags == nil {
		return errUnitializedCommand
	}
	// fill in flags from command line
	cmd.fs.Parse(cmd.rawArgs)

	submission, err := env.SubmissionModel.Create(cmd.formID)
	if err != nil {
		return err
	}

	// interactive mode to assign flags's values
	if cmd.fs.NFlag() == 0 {
		for i := 0; i < len(cmd.flags); i++ {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("(flag usage)\n%s\n\n(flag name)\n%s:\n", cmd.flags[i].usage, cmd.flags[i].name)
			scanner.Scan()
			cmd.flags[i].value = scanner.Text()
			fmt.Println()
		}
	}

	// create entries through flags
	for _, flag := range cmd.flags {
		entry, err := env.EntryModel.Create(submission.GetID(), flag.labelID, flag.value)
		if err != nil {
			return err
		}
		fmt.Println("entry:", entry)
	}
	return nil
}

type helpable interface {
	GetName() string
	GetUsage() string
}

func main() {
	env, err := ksat.NewEnv(ksat.LocalSqlite)
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
		helpables := make([]helpable, len(forms))
		for i, form := range forms {
			helpables[i] = form
		}
		if len(helpables) == 0 {
			fmt.Println("main: no subcommands found...")
			return
		}
		for _, item := range helpables {
			fmt.Printf("%s\n	-%s\n", item.GetName(), item.GetUsage())
		}
	}
	stdFlag.Parse()

	// process subcommand
	cmd, err := newCommand(env, os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.process(env); err != nil {
		log.Fatal(err)
	}
}
