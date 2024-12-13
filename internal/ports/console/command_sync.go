package console

import (
	"context"
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"
)

// Command prototype
type CommandSync struct{}

func (ec *CommandSync) Execute(ctx context.Context, args []string) (string, error) {
	ca, ok := ctx.Value(domain.CtxApp).(*ConsoleApp)
	if !ok {
		return "", fmt.Errorf("failed to get main application object")
	}

	resp, err := http_client.ExecuteCommand(domain.S_CMD_SYNC, args, nil)
	if err != nil {
		return "", err
	}

	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return "", err
	}
	ca.UpdateRecordList()

	result := fmt.Sprintf("Sync data: %s", resp.Status)

	return result, nil
}

func (ec *CommandSync) Help() string {
	return "synchronize data records with server"
}
