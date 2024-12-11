package console

import (
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"

	"github.com/rivo/tview"
)

func newTextDataForm(ca *ConsoleApp) *tview.Form {
	form := tview.NewForm()

	// entry fields
	form.AddInputField("Name", "", 30, nil, nil)

	form.AddTextArea("Text", "", 0, 15, 0, nil)

	form.AddTextArea("Description", "", 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			text := form.GetFormItemByLabel("Text").(*tview.TextArea).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			r := new(domain.KeeperRecord)
			r.Name = name
			r.Description = description
			r.SecretType = "text"
			r.IsDeleted = false

			errMsg := "New Text error: %s"

			k := new(domain.KeeperText)
			k.Text = text

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
			resp, err := http_client.ExecuteCommand("new", args, &jsonData)
			if err != nil {
				ca.AppendConsole(fmt.Sprintf(errMsg, err))
				return
			}

			message := resp.Status
			//message := fmt.Sprintf("Name: %s\nText: %s\nDescription: %s\n", name, text, description)
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
