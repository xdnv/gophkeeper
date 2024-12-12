package console

import (
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"

	"github.com/rivo/tview"
)

func newCredentialsForm(ca *ConsoleApp, r *domain.KeeperRecord) *tview.Form {
	form := tview.NewForm()

	var k = domain.SecretCredentials{}
	var title = "Credentials (new)"

	var newRecord bool = (r == nil)
	if newRecord {
		r = new(domain.KeeperRecord)
		r.SecretType = domain.SECRET_CREDENTIALS
		r.IsDeleted = false
	} else {
		// if we can't read Secret, we use empty structure
		_ = json.Unmarshal([]byte(r.Secret), &k)
		title = "Credentials " + r.Reference()
	}

	// entry fields
	form.AddInputField("Name", r.Name, 30, nil, nil)

	form.AddInputField("Address", k.Address, 50, nil, nil).
		AddInputField("Login", k.Login, 50, nil, nil).
		AddInputField("Password", k.Password, 50, nil, nil)

	form.AddTextArea("Description", r.Description, 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			r.Name = form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			k.Address = form.GetFormItemByLabel("Address").(*tview.InputField).GetText()
			k.Login = form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
			k.Password = form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			r.Description = form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			errMsg := "error: %s"

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
			//message := fmt.Sprintf("Name: %s\nAddress: %s\nLogin: %s\nPassword: %s\nDescription: %s\n", name, address, login, password, description)
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
