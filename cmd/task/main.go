package main

import (
	"bufio"
	"encoding/json"
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
		if err := s.Err(); err != io.EOF && err != nil {
			return err
		}
	}
	return nil
}

// submitForm execute submits the form and creates a submission and all entries. All entries with empty text are not submitted
func (cmd *command) submitForm(env *ksat.Env) error {
	submission, err := env.SubmissionModel.Create(cmd.formID)
	fmt.Println(
		fmt.Sprintf("form '%s' submitted at time:%v", cmd.fs.Name(), submission.GetCreateAt()),
	)
	if err != nil {
		return err
	}
	for _, flag := range cmd.flags {
		if flag.arg == "" {
			continue
		}
		for _, arg := range strings.Split(flag.arg, cmd.flagArgSeperator) {
			entry, err := env.EntryModel.Create(submission.GetID(), flag.labelID, arg)
			if err != nil {
				return err
			}
			fmt.Print(fmt.Sprintf("\t%s: %s", flag.name, entry.GetTxt()))
		}
	}
	return nil
}

// newForm a new form as a subcommand. If the labels flag is empty, creates it without labels
func newForm(cmd *command, env *ksat.Env) error {
	type jsonLabel struct {
		Repeatable  bool // `json:"omitempty"`
		Name, Usage string
	}

	// validate form
	if err := ksat.ValidateName(cmd.flags[0].arg); err != nil {
		fmt.Println("create:" + cmd.flags[0].arg + "'")
		return err
	}
	if err := ksat.ValidateUsage(cmd.flags[1].arg); err != nil {
		fmt.Println("create:", cmd.flags[1].arg)
		return err
	}

	// validate labels
	var labels []jsonLabel
	if cmd.flags[2].arg != "" {
		if err := json.Unmarshal([]byte(cmd.flags[2].arg), &labels); err != nil {
			return err
		}
	}
	names := map[string]bool{}
	for _, label := range labels {
		if err := ksat.ValidateName(label.Name); err != nil {
			fmt.Println("create:label:", label.Name)
			return err
		}
		if err := ksat.ValidateUsage(label.Usage); err != nil {
			fmt.Println("create:label", label.Usage)
			return err
		}
		if _, ok := names[label.Name]; !ok {
			names[label.Name] = true
		} else {
			fmt.Println("create:label:duplicate name not allowed", label.Usage)
			return fmt.Errorf("name '%s' was used in at least two seperate labels objects", label.Name)
		}
	}

	// create form and label
	form, err := env.FormModel.Create(cmd.flags[0].arg, cmd.flags[1].arg)
	if err != nil {
		return err
	}
	for i, label := range labels {
		if _, err := env.LabelModel.Create(
			form.GetID(),
			int64(i+1),
			label.Repeatable,
			label.Name,
			label.Usage,
		); err != nil {
			return err
		}
	}
	return nil
}

// viewForm views a form's entries. Ensure that form name passed in is valid.
func viewForm(cmd *command, env *ksat.Env) error {
	if err := ksat.ValidateName(cmd.flags[0].arg); err != nil {
		fmt.Println("read:" + cmd.flags[0].arg)
		return err
	}
	form, err := env.FormModel.GetByName(cmd.flags[0].arg)
	if err != nil {
		return err
	}
	if form == (ksat.Form{}) {
		return fmt.Errorf("form '%s' does not exist", cmd.fs.Name())
	}
	submissions, err := env.SubmissionModel.GetSubmissions(form.GetID())
	if err != nil {
		return err
	}
	labels, err := env.LabelModel.GetLabels(form.GetID())
	if err != nil {
		return err
	}
	if len(submissions) == 0 {
		fmt.Println("no submission for this form yet")
		return nil
	}
	for _, submission := range submissions {
		fmt.Printf("submission:%v\n", submission.GetCreateAt())
		for _, label := range labels {
			entries, err := env.GetEntries(submission.GetID(), label.GetID())
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				continue
			}
			for _, entry := range entries {
				fmt.Printf("\t%s: %s \n", label.GetName(), entry.GetTxt())
			}
		}
	}
	return nil
}

// removeForm remove a form, labels, and entries. Ensure that form name passed in is valid.
func removeForm(cmd *command, env *ksat.Env) error {
	if err := ksat.ValidateName(cmd.flags[0].arg); err != nil {
		fmt.Println("read:" + cmd.flags[0].arg)
		return err
	}
	if err := env.FormModel.DeleteByName(cmd.flags[0].arg); err != nil {
		return err
	}
	fmt.Printf("removed form '%s' successfully!\n", cmd.fs.Name())
	return nil
}

/*
func updateForm(cmd *command, env *ksat.Env) error {
	// ONLY ACCEPT ONE FLAG

	// update form
	type jsonLabel struct {
		id         int64
		Name    string
		Usage      string
		Repeatable bool
	}
	type jsonForm struct {
		Name, NewUsage, NewName string
		UpdateLabel jsonLabel
		DeleteLabels    []jsonLabel
	}
	var jform jsonForm
	if cmd.flags[0].arg == "" {
		return errors.New("nothing passed into form flag for update command")
	}
	if err := json.Unmarshal([]byte(cmd.flags[0].arg), &jform); err != nil {
		return err
	}
	// validate the update of the form
	var formID int64
	if jform.Name != "" {
		form, err := env.FormModel.GetByName(jform.Name)
		if err != nil {
			return err
		}
		if form == (ksat.Form{}) {
			return errors.New("form name does not exists")
		}
		if jform.NewName != "" {
			jform.NewName = form.GetName()
		}
		if jform.NewUsage != "" {
			jform.NewUsage = form.GetUsage()
		}
		formID = form.GetID()
	}
	if err := ksat.ValidateName(jform.NewName); err != nil {
		fmt.Println("form newName isn't valid")
		return err
	}
	if err := ksat.ValidateUsage(jform.NewUsage); err != nil {
		fmt.Println("form newUsage isn't valid")
		return err
	}
	if formID == 0 {
		form, err := env.FormModel.Create(jform.NewName, jform.NewUsage)
		if err != nil {
			return err
		}
		formID = form.GetID()
	} else {
		// TODO env.FormModel.Update(form.GetID(), jform.NewName, jform.NewUsage)
	}
	if len(jform.DeleteLabels) == 0 && jform.UpdateLabel {
		return nil
	}
	if len(jform.DeleteLabels) != 0 && len(jform.UpdateLabels) != 0 {
		return errors.New("you can only pick either to delete labels or update labels")
	}
	if len(jform.DeleteLabels) != 0 {
		names := map[string]bool{}
		for _, label := range jform.DeleteLabels {
			if err := ksat.ValidateName(label.Name); err != nil {
				fmt.Println("a label called '"+label.Name+"'")
				return err
			}
			if _, ok := names[label.Name]; ok {
				return errors.New("duplicated name for labels in deleteLabels")
			}
			names[label.Name] = true
		}
		// TODO delete label env.LabelModel
		return nil
	}
	if len(jform.UpdateLabels) != 0 {
		names := map[string]bool{}
		for _, label := range jform.UpdateLabels {
			if err := ksat.ValidateName(label.Name); err != nil {
				fmt.Println("a label called '"+label.Name+"'")
				return err
			}
			if _, ok := names[label.Name]; ok {
				return errors.New("duplicated name for labels in deleteLabels")
			}
			names[label.Name] = true
		}
		for _, label := range jform.UpdateLabels {
			if err := ksat.ValidateName(label.Name); err != nil {
				fmt.Println("a label called '"+label.Name+"'")
				return err
			}
			if _, ok := names[label.Name]; ok {
				return errors.New("duplicated name for labels in deleteLabels")
			}
			names[label.Name] = true
		}
	}

	// check if form name != "", if so, do below, otherwise try to create a new form (will do in the update)
		// check if form name is valid
		// check if form name exist
		// assign usage of form to usage of jform if NewUsage was empty
		// assign name of form to name of jform if NewName was empty

	// check if form new name is valid (does not already exists)
	// check if form new name is valid (syntax)
	// check if form newUsage is valid

	// get all labels for form
	// for each label attempted to be update, called jlabels,  we will do the following
}
*/

// main program
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
	// process commands
	switch cmd.fs.Name() {
	case "new":
		if err := newForm(cmd, env); err != nil {
			log.Fatal(err)
		}
	case "view":
		if err := viewForm(cmd, env); err != nil {
			log.Fatal(err)
		}
	case "remove":
		if err := removeForm(cmd, env); err != nil {
			log.Fatal(err)
		}
	default:
		if err := cmd.submitForm(env); err != nil {
			log.Fatal(err)
		}
	}
}
