package console

import (
	"context"
	"fmt"
	"internal/domain"
	"strings"
)

// Command prototype
type CommandNew struct{}

func (ec *CommandNew) Execute(ctx context.Context, args []string) (string, error) {
	ca, ok := ctx.Value(domain.CtxApp).(*ConsoleApp)
	if !ok {
		return "", fmt.Errorf("failed to get main application object")
	}

	if len(args) == 0 {
		return "Error: no argument given", nil
	}

	datatype := strings.ToLower(args[0])

	result := ""

	//TODO: add suspend
	// app.Suspend(func() {
	// 	if err := app.Run(); err != nil {
	// 		panic(err)
	// 	}
	// })

	switch datatype {
	case TYPE_CREDENTIALS:
		ca.ActivateNewCredentialsPage(nil)
	case TYPE_CREDITCARD:
		ca.ActivateNewCreditCardPage(nil)
	case TYPE_TEXT:
		ca.ActivateNewTextDataPage(nil)
	case TYPE_BINARY:
		ca.ActivateNewBinaryDataPage(nil)
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
