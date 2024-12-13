package console

import (
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"

	"github.com/aerogu/tvchooser"
	"github.com/rivo/tview"
)

func newBinaryDataForm(ca *ConsoleApp, r *domain.KeeperRecord) *tview.Form {
	form := tview.NewForm()

	var k = new(domain.KeeperBinary)
	var title = "Binary (new)"

	//whether we reload file
	var dataChanged = true

	var newRecord bool = (r == nil)
	if newRecord {
		r = new(domain.KeeperRecord)
		r.SecretType = domain.SECRET_BINARY
		r.IsDeleted = false
	} else {
		// if we can't read Secret, we use empty structure
		_ = json.Unmarshal([]byte(r.Secret), &k)
		title = "Binary " + r.Reference()
		dataChanged = false
	}

	// entry fields
	form.AddInputField("Name", r.Name, 30, nil, nil)

	form.AddTextArea("Path", k.FileName, 0, 4, 0, nil).
		AddButton("Select file", func() {
			path := tvchooser.FileChooser(ca.Application, false)
			if path != "" {
				form.GetFormItemByLabel("Path").(*tview.TextArea).SetText(path, false)
				dataChanged = true
			}
		})

	form.AddTextArea("Description", r.Description, 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			r.Name = form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Path").(*tview.TextArea).GetText()
			r.Description = form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			errMsg := "error: %s"

			//check for file change, reload if there's new selected
			if dataChanged {
				var err error
				k, err = domain.NewBinarySecret(path)
				if err != nil {
					ca.AppendConsole(fmt.Sprintf(errMsg, err))
					return
				}
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
			resp, err := http_client.ExecuteCommand(domain.S_CMD_UPDATE, args, &jsonData)
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

	form.SetBorder(true).SetTitle(title)
	return form
}
