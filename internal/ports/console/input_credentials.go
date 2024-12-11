package console

import (
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"

	"github.com/rivo/tview"
)

func newCredentialsForm(ca *ConsoleApp) *tview.Form {
	form := tview.NewForm()

	// entry fields
	form.AddInputField("Name", "", 30, nil, nil)

	form.AddInputField("Address", "", 50, nil, nil).
		AddInputField("Login", "", 50, nil, nil).
		AddInputField("Password", "", 50, nil, nil)

	form.AddTextArea("Description", "", 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			address := form.GetFormItemByLabel("Address").(*tview.InputField).GetText()
			login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			r := new(domain.KeeperRecord)
			r.Name = name
			r.Description = description
			r.SecretType = "credentials"
			r.IsDeleted = false

			errMsg := "New Credentials error: %s"

			k := new(domain.KeeperCredentials)
			k.Address = address
			k.Login = login
			k.Password = password

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

	form.SetBorder(true).SetTitle("New credentials")
	return form
}
