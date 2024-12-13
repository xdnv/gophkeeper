package console

import (
	"context"
	"fmt"
	"strings"

	"internal/domain"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var items []domain.KeeperRecord

var cp *CommandProcessor

type ConsoleApp struct {
	*tview.Application
	ctx       context.Context
	list      *tview.List
	content   *tview.TextView
	console   *tview.TextView
	input     *tview.InputField
	flex      *tview.Flex
	inConsole bool
	version   *domain.Version
	loginForm *tview.Flex
	mainForm  *tview.Flex
}

func CenterVertically(text *tview.TextView) *tview.TextView {
	text.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		y += h / 2
		return x, y, w, h
	})
	return text
}

func NewApp(ctx context.Context) *ConsoleApp {
	return &ConsoleApp{
		Application: tview.NewApplication(),
		ctx:         ctx,
		version:     domain.GetVersion(),
	}
}

func (ca *ConsoleApp) InitCommands() {
	cp = NewCommandParser(ca)
	cp.RegisterCommand("sync", &CommandSync{})
	cp.RegisterCommand("list", &CommandList{})
	cp.RegisterCommand("new", &CommandNew{})
	cp.RegisterCommand("edit", &CommandEdit{})
	cp.RegisterCommand("dump", &CommandDump{})
	cp.RegisterCommand("delete", &CommandDelete{})
	cp.RegisterCommand("ping", &CommandPing{})
	cp.RegisterCommand("exit", &CommandExit{})
}

func (ca *ConsoleApp) Init() {
	ca.InitCommands()
	ca.ActivateLoginPage(false)
	// app.loginForm = createLoginForm(app)
	// app.SetRoot(app.loginForm, true) // Set Login Form as main form
}

// Update list from server
func SyncRecordList(ca *ConsoleApp) {
	ca.HandleCommand("sync")
	ca.UpdateRecordList()
}

// Search list record using ListNR, ShortID and ID
func SearchByID(id string) (*domain.KeeperRecord, error) {

	//ListNR
	if strings.HasPrefix(id, "#") {
		rowNum := strings.TrimPrefix(id, "#")
		for _, record := range items {
			if record.ListNR == rowNum {
				return &record, nil
			}
		}
	}

	//ShortID
	if len(id) == 8 {
		for _, record := range items {
			if record.ShortID == id {
				return &record, nil
			}
		}
	}

	//ID
	if len(id) == 36 {
		for _, record := range items {
			if record.ID == id {
				return &record, nil
			}
		}
	}

	return nil, fmt.Errorf("no record found with identifier %s", id)
}

func (ca *ConsoleApp) ActivateLoginPage(clear bool) {
	if ca.loginForm == nil {
		ca.loginForm = newLoginForm(ca)
	}
	if clear {
		form := ca.loginForm.GetItem(1).(*tview.Form)
		form.GetFormItemByLabel("Username").(*tview.InputField).SetText("")
		form.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
	}
	ca.SetRoot(ca.loginForm, true).SetFocus(ca.loginForm).Run()
}

func (ca *ConsoleApp) ActivateMainPage() {
	if ca.mainForm == nil {
		ca.mainForm = newMainLayout(ca)
	}
	ca.ResetFocus()
	enableCapture(ca)
	//app.SetRoot(app.mainForm, true).SetFocus(app.mainForm).Run()
	ca.SetRoot(ca.mainForm, true).SetFocus(ca.list)
	SyncRecordList(ca)
	ca.Run()
}

func (ca *ConsoleApp) ActivateNewCreditCardPage(r *domain.KeeperRecord) {
	form := newCreditCardForm(ca, r)
	disableCapture(ca)
	ca.SetRoot(form, true).SetFocus(form).Run()
}

// creates form for a new or editable "credentials" object
func (ca *ConsoleApp) ActivateNewCredentialsPage(r *domain.KeeperRecord) {
	form := newCredentialsForm(ca, r)
	disableCapture(ca)
	ca.SetRoot(form, true).SetFocus(form).Run()
}

func (ca *ConsoleApp) ActivateNewTextDataPage(r *domain.KeeperRecord) {
	form := newTextDataForm(ca, r)
	disableCapture(ca)
	ca.SetRoot(form, true).SetFocus(form).Run()
}

func (ca *ConsoleApp) ActivateNewBinaryDataPage(r *domain.KeeperRecord) {
	form := newBinaryDataForm(ca, r)
	disableCapture(ca)
	ca.SetRoot(form, true).SetFocus(form).Run()
}

func (app *ConsoleApp) ClearConsole() {
	app.console.SetText("")
}

func (ca *ConsoleApp) AppendConsole(message string) {
	ca.console.SetText(fmt.Sprintf("%s\n%s", ca.console.GetText(false), message))
}

func (ca *ConsoleApp) UpdateRecordList() {
	ca.list.Clear()
	for _, item := range items {
		ca.list.AddItem(fmt.Sprintf("#%s. %s", item.ListNR, item.Name), item.Description, 0, nil)
	}
}

func (ca *ConsoleApp) HandleCommand(cmdLine string) {
	if len(cmdLine) == 0 {
		return
	}
	ctx := context.WithValue(ca.ctx, domain.CtxApp, ca)
	ca.AppendConsole(fmt.Sprintf(">%s", cmdLine))
	result, err := cp.Parse(ctx, cmdLine)
	if err != nil {
		ca.AppendConsole(fmt.Sprintf("Error executing command: %s\n", err.Error()))
		return
	}
	ca.AppendConsole(result)
}

// Switch between console and UI when ~ is pressed
func (ca *ConsoleApp) ToggleConsoleFocus() {
	if ca.inConsole {
		ca.flex.ResizeItem(ca.flex.GetItem(0), 0, 4)
		ca.flex.ResizeItem(ca.flex.GetItem(1), 0, 1)
		ca.Application.SetFocus(ca.list)
	} else {
		ca.flex.ResizeItem(ca.flex.GetItem(0), 0, 1)
		ca.flex.ResizeItem(ca.flex.GetItem(1), 0, 4)
		ca.Application.SetFocus(ca.input)
	}
	ca.inConsole = !ca.inConsole
}

// Switch between console and UI when ~ is pressed
func (ca *ConsoleApp) ResetFocus() {
	ca.inConsole = true
	ca.ToggleConsoleFocus()
}

// // Keypress handler to switch between windows
// func (app *ConsoleApp) HandleKeyPress(event *tcell.KeyEvent) *tview.Application {
// 	if event.Key == tcell.KeyRunes && event.Rune == '~' { // Проверяем нажатие тильды (~)
// 		app.ToggleFocus()
// 	}
// 	return app.Application
// }
