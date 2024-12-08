module console

go 1.22.7

require (
	github.com/gdamore/tcell/v2 v2.7.1
	github.com/rivo/tview v0.0.0-20241103174730-c76f7879f592
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

require (
	github.com/aerogu/tvchooser v1.1.0
	internal/domain v1.0.0
)

replace internal/domain => ../../domain
