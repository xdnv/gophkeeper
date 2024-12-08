package console

import (
	"context"
	"fmt"

	"internal/domain"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	correctUsername = "user"
	correctPassword = "password"
)

var registeredUsers = map[string]string{
	correctUsername: correctPassword,
}

type ListItem struct {
	ID      string
	Text    string
	Comment string
}

var items = []ListItem{
	{"1", "Item 1", "Description 1"},
	{"2", "Item 2", "Description 2 very long description to read only on big big screen"},
	{"3", "Item 3", "Description 3"},
	{"4", "Item 4", "Description 4"},
	{"5", "Item 5", "Description 5 www.blablablablabla.com/somenewpage/wow?abcde=12345"},
	{"6", "Item 6", "Description 6"},
	{"7", "Item 7", "Description 7"},
	{"8", "Item 8", "Description 8"},
	{"9", "Item 9", "Description 9"},
	{"10", "Item 10", "Description 10"},
	{"11", "Item 11", "Description 11"},
	{"12", "Item 12", "Description 12"},
	{"13", "Item 13", "Description 13"},
	{"14", "Item 14", "Description 14"},
	{"15", "Item 15", "Description 15"},
	{"16", "MCRD 1234", "MC 1234: This card is so important to me so its description cannot fit to even biggest screen in the room. Or, maybe, in the world too. Wow! Exceptionally long description."},
	{"17", "VISA 5678", "Card #2"},
	{"18", "VISA 9012", "Card #3"},
	{"19", "MIR 3456", "Card #4"},
	{"20", "MIR 7890", "Card #5"},
	{"21", "VISA 1234", "Card #6"},
}

var cp *CommandParser

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

type key string

const (
	appCtx key = "app"
)

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

func (app *ConsoleApp) InitCommands() {
	cp = NewCommandParser(app)
	cp.RegisterCommand("list", &CommandList{})
	cp.RegisterCommand("new", &CommandNew{})
	cp.RegisterCommand("exit", &CommandExit{})
}

func (app *ConsoleApp) Init() {
	app.InitCommands()
	app.ActivateLoginPage(false)
	// app.loginForm = createLoginForm(app)
	// app.SetRoot(app.loginForm, true) // Set Login Form as main form
}

func (app *ConsoleApp) ActivateLoginPage(clear bool) {
	if app.loginForm == nil {
		app.loginForm = newLoginForm(app)
	}
	if clear {
		form := app.loginForm.GetItem(1).(*tview.Form)
		form.GetFormItemByLabel("Username").(*tview.InputField).SetText("")
		form.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
	}
	app.SetRoot(app.loginForm, true).SetFocus(app.loginForm).Run()
}

func (app *ConsoleApp) ActivateMainPage() {
	if app.mainForm == nil {
		app.mainForm = newMainLayout(app)
	}
	app.ResetFocus()
	enableCapture(app)
	//app.SetRoot(app.mainForm, true).SetFocus(app.mainForm).Run()
	app.SetRoot(app.mainForm, true).SetFocus(app.list).Run()
}

func (app *ConsoleApp) ActivateNewCreditCardPage() {
	form := newCreditCardForm(app)
	disableCapture(app)
	app.SetRoot(form, true).SetFocus(form).Run()
}

func (app *ConsoleApp) ActivateNewCredentialsPage() {
	form := newCredentialsForm(app)
	disableCapture(app)
	app.SetRoot(form, true).SetFocus(form).Run()
}

func (app *ConsoleApp) ActivateNewTextDataPage() {
	form := newTextDataForm(app)
	disableCapture(app)
	app.SetRoot(form, true).SetFocus(form).Run()
}

func (app *ConsoleApp) ActivateNewBinaryDataPage() {
	form := newBinaryDataForm(app)
	disableCapture(app)
	app.SetRoot(form, true).SetFocus(form).Run()
}

func (app *ConsoleApp) ClearConsole() {
	app.console.SetText("")
}

func (app *ConsoleApp) AppendConsole(message string) {
	app.console.SetText(fmt.Sprintf("%s\n%s", app.console.GetText(false), message))
}

func (app *ConsoleApp) handleCommand(cmdLine string) {
	if len(cmdLine) == 0 {
		return
	}
	ctx := context.WithValue(app.ctx, appCtx, app)
	app.AppendConsole(fmt.Sprintf(">%s", cmdLine))
	result, err := cp.Parse(ctx, cmdLine)
	if err != nil {
		app.AppendConsole(fmt.Sprintf("Error executing command: %s\n", err.Error()))
		return
	}
	app.AppendConsole(result)
}

// Switch between console and UI when ~ is pressed
func (app *ConsoleApp) ToggleConsoleFocus() {
	if app.inConsole {
		app.flex.ResizeItem(app.flex.GetItem(0), 0, 4)
		app.flex.ResizeItem(app.flex.GetItem(1), 0, 1)
		app.Application.SetFocus(app.list)
	} else {
		app.flex.ResizeItem(app.flex.GetItem(0), 0, 1)
		app.flex.ResizeItem(app.flex.GetItem(1), 0, 4)
		app.Application.SetFocus(app.input)
	}
	app.inConsole = !app.inConsole
}

// Switch between console and UI when ~ is pressed
func (app *ConsoleApp) ResetFocus() {
	app.inConsole = true
	app.ToggleConsoleFocus()
}

// // Keypress handler to switch between windows
// func (app *ConsoleApp) HandleKeyPress(event *tcell.KeyEvent) *tview.Application {
// 	if event.Key == tcell.KeyRunes && event.Rune == '~' { // Проверяем нажатие тильды (~)
// 		app.ToggleFocus()
// 	}
// 	return app.Application
// }
