package console

import (
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"

	"github.com/aerogu/tvchooser"
	"github.com/rivo/tview"
)

func newBinaryDataForm(ca *ConsoleApp) *tview.Form {
	form := tview.NewForm()

	// entry fields
	form.AddInputField("Name", "", 30, nil, nil)

	form.AddTextArea("Path", "", 0, 4, 0, nil).
		AddButton("Select file", func() {
			path := tvchooser.FileChooser(ca.Application, false)
			if path != "" {
				form.GetFormItemByLabel("Path").(*tview.TextArea).SetText(path, false)
			}
		})

	form.AddTextArea("Description", "", 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Path").(*tview.TextArea).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			r := new(domain.KeeperRecord)
			r.Name = name
			r.Description = description
			r.SecretType = "binary"
			r.IsDeleted = false

			errMsg := "New Binary error: %s"

			k, err := domain.NewBinarySecret(path)
			if err != nil {
				ca.AppendConsole(fmt.Sprintf(errMsg, err))
				return
			}

			jsonDataCr, err := json.Marshal(k)
			if err != nil {
				ca.AppendConsole(fmt.Sprintf(errMsg, err))
				return
			}
			r.Secret = string(jsonDataCr)

			jsonData, err := json.Marshal(r)
			if err != nil {
				ca.AppendConsole(fmt.Sprintf(errMsg, err))
				return
			}

			args := []string{r.SecretType}
			resp, err := http_client.ExecuteCommand(domain.S_CMD_NEW, args, &jsonData)
			if err != nil {
				ca.AppendConsole(fmt.Sprintf(errMsg, err))
				return
			}

			message := resp.Status
			//message := fmt.Sprintf("Name: %s\nPath: %s\nDescription: %s\n", name, path, description)
			modal := tview.NewModal().
				SetText(message).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					ca.AppendConsole(message)
					ca.ActivateMainPage()
				})
			if err := ca.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
				panic(err)
			}
		})

	form.AddButton("Return", func() {
		ca.AppendConsole("Cancelled")
		ca.ActivateMainPage()
	})

	form.SetBorder(true).SetTitle("New text data")
	return form
}
