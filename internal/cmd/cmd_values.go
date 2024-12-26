package cmd

import (
	"fmt"
	"os/exec"

	"github.com/direnv/direnv/v2/gzenv"
)

// CmdValue represents a single recorded command value
type CmdValue struct {
	Cmd     string `json:"command"`
	Value   string `json:"value"`
}

// CmdValues represent a record of all the known commands to watch
type CmdValues struct {
	list *[]CmdValue
}

// NewCmdValues creates a new empty CmdValues
func NewCmdValues() (values CmdValues) {
	list := make([]CmdValue, 0)
	values.list = &list
	return
}

// Update gets the latest values and updates the record.
func (values *CmdValues) Update(cmd string, config *Config) (err error) {
	var output string

	output, err = runCmd(cmd, config)

	err = values.NewCmd(cmd, output)

	return
}

// NewCmd adds the command to watch, with the given output
func (values *CmdValues) NewCmd(cmd string, value string) (err error) {
	var cmdValue *CmdValue

	for idx := range *(values.list) {
		if (*values.list)[idx].Cmd == cmd {
			cmdValue = &(*values.list)[idx]
			break
		}
	}
	if cmdValue == nil {
		newCmdValue := append(*values.list, CmdValue{Cmd: cmd})
		values.list = &newCmdValue
		cmdValue = &((*values.list)[len(*values.list)-1])
	}

	cmdValue.Value = value

	return
}

type cmdCheckFailed struct {
	message string
}

func (err cmdCheckFailed) Error() string {
	return err.message
}

// Check validates all the recorded commands
func (values *CmdValues) Check(config *Config) (err error) {
	if len(*values.list) == 0 {
		return
	}
	for idx := range *values.list {
		err = (*values.list)[idx].Check(config)
		if err != nil {
			return
		}
	}
	return
}

// CheckOne compares notes between the given command and the recorded output
func (values *CmdValues) CheckOne(cmd string, config *Config) (err error) {
	for idx := range *values.list {
		if value := (*values.list)[idx]; value.Cmd == cmd {
			err = value.Check(config)
			return
		}
	}
	return cmdCheckFailed{fmt.Sprintf("Command %q is unknown", cmd)}
}

// Check verifies that the command hasn’t changed
func (value CmdValue) Check(config *Config) (err error) {
	output, err := runCmd(value.Cmd, config)

	switch {
	case err != nil:
		logDebug("Cmd Check: %s: ERR: %v", value.Cmd, err)
		return err
	case output != value.Value:
		logDebug("Check: %s: stale (value: %v, previousvalue: %v)",
			value.Cmd, output, value.Value)
		return cmdCheckFailed{fmt.Sprintf("Command %q has changed", value.Cmd)}
	}
	logDebug("Check: %s: up to date", value.Cmd)
	return nil
}

// Formatted shows the command in a user-friendly format.
func (value *CmdValue) Formatted() string {
	return fmt.Sprintf("%q - %s", value.Cmd, value.Value)
}

// Marshal dumps the commands into gzenv format
func (values *CmdValues) Marshal() string {
	return gzenv.Marshal(*values.list)
}

// Unmarshal loads the watches back from gzenv
func (values *CmdValues) Unmarshal(from string) error {
	return gzenv.Unmarshal(from, values.list)
}

func runCmd(cmd string, config *Config) (string, error) {
	output, err := exec.Command(config.BashPath, "-c", cmd).Output()
	return string(output), err;
}
