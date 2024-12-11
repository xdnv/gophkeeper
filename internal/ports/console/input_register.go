package console

import (
	"fmt"
	"internal/domain"
	"internal/transport/http_client"

	"github.com/rivo/tview"
)

func newRegistrationForm(ca *ConsoleApp) *tview.Form {
	form := tview.NewForm()
	form.
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddInputField("Email", "", 50, nil, nil).
		AddButton("Register", func() {
			username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			email := form.GetFormItemByLabel("Email").(*tview.InputField).GetText()

			uac := new(domain.UserAccountData)
			uac.Login = username
			uac.Password = password
			uac.Email = email
			err := http_client.Register(uac)

			if err != nil {
				modal := tview.NewModal().
					SetText(fmt.Sprintf("Failed to register:\n%s", err.Error())).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						ca.SetRoot(form, true) // Return to registration form
					})
				if err := ca.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
					panic(err)
				}
			} else {
				modal := tview.NewModal().
					SetText("Registration successful").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						ca.ActivateLoginPage(true)
					})
				if err := ca.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
					panic(err)
				}
			}
		}).
		AddButton("Return", func() {
			ca.ActivateLoginPage(false)
			//ca.SetRoot(createLoginForm(ca), true) // Return to auth form
		}).
		SetTitle("Registration").
		SetBorder(true)

	return form
}
