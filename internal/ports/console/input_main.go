package console

import (
	"fmt"
	"internal/domain"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newMainLayout(ca *ConsoleApp) *tview.Flex {
	ca.list = tview.NewList()
	ca.list.
		SetTitle("Records").
		SetBorder(true)

	ca.content = tview.NewTextView()
	ca.content.
		SetText("select list record to view its contents.").
		SetScrollable(true).
		SetMaxLines(150).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetTitle("Data").
		SetBorder(true)

	ca.list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := items[index]
		ca.content.SetText(fmt.Sprintf("#%s. %s\n%s\n%s [%s]\n%s", item.ListNR, item.Name, item.Description, item.ShortID, item.ID, "(press ENTER to decrypt data)"))
	})

	ca.list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := items[index]

		decrypted := ""
		secret, err := domain.KeepDeserialized(item.SecretType, []byte(item.Secret))
		if err == nil {
			decrypted = domain.KeepReadable(*secret)
		} else {
			decrypted = fmt.Sprintf("error decrypting data: %s", err)
		}

		ca.content.SetText(fmt.Sprintf("#%s. %s\n%s\n%s [%s]\n%s", item.ListNR, item.Name, item.Description, item.ShortID, item.ID, "DECRYPTED DATA:\n"+decrypted))
	})

	ca.console = tview.NewTextView()
	ca.console.
		SetText("").
		SetScrollable(true).
		SetMaxLines(150).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetTitle("Console").
		SetBorder(true)

	ca.input = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)

	ca.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			command := ca.input.GetText()
			if strings.TrimSpace(command) != "" {
				ca.input.SetText("")              // Clear input field
				ca.HandleCommand(command)         // Process command
				ca.console.ScrollToEnd()          // Scroll to bottom
				ca.Application.SetFocus(ca.input) // Return focus
			}
			return nil // Stop event processing
		case tcell.KeyEscape:
			ca.input.SetText("") // Clear input when ESC is pressed
			return nil           // Stop event processing
		default:
			return event // Give event to next chained event processor
		}
	})

	// SetChangedFunc(func(text string) {
	// 	if text == "" {
	// 		return
	// 	}
	// 	handleCommand(app, text)
	// 	app.input.SetText("")
	// 	app.console.ScrollToEnd()
	// 	app.Application.SetFocus(app.input)
	// })

	instruction := tview.NewTextView().SetText(
		"TAB: switch list/console | PgUp/Dn: scroll output | CTRL+C: exit",
	).SetTextColor(tview.Styles.SecondaryTextColor)

	listFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ca.list, 0, 1, true).
		AddItem(ca.content, 0, 3, false)

	consoleFlex := tview.NewFlex().
		AddItem(ca.console, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(listFlex, 0, 4, false).
		AddItem(consoleFlex, 0, 1, false).
		AddItem(ca.input, 1, 0, false).
		AddItem(instruction, 1, 0, false)

	ca.flex = flex

	// // Update list from server
	// SyncRecordList(ca)

	ca.Application.SetRoot(flex, true)
	enableCapture(ca)

	return flex
}

func enableCapture(ca *ConsoleApp) {
	ca.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			ca.ToggleConsoleFocus()
			return nil
		case tcell.KeyPgUp: // Scroll up
			r, _ := ca.console.GetScrollOffset()
			ca.console.ScrollTo(r-10, 0) // 10 lines up
			return nil
		case tcell.KeyPgDn: // Scroll down
			r, _ := ca.console.GetScrollOffset()
			ca.console.ScrollTo(r+10, 0) // 10 lines down
			return nil
		case tcell.KeyRune:
			if event.Rune() == '~' || event.Rune() == '`' {
				ca.ToggleConsoleFocus()
				return nil
			}
			//app.console.SetText(fmt.Sprintf("Got rune: %c\n", event.Rune()))
			return event
		default:
			return event
		}
	})
}

func disableCapture(ca *ConsoleApp) {
	ca.SetInputCapture(nil)
}
