package console

import (
	"context"
	"encoding/json"
	"fmt"
	"internal/domain"
)

// Command prototype
type CommandDump struct{}

func (ec *CommandDump) Execute(ctx context.Context, args []string) (string, error) {
	// ca, ok := ctx.Value(domain.CtxApp).(*ConsoleApp)
	// if !ok {
	// 	return "", fmt.Errorf("failed to get main application object")
	// }

	if len(args) < 2 {
		return "", fmt.Errorf("please provide Id and path to save file")
	}

	searchID := args[0]
	filePath := args[1]
	item, err := SearchByID(searchID)
	if err != nil {
		return "", err
	}

	recordName := item.Reference() //fmt.Sprintf("#%s. %s (%s)", item.ListNR, item.Name, item.ShortID)

	if item.SecretType != domain.SECRET_BINARY {
		return "", fmt.Errorf("not a binary secret: %s", recordName)
	}

	var secret domain.KeeperBinary
	err = json.Unmarshal([]byte(item.Secret), &secret)
	if err != nil {
		return "", err
	}

	err = domain.DumpBinary(&secret, filePath)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("succesfully dumped [%s] to %s", recordName, filePath)
	return result, nil
}

func (ec *CommandDump) Help() string {
	return "dump binary secret to file. syntax: dump <listID|short#|long#> <filename>"
}
