package console

import (
	"context"
	"fmt"
	"internal/app"
	"internal/transport/http_client"
)

// Command prototype
type CommandPing struct{}

func (ec *CommandPing) Execute(ctx context.Context, args []string) (string, error) {
	// ca, ok := ctx.Value(domain.CtxApp).(*ConsoleApp)
	// if !ok {
	// 	return "", fmt.Errorf("failed to get main application object")
	// }

	resp, err := http_client.Ping()

	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Ping server %s: %s", app.Cc.Endpoint, resp.Status)

	return result, nil
}

func (ec *CommandPing) Help() string {
	return "ping established connection to server"
}
