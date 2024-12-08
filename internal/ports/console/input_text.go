package console

import (
	"fmt"

	"github.com/rivo/tview"
)

func createTextDataForm(app *ConsoleApp) *tview.Form {
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

			message := fmt.Sprintf("Name: %s\nText: %s\nDescription: %s\n", name, text, description)
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

	form.SetBorder(true).SetTitle("New text data")
	return form
}
