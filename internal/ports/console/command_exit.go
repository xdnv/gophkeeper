package console

import (
	"context"
	"fmt"
)

type CommandExit struct{}

func (ec *CommandExit) Execute(ctx context.Context, args []string) (string, error) {
	app, ok := ctx.Value(appCtx).(*ConsoleApp)
	if !ok {
		return "", fmt.Errorf("failed to get main application object")
	}
	app.Stop()
	return "", nil
}

func (ec *CommandExit) Help() string {
	return "exits porgram"
}
