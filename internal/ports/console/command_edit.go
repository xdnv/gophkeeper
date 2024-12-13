package console

import (
	"context"
	"fmt"
	"internal/domain"
)

// Command prototype
type CommandEdit struct{}

func (ec *CommandEdit) Execute(ctx context.Context, args []string) (string, error) {
	ca, ok := ctx.Value(domain.CtxApp).(*ConsoleApp)
	if !ok {
		return "", fmt.Errorf("failed to get main application object")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("please provide record ID")
	}

	searchID := args[0]
	item, err := SearchByID(searchID)
	if err != nil {
		return "", err
	}

	recordName := item.Reference()

	// sending full ID
	//serverArgs := []string{item.ID}

	//TODO: add suspend
	// app.Suspend(func() {
	// 	if err := app.Run(); err != nil {
	// 		panic(err)
	// 	}
	// })

	switch item.SecretType {
	case TYPE_CREDENTIALS:
		ca.ActivateNewCredentialsPage(item)
	case TYPE_CREDITCARD:
		ca.ActivateNewCreditCardPage(item)
	case TYPE_TEXT:
		ca.ActivateNewTextDataPage(item)
	case TYPE_BINARY:
		ca.ActivateNewBinaryDataPage(item)
	default:
		return "", fmt.Errorf("unknown argument: %s", item.SecretType)
	}

	result := fmt.Sprintf("Succesfully updated [%s]", recordName)
	return result, nil
}

func (ec *CommandEdit) Help() string {
	return "edit secret entry. syntax: edit <listID|short#|long#>"
}
