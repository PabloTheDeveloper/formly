package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pablothedeveloper/formly"
)

func main() {
	env, err := formly.NewLocalSqLiteEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()
	flag.CommandLine.Usage = func() {
		fmt.Println("form [--help] [command] [--help] [<args>] [<options>]")
		fmt.Println("commands:\n" +
			"  create\t- creates a form\n" +
			"  delete\t- deletes an existing form\n" +
			"  label\t\t- adds a label to an existing form\n" +
			"  review\t- reviews content of an existing form\n" +
			"  submit\t- submits an existing form\n" +
			"  submissions\t- views prior submissions of a form")
	}
	flag.Parse()
	if flag.NArg() == 0 && flag.NFlag() == 0 {
		fmt.Println("default behavior...done!")
		return
	}
	cmd := flag.NewFlagSet(flag.Arg(0), flag.ExitOnError)
	cmd.Usage = func() {
		switch cmd.Name() {
		case "create":
			fmt.Println("usage: form create <form-name> <form-usage>")
		case "delete":
			fmt.Println("usage: form delete <form-name>")
		case "label":
			fmt.Println("usage: form label <form-name> [--repeatable] <label-name> <label-usage>")
		case "submit":
			fmt.Println("usage: form submit <form-name>")
		case "submissions":
			fmt.Println("usage: form submissions <form-name>")
		case "modify":
			fmt.Println(
				"usage: form modify <form-name> [--name] [--usage]" +
					"<label-name> [--name] [--usage]",
			)
		default:
			flag.CommandLine.Usage()
		}
	}
	cmd.Parse(flag.Args()[1:])
	switch cmd.Name() {
	case "create":
		if cmd.NArg() < 2 {
			fmt.Println("fatal: Must specify a form name and form usage")
			cmd.Usage()
			return
		}
		if err := create(env, cmd.Arg(0), cmd.Arg(1)); err != nil {
			fmt.Println(err)
		}
	case "delete":
		subcmd, err := newSubCommand(env, cmd, cmd.Args()...)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := delete(env, subcmd.form.ID); err != nil {
			fmt.Println(err)
		}
	case "label":
		subcmd, err := newSubCommand(env, cmd, cmd.Args()...)
		if err != nil {
			fmt.Println(err)
			return
		}
		repeatable := subcmd.fs.Bool("repeatable", false, "whether label repeats")
		subcmd.parse()
		name := subcmd.fs.Arg(0)
		usage := subcmd.fs.Arg(1)
		if subcmd.fs.NArg() < 2 {
			fmt.Println("fatal: Must specify a label name and label usage")
			cmd.Usage()
			return
		}
		if err := label(env, subcmd.form.ID, int64(len(subcmd.labels)+1), *repeatable, name, usage); err != nil {
			fmt.Println(err)
		}
	case "review":
		subcmd, err := newSubCommand(env, cmd, cmd.Args()...)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("form created: %v\n", subcmd.form)
		fmt.Printf("labels created: %v\n", subcmd.labels)
	case "submit":
		subcmd, err := newSubCommand(env, cmd, cmd.Args()...)
		if err != nil {
			fmt.Println(err)
			return
		}
		subcmd.setupFormFlags()
		if err := subcmd.parseFormFlags(); err != nil {
			fmt.Println(err)
			return
		}
		if err := subcmd.submitForm(env); err != nil {
			fmt.Println(err)
			return
		}
	case "submissions":
		subcmd, err := newSubCommand(env, cmd, cmd.Args()...)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := submissions(env, subcmd.form.ID); err != nil {
			fmt.Println(err)
		}
	case "modify":
		subcmd, err := newSubCommand(env, cmd, cmd.Args()...)
		if err != nil {
			fmt.Println(err)
			return
		}
		newName := subcmd.fs.String("name", subcmd.form.Name, "new name for form")
		newUsage := subcmd.fs.String("usage", subcmd.form.Usage, "new usage for form")
		subcmd.parse()

		if err := modify(env, subcmd.form.ID, *newName, *newUsage); err != nil {
			fmt.Println(err)
			return
		}
		if subcmd.fs.NArg() == 0 {
			return
		}
		found := -1
		for i, label := range subcmd.labels {
			if label.Name == subcmd.fs.Arg(0) {
				found = i
				break
			}
		}
		if found == -1 {
			fmt.Printf("label '%s' for form '%s' not found\n", subcmd.fs.Arg(0), subcmd.fs.Name())
			return
		}
		subsubcmd := flag.NewFlagSet(subcmd.fs.Arg(0), flag.ExitOnError)
		newLabelName := subsubcmd.String("name", subcmd.labels[found].Name, "new name for label")
		newLabelUsage := subsubcmd.String("usage", subcmd.labels[found].Usage, "new usage for label")
		position := subsubcmd.Int64("position", subcmd.labels[found].Position, "new position for label")
		repeatable := subsubcmd.Bool("repeatable", false, "whether label repeats")
		subsubcmd.Parse(subcmd.fs.Args()[1:])
		if err := modifylabel(env, subcmd.form.ID, subcmd.labels[found].ID, *position, *repeatable, *newLabelName, *newLabelUsage); err != nil {
			fmt.Println(err)
			return
		}
	case "":
		fmt.Println("No command passed in")
	default:
		fmt.Printf("subcommand '%s' does not exist\n", cmd.Name())
	}
}

type subcommand struct {
	fs                     *flag.FlagSet
	form                   formly.Form
	labels                 []formly.Label
	flags                  []formFlag
	entries                []formly.Entry
	unParsedArgs           []string
	repeatableArgSeperator string
}
type formFlag struct {
	labelID    int64
	repeatable bool
	name, txt  string
}

func newSubCommand(env *formly.Env, cmd *flag.FlagSet, args ...string) (scmd subcommand, err error) {
	if cmd.NArg() == 0 && cmd.NFlag() == 0 {
		err = fmt.Errorf("default behavior for cmd...done")
		return
	}
	scmd.form, err = env.FormModel.GetByName(args[0])
	if err != nil {
		return
	}
	scmd.labels, err = env.LabelModel.GetLabels(scmd.form.ID)
	if err != nil {
		return
	}
	scmd.entries = make([]formly.Entry, 0)
	scmd.fs = flag.NewFlagSet(args[0], flag.ExitOnError)
	scmd.unParsedArgs = args[1:]
	scmd.repeatableArgSeperator = ",/"
	return
}

var errRepeatableFlagSeperator = errors.New("flag contains a seperator while not being repeatable")

func (scmd *subcommand) parse() {
	scmd.fs.Parse(scmd.unParsedArgs)
}
func (scmd *subcommand) setupFormFlags() {
	for i, label := range scmd.labels {
		scmd.flags = append(
			scmd.flags,
			formFlag{labelID: label.ID, repeatable: label.Repeatable, name: label.Name, txt: ""},
		)
		scmd.fs.StringVar(
			&scmd.flags[i].txt,
			scmd.labels[i].Name,
			"",
			scmd.labels[i].Usage,
		)
	}
}
func (scmd *subcommand) parseFormFlags() error {
	if err := scmd.fs.Parse(scmd.unParsedArgs); err != nil {
		return err
	}
	for _, fflag := range scmd.flags {
		if !fflag.repeatable && strings.Contains(fflag.txt, scmd.repeatableArgSeperator) {
			return errRepeatableFlagSeperator
		}
	}
	if scmd.fs.NFlag() != 0 {
		return nil
	}
	for i, flag := range scmd.flags {
		s := bufio.NewScanner(os.Stdin)
		inputs := []string{}
		prompt := flag.name + ":\n"
		for fmt.Print(prompt); s.Scan(); fmt.Print(prompt) {
			txt := s.Text()
			if strings.Contains(txt, scmd.repeatableArgSeperator) {
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
		scmd.flags[i].txt = strings.Join(inputs, scmd.repeatableArgSeperator)
		if err := s.Err(); err != io.EOF && err != nil {
			return err
		}
	}
	return nil
}
func (scmd *subcommand) submitForm(env *formly.Env) error {
	submission, err := env.SubmissionModel.Create(scmd.form.ID)
	if err != nil {
		return err
	}
	fmt.Println(
		fmt.Sprintf("form '%s' submitted at time:%v", scmd.fs.Name(), submission.CreateAt),
	)
	for _, flag := range scmd.flags {
		if flag.txt == "" {
			continue
		}
		for _, arg := range strings.Split(flag.txt, scmd.repeatableArgSeperator) {
			entry, err := env.EntryModel.Create(submission.ID, flag.labelID, arg)
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("\t%s: %s", flag.name, entry.Txt))
		}
	}
	return nil
}
func create(env *formly.Env, name, usage string) error {
	form, err := env.FormModel.Create(name, usage)
	if err != nil {
		return err
	}
	fmt.Printf("form created: %v\n", form)
	return nil
}
func delete(env *formly.Env, formID int64) error {
	form, err := env.FormModel.DeleteByID(formID)
	if err != nil {
		return err
	}
	fmt.Printf("deleted form: %v\n", form)
	return nil
}
func label(env *formly.Env, formID, position int64, repeatable bool, name, usage string) error {
	fmt.Println("trying to add label...")
	label, err := env.LabelModel.Create(formID, position, repeatable, name, usage)
	if err != nil {
		return err
	}
	fmt.Printf("label created: %v\n", label)
	return nil
}
func submissions(env *formly.Env, formID int64) error {
	submissions, err := env.SubmissionModel.GetSubmissions(formID)
	if err != nil {
		return err
	}
	labels, err := env.LabelModel.GetLabels(formID)
	if err != nil {
		return err
	}
	if len(submissions) == 0 {
		fmt.Println("no submission for this form yet")
		return nil
	}
	for _, submission := range submissions {
		fmt.Printf("submission:%v\n", submission.CreateAt)
		for _, label := range labels {
			entries, err := env.GetEntries(submission.ID, label.ID)
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				continue
			}
			for _, entry := range entries {
				fmt.Printf("\t%s: %s \n", label.Name, entry.Txt)
			}
		}
	}
	return nil
}
func modify(env *formly.Env, formID int64, newName, newUsage string) error {
	form, err := env.FormModel.Update(formID, newName, newUsage)
	if err != nil {
		return err
	}
	fmt.Printf("updated form: %v\n", form)
	return nil
}
func modifylabel(env *formly.Env, formID, labelID, position int64, repeatable bool, newName, newUsage string) error {
	labels, err := env.LabelModel.Update(formID, labelID, position, repeatable, newName, newUsage)
	if err != nil {
		return err
	}
	fmt.Printf("updated label(s): %v\n", labels)
	return nil
}
