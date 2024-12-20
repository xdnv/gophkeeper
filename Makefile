# Get version from git hash
git_hash := $(shell git rev-parse --short HEAD || echo 'development')

# Get current date
current_time = $(shell date +"%Y-%m-%d:T%H:%M:%S")

# Add linker flags
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_hash}'

# Build binaries for current OS and Linux
.PHONY:
build:
    @echo "Building binaries..."
    go build -ldflags=${linker_flags} -o=./bin/binver ./cmd/server/main.go
    GOOS=windows GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/win_amd64/binver ./cmd/server/main.go