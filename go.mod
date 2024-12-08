module github.com/xdnv/gophkeeper

go 1.22.7

require (
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

require go.uber.org/multierr v1.10.0 // indirect

require internal/adapters/cryptor v1.0.0

replace internal/adapters/cryptor => ./internal/adapters/cryptor

require internal/adapters/logger v1.0.0

replace internal/adapters/logger => ./internal/adapters/logger

replace internal/adapters/retrier => ./internal/adapters/retrier

require internal/app v1.0.0

replace internal/app => ./internal/app

require internal/domain v1.0.0

replace internal/domain => ./internal/domain

require (
	internal/ports/console v1.0.0
	internal/ports/storage v1.0.0
	internal/transport/http_server v1.0.0
)

replace internal/transport/http_server => ./internal/transport/http_server

replace internal/ports/storage => ./internal/ports/storage

replace internal/ports/console => ./internal/ports/console

require (
	github.com/aerogu/tvchooser v1.1.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.7.1 // indirect
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.1 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/tview v0.0.0-20241103174730-c76f7879f592 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
)
