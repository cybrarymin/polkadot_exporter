## help: print the help message
-include .envrc # -include will include .envrc but if it doesn't exist it won't return error. .envrc usually is not commited in git so to avoid pipeline failure we do this

#================================================================#
# HELPERS
#================================================================#

# always use helo as the first target. Because make command without any target will run first target defined in it. "make" will equal to "make help"
.PHONY: help # .PHONY for each target will teach make if we have a local file or directory that names help pls don't consider them and use the target we defined cause make command can't dinstingush the directory or file from targets we define inside makefile and it get's confused
help: # @ before the command will not echo the command itself when we run make <target> command
	@echo "Usage:" 
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: prerequsite_confirm
prerequsite_confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]



#================================================================#
# DEVELOPMENT
#================================================================#
## build/exporter: building the application
current_time = $(shell date +"%Y-%m-%dT%H:%M:%S%z")
git_version = $(shell git describe --always --long --dirty --tags 2>/dev/null; if [[ $$? != 0 ]]; then git describe --always --dirty; fi)

Linkerflags = -s -X github.com/cybrarymin/polkadot_exporter/cmd/srv.BuildTime=${current_time} -X github.com/cybrarymin/greenlight/cmd/srv.Version=${git_version}
.PHONY: build/exporter
build/exporter:
	@go mod tidy
	GOOS=linux GOARCH=amd64 go build -ldflags="${Linkerflags}" -o=./bin/polkadot-exporter-linux-amd64 ./
	GOOS=darwin GOARCH=arm64 go build -ldflags="${Linkerflags}" -o=./bin/polkadot-exporter-darwin-arm64 ./
	go build -o=./bin/polkadot-exporter-local-compatible -ldflags="${Linkerflags}" ./

## run/exporter: run the application
.PHONY: run/exporter
run/exporter/http:
	@go run main.go 

run/exporter/https:
	@go run main.go --listen-addr ${LISTEN_ADDR} --crt ${CERT_FILE} --crt-key ${CERT_KEY_FILE}

#================================================================#
# QUALITY CHECK , LINTING, Vendoring
#================================================================#
.PHONY: audit
audit:
	@echo "Tidying and verifying golang packages and module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting codes..."
	go fmt ./...
	@echo "Vetting codes..."
	go vet ./...
	@echo "Static Checking of the code..."
	staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...

.PHONY: vendor
vendor:
	@echo "Tidying and verifying golang packages and module dependencies..."
	go mod verify
	go mod tidy
	@echo "Vendoring all golang dependency modules and packages..."
	go mod vendor

