package console

import (
	"context"
	"fmt"
	"strings"
)

// Command interface
type Command interface {
	Execute(ctx context.Context, args []string) (string, error)
	Help() string
}

// registered Command storage
type CommandParser struct {
	app      *ConsoleApp
	commands map[string]Command
}

// alias of built-in "help" command
const CMD_HELP = "help"

const TYPE_CREDENTIALS = "credentials"
const TYPE_CREDITCARD = "creditcard"
const TYPE_TEXT = "text"
const TYPE_BINARY = "binary"

//var datatypes = []string{TYPE_LOGIN, TYPE_CREDITCARD, TYPE_TEXT, TYPE_BINARY}

// Constructor
func NewCommandParser(ca *ConsoleApp) *CommandParser {
	return &CommandParser{
		app:      ca,
		commands: make(map[string]Command),
	}
}

// Registers new command
func (cp *CommandParser) RegisterCommand(name string, cmd Command) {
	// CMD_HELP is protected
	if name == CMD_HELP {
		return
	}
	cp.commands[name] = cmd
}

// Parses command with arguments, returns string result. Has built-in "help" command.
func (cp *CommandParser) Parse(ctx context.Context, input string) (string, error) {
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return "", nil
	}

	command := strings.ToLower(tokens[0])
	args := tokens[1:]

	if command == CMD_HELP {
		return cp.PrintHelpEntries()
	}

	if cmd, exists := cp.commands[command]; exists {
		return cmd.Execute(ctx, args)
	}

	//unknown command is not an error
	return fmt.Sprintf("Unknown command: %s. Type '%s' for help.", command, CMD_HELP), nil
}

// Executes command directly
func (cp *CommandParser) ExecuteCommand(ctx context.Context, command string, args []string) (string, error) {
	if cmd, exists := cp.commands[command]; exists {
		return cmd.Execute(ctx, args)
	}
	return "", fmt.Errorf("command not found: %s", command)
}

// Returns help entries from all registered commands as slice of strings
func (cp *CommandParser) GetHelpEntries() []string {
	var entries []string
	entries = append(entries, "- help: show available commands and their options")

	for name, cmd := range cp.commands {
		entries = append(entries, fmt.Sprintf("- %s: %s", name, cmd.Help()))
	}
	return entries
}

// Prints to string help entries from all registered commands
func (cp *CommandParser) PrintHelpEntries() (string, error) {
	entries := cp.GetHelpEntries()
	result := "Available commands:\n" + strings.Join(entries, "\n")
	return result, nil
}
