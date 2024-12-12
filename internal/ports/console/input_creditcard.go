package console

import (
	"encoding/json"
	"fmt"
	"internal/domain"
	"internal/transport/http_client"
	"regexp"

	"github.com/rivo/tview"
)

func newCreditCardForm(ca *ConsoleApp, r *domain.KeeperRecord) *tview.Form {
	form := tview.NewForm()

	var k = domain.KeeperCreditcard{}
	var title = "Credit card (new)"

	var newRecord bool = (r == nil)
	if newRecord {
		r = new(domain.KeeperRecord)
		r.SecretType = domain.SECRET_CREDITCARD
		r.IsDeleted = false
	} else {
		// if we can't read Secret, we use empty structure
		_ = json.Unmarshal([]byte(r.Secret), &k)
		title = "Credit card " + r.Reference()
	}

	// entry fields
	form.AddInputField("Name", r.Name, 30, nil, nil)

	form.AddInputField("Card Number", k.CardNumber, 20,
		func(textToCheck string, lastChar rune) bool {
			// check for digits and spaces
			if matched, _ := regexp.MatchString(`^(?:[0-9]{0,4} ?){0,4}$`, textToCheck); !matched {
				return false
			}
			return true
		},
		nil)

	form.AddInputField("Expiration Date (MM/YY)", k.ExpirationDate, 6,
		func(textToCheck string, lastChar rune) bool {
			if matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])?\/?[0-9]{0,2}$`, textToCheck); !matched {
				return false
			}
			return true
		},
		nil)

	form.AddInputField("CVV", k.SecurityCode, 5,
		func(textToCheck string, lastChar rune) bool {
			if matched, _ := regexp.MatchString(`^[0-9]{0,3}$`, textToCheck); !matched {
				return false
			}
			return true
		},
		nil)

	form.AddTextArea("Description", r.Description, 0, 5, 0, nil)

	form.AddButton("Submit",
		func() {
			r.Name = form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			k.CardNumber = form.GetFormItemByLabel("Card Number").(*tview.InputField).GetText()
			k.ExpirationDate = form.GetFormItemByLabel("Expiration Date (MM/YY)").(*tview.InputField).GetText()
			k.SecurityCode = form.GetFormItemByLabel("CVV").(*tview.InputField).GetText()
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
			//message := fmt.Sprintf("Name: %s\nCard Number: %s\nExpiration Date: %s\nCVV: %s\nDescription: %s\n", name, cardNumber, expirationDate, cvv, description)
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
