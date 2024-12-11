package console

import (
	"context"
	"fmt"
	"internal/domain"
)

// Command prototype
type CommandExit struct{}

func (ec *CommandExit) Execute(ctx context.Context, args []string) (string, error) {
	ca, ok := ctx.Value(domain.CtxApp).(*ConsoleApp)
	if !ok {
		return "", fmt.Errorf("failed to get main application object")
	}
	ca.Stop()
	return "", nil
}

func (ec *CommandExit) Help() string {
	return "exits porgram"
}
