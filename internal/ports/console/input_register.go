package console

import (
	"github.com/rivo/tview"
)

func newRegistrationForm(app *ConsoleApp) *tview.Form {
	form := tview.NewForm()
	form.
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Register", func() {
			username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			if _, exists := registeredUsers[username]; exists {
				modal := tview.NewModal().
					SetText("Sorry, this username already exists.\nPlease choose another or use login screen.").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetRoot(form, true) // Return to registration form
					})
				if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
					panic(err)
				}
			} else {
				registeredUsers[username] = password // Store new user & password
				modal := tview.NewModal().
					SetText("Registration successful").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.ActivateLoginPage(true)
						//app.SetRoot(createLoginForm(app), true) // Return to auth form
					})
				if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
					panic(err)
				}
			}
		}).
		AddButton("Return", func() {
			app.ActivateLoginPage(false)
			//app.SetRoot(createLoginForm(app), true) // Return to auth form
		}).
		SetTitle("Registration").
		SetBorder(true)

	return form
}
