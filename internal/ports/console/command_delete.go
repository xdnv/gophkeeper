package console

import (
	"context"
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"
)

// Command prototype
type CommandDelete struct{}

func (ec *CommandDelete) Execute(ctx context.Context, args []string) (string, error) {
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
	serverArgs := []string{item.ID}
	resp, err := http_client.ExecuteCommand(domain.S_CMD_DELETE, serverArgs, nil)
	if err != nil {
		return "", err
	}

	// sync after delete
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return "", err
	}
	ca.UpdateRecordList()

	result := fmt.Sprintf("Succesfully deleted [%s]", recordName)
	return result, nil
}

func (ec *CommandDelete) Help() string {
	return "delete secret entry. syntax: delete <listID|short#|long#>"
}
