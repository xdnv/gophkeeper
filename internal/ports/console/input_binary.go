package console

import (
	"fmt"

	"github.com/aerogu/tvchooser"
	"github.com/rivo/tview"
)

func newBinaryDataForm(app *ConsoleApp) *tview.Form {
	form := tview.NewForm()

	// entry fields
	form.AddInputField("Name", "", 30, nil, nil)

	form.AddTextArea("Path", "", 0, 4, 0, nil).
		AddButton("Select file", func() {
			path := tvchooser.FileChooser(app.Application, false)
			if path != "" {
				form.GetFormItemByLabel("Path").(*tview.TextArea).SetText(path, false)
			}
		})

	form.AddTextArea("Description", "", 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Path").(*tview.TextArea).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			message := fmt.Sprintf("Name: %s\nPath: %s\nDescription: %s\n", name, path, description)
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
