package console

import (
	"fmt"
	"regexp"

	"github.com/rivo/tview"
)

func newCreditCardForm(app *ConsoleApp) *tview.Form {
	form := tview.NewForm()

	// entry fields
	form.AddInputField("Card Number", "", 20,
		func(textToCheck string, lastChar rune) bool {
			// check for digits and spaces
			if matched, _ := regexp.MatchString(`^(?:[0-9]{0,4} ?){0,4}$`, textToCheck); !matched {
				return false
			}
			return true
		},
		nil)

	form.AddInputField("Expiration Date (MM/YY)", "", 6,
		func(textToCheck string, lastChar rune) bool {
			if matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])?\/?[0-9]{0,2}$`, textToCheck); !matched {
				return false
			}
			return true
		},
		nil)

	form.AddInputField("CVV", "", 5,
		func(textToCheck string, lastChar rune) bool {
			if matched, _ := regexp.MatchString(`^[0-9]{0,3}$`, textToCheck); !matched {
				return false
			}
			return true
		},
		nil)

	form.AddTextArea("Description", "", 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			cardNumber := form.GetFormItemByLabel("Card Number").(*tview.InputField).GetText()
			expirationDate := form.GetFormItemByLabel("Expiration Date (MM/YY)").(*tview.InputField).GetText()
			cvv := form.GetFormItemByLabel("CVV").(*tview.InputField).GetText()
			description := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			//if validateCard(cardNumber) && validateExpiration(expirationDate) && validateCVV(cvv) {
			message := fmt.Sprintf("Card Number: %s\nExpiration Date: %s\nCVV: %s\nDescription: %s\n", cardNumber, expirationDate, cvv, description)
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

	form.SetBorder(true).SetTitle("New credit card data")
	return form
}

// func createCreditCardForm_ex(app *ConsoleApp) *tview.Form {

// 	form := tview.NewForm().
// 		AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
// 		AddInputField("First name", "", 20, nil, nil).
// 		AddInputField("Last name", "", 20, nil, nil).
// 		AddTextArea("Address", "", 40, 0, 0, nil).
// 		AddTextView("Notes", "This is just a demo.\nYou can enter whatever you wish.", 40, 2, true, false).
// 		AddCheckbox("Age 18+", false, nil).
// 		AddPasswordField("Password", "", 10, '*', nil).
// 		AddButton("Save", nil).
// 		AddButton("Quit", func() {
// 			app.Stop()
// 		})

// 	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

// 	// if err := app.SetRoot(form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
// 	// 	panic(err)
// 	// }

// 	return form
// }
