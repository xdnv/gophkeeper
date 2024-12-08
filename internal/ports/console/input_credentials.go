package console

import (
	"fmt"

	"github.com/rivo/tview"
)

func newCredentialsForm(app *ConsoleApp) *tview.Form {
	form := tview.NewForm()

	// entry fields
	form.AddInputField("Address", "", 50, nil, nil).
		AddInputField("Login", "", 50, nil, nil).
		AddInputField("Password", "", 50, nil, nil)

	form.AddTextArea("Description", "", 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			address := form.GetFormItemByLabel("Address").(*tview.InputField).GetText()
			login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			message := fmt.Sprintf("Address: %s\nLogin: %s\nPassword: %s\nDescription: %s\n", address, login, password, description)
			modal := tview.NewModal().
				SetText(message).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.AppendConsole(message)
					app.ActivateMainPage()
				})
			if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
				panic(err)
			}
		})

	form.AddButton("Return", func() {
		app.AppendConsole("Cancelled")
		app.ActivateMainPage()
	})

	form.SetBorder(true).SetTitle("New credentials")
	return form
}
