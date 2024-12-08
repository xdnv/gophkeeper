package console

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newMainLayout(app *ConsoleApp) *tview.Flex {
	app.list = tview.NewList()

	for _, item := range items {
		app.list.AddItem(fmt.Sprintf("#%s. %s", item.ID, item.Text), item.Comment, 0, nil)
	}
	app.list.
		SetTitle("Records").
		SetBorder(true)

	app.content = tview.NewTextView()
	app.content.
		SetText("select list record to view its contents.").
		SetScrollable(true).
		SetMaxLines(150).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetTitle("Data").
		SetBorder(true)

	app.list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := items[index]
		app.content.SetText(fmt.Sprintf("#%s. %s\n%s\n%s", item.ID, item.Text, item.Comment, "(press ENTER to decrypt data)"))
		//app.console.SetText("Selected ID: " + itemID)
	})

	app.list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		itemID := items[index].ID
		app.content.SetText("DECRYPTED ID: " + itemID)
		//app.console.SetText("Selected ID: " + itemID)
	})

	app.console = tview.NewTextView()
	app.console.
		SetText("").
		SetScrollable(true).
		SetMaxLines(150).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetTitle("Console").
		SetBorder(true)

	app.input = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)

	app.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			command := app.input.GetText()
			if strings.TrimSpace(command) != "" {
				app.input.SetText("")               // Очищаем поле ввода
				app.handleCommand(command)          // Обрабатываем команду
				app.console.ScrollToEnd()           // Прокручиваем консоль вниз после ввода команды
				app.Application.SetFocus(app.input) // Возвращаем фокус на поле ввода
			}
			return nil // Не передаем событие дальше
		case tcell.KeyEscape:
			app.input.SetText("") // Очищаем поле ввода при нажатии Esc
			return nil            // Не передаем событие дальше
		default:
			return event // Передаем событие дальше для других обработчиков
		}
	})

	// SetChangedFunc(func(text string) {
	// 	if text == "" {
	// 		return
	// 	}
	// 	handleCommand(app, text)
	// 	app.input.SetText("")               // Очищаем ввод после обработки команды
	// 	app.console.ScrollToEnd()           // Прокручиваем консоль вниз после ввода команды
	// 	app.Application.SetFocus(app.input) // Возвращаем фокус на поле ввода
	// })

	instruction := tview.NewTextView().SetText(
		"TAB: switch list/console | PgUp/Dn: scroll output | CTRL+C: exit",
	).SetTextColor(tview.Styles.SecondaryTextColor)

	listFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(app.list, 0, 1, true).
		AddItem(app.content, 0, 3, false)

	consoleFlex := tview.NewFlex().
		AddItem(app.console, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(listFlex, 0, 4, false).
		AddItem(consoleFlex, 0, 1, false).
		AddItem(app.input, 1, 0, false).
		AddItem(instruction, 1, 0, false)

	app.flex = flex

	app.Application.SetRoot(flex, true)
	enableCapture(app)

	return flex
}

func enableCapture(app *ConsoleApp) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.ToggleConsoleFocus()
			return nil
		case tcell.KeyPgUp: // Scroll up
			r, _ := app.console.GetScrollOffset()
			app.console.ScrollTo(r-10, 0) // 10 lines up
			return nil
		case tcell.KeyPgDn: // Scroll down
			r, _ := app.console.GetScrollOffset()
			app.console.ScrollTo(r+10, 0) // 10 lines down
			return nil
		case tcell.KeyRune:
			if event.Rune() == '~' || event.Rune() == '`' {
				app.ToggleConsoleFocus()
				return nil
			}
			//app.console.SetText(fmt.Sprintf("Got rune: %c\n", event.Rune()))
			return event
		default:
			return event
		}
	})
}

func disableCapture(app *ConsoleApp) {
	app.SetInputCapture(nil)
}
