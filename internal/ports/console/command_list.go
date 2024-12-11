package console

import (
	"context"
	"fmt"
	"internal/domain"
)

// Command prototype
type CommandList struct{}

func (ec *CommandList) Execute(ctx context.Context, args []string) (string, error) {
	ca, ok := ctx.Value(domain.CtxApp).(*ConsoleApp)
	if !ok {
		return "", fmt.Errorf("failed to get main application object")
	}

	result := "Available records:\n"
	for i := 0; i < ca.list.GetItemCount(); i++ {
		text, _ := ca.list.GetItemText(i)
		result += fmt.Sprintf("- %s\n", text)
	}

	return result, nil
}

func (ec *CommandList) Help() string {
	return "list available entries. Optionally can filter by type: [credentials][creditcard][text][binary] and text substring"
}
