package console

import (
	"fmt"
	"internal/app"
	"internal/domain"
	"internal/transport/http_client"

	"github.com/rivo/tview"
)

func newLoginForm(ca *ConsoleApp) *tview.Flex {
	loginForm := tview.NewForm()

	header := fmt.Sprintf("the GophKeeper v.%s (%s, %s)\n\nHello user! Please name youself.", ca.version.Version, ca.version.Date, ca.version.Commit)

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
			password := loginForm.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			uac := new(domain.UserAccountData)
			uac.Login = username
			uac.Password = password
			ar, err := http_client.Login(uac)

			if err != nil {
				modal := tview.NewModal().
					SetText(fmt.Sprintf("Failed to login:\n%s", err.Error())).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						{
							loginForm.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
							ca.SetRoot(centeredLoginForm, true) // Return to auth form
						}
					})
				if err := ca.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
					panic(err)
				}
			} else {
				app.Cc.AuthToken = ar.Token
				app.Cc.SessionSigningKey = &ar.PublicKey
				ca.ActivateMainPage()
			}
		}).
		AddButton("Register", func() {
			ca.SetRoot(newRegistrationForm(ca), true) // Go to registration form
		}).
		AddButton("Exit", func() {
			ca.Stop()
		}).
		SetTitle("Login").
		SetBorder(true)

	return centeredLoginForm
}
