package console

import (
	"context"
	"fmt"
	"strings"
)

// ExampleCommand - пример реализации команды
type CommandNew struct{}

func (ec *CommandNew) Execute(ctx context.Context, args []string) (string, error) {
	app, ok := ctx.Value(appCtx).(*ConsoleApp)
	if !ok {
		return "", fmt.Errorf("failed to get main application object")
	}

	if len(args) == 0 {
		return "Error: no argument given", nil
	}

	datatype := strings.ToLower(args[0])

	result := ""

	switch datatype {
	case TYPE_CREDENTIALS:
		app.ActivateNewCredentialsPage()
	case TYPE_CREDITCARD:
		app.ActivateNewCreditCardPage()
	case TYPE_TEXT:
		app.ActivateNewTextDataPage()
	case TYPE_BINARY:
		app.ActivateNewBinaryDataPage()
	default:
		result := fmt.Sprintf("Error: unknown argument: %s", datatype)
		return result, nil
	}

	//result := fmt.Sprintf("Executing New with arguments: %s, len: %v", strings.Join(args, "->"), len(args))
	return result, nil
}

func (ec *CommandNew) Help() string {
	return "create new object of type: [credentials][creditcard][text][binary]"
}
