package console

import (
	"fmt"

	"github.com/rivo/tview"
)

func newLoginForm(app *ConsoleApp) *tview.Flex {
	loginForm := tview.NewForm()

	header := fmt.Sprintf("the GophKeeper v.%s (%s, %s)\n\nHello user! Please name youself.", app.version.Version, app.version.Date, app.version.Commit)

	centeredLoginForm := tview.NewFlex().
		SetDirection(tview.FlexRow).
		//AddItem(tview.NewTextView().SetText(" ").SetTextAlign(tview.AlignCenter), 0, 1, false).
		AddItem(CenterVertically(tview.NewTextView().SetText(header).SetTextAlign(tview.AlignCenter)), 0, 1, false).
		AddItem(loginForm, 0, 1, true).
		AddItem(tview.NewTextView().SetText(" ").SetTextAlign(tview.AlignCenter), 0, 1, false)

	loginForm.
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Enter", func() {
			username := loginForm.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			pwd := loginForm.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			if password, exists := registeredUsers[username]; exists && password == pwd {
				app.ActivateMainPage()
			} else {
				modal := tview.NewModal().
					SetText("Wrong username or password").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						{
							loginForm.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
							app.SetRoot(centeredLoginForm, true) // Return to auth form
						}
					})
				if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
					panic(err)
				}
			}
		}).
		AddButton("Register", func() {
			app.SetRoot(newRegistrationForm(app), true) // Go to registration form
		}).
		AddButton("Exit", func() {
			app.Stop()
		}).
		SetTitle("Login").
		SetBorder(true)

	return centeredLoginForm
}
