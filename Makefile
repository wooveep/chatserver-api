PROJECT_NAME=chatserver-api-app
BIN_NAMES=chatserver-api
GOARCHS=amd64
GOARCHS_MAC=amd64
go_version = $(shell go version)
commit_id = $(shell git rev-parse HEAD)
branch_name = $(shell git name-rev --name-only HEAD)
build_time = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
app_version = 0.0.1-dev
version_package = chatserver-api/utils/version

dev: mac

default: all

all: windows linux mac

prepare:
	@mkdir -p dist

windows: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS); do \
			mkdir -p dist/windows_$$GOARCH; \
			OOSG=windows GOARCH=$$GOARCH go build -ldflags \
			"-X ${version_package}.CommitId=${commit_id} \
			-X ${version_package}.BranchName=${branch_name} \
			-X ${version_package}.BuildTime=${build_time} \
			-X ${version_package}.AppVersion=${app_version}" \
			-o dist/windows_$$GOARCH/$$BIN_NAME.exe cmd/main.go; \
		done \
	done

linux: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS); do \
			mkdir -p dist/linux_$$GOARCH; \
			GOOS=linux GOARCH=$$GOARCH go  build -ldflags \
			"-X ${version_package}.CommitId=${commit_id} \
			-X ${version_package}.BranchName=${branch_name} \
			-X ${version_package}.BuildTime=${build_time} \
			-X ${version_package}.AppVersion=${app_version}"  \
			 -o dist/linux_$$GOARCH/$$BIN_NAME cmd/main.go; \
		done \
	done

mac: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS_MAC); do \
			mkdir -p dist/mac_$$GOARCH; \
			GOOS=darwin GOARCH=$$GOARCH go  build -ldflags \
			"-X ${version_package}.CommitId=${commit_id} \
			-X ${version_package}.BranchName=${branch_name} \
			-X ${version_package}.BuildTime=${build_time} \
			-X ${version_package}.AppVersion=${app_version}"  \
			-o dist/mac_$$GOARCH/$$BIN_NAME cmd/main.go; \
		done \
	done

package: all
	for GOARCH in $(GOARCHS); do \
		zip -q -r dist/$(PROJECT_NAME)-windows-$$GOARCH.zip dist/windows_$$GOARCH/; \
		zip -q -r dist/$(PROJECT_NAME)-linux-$$GOARCH.zip dist/linux_$$GOARCH/; \
	done

	for GOARCH in $(GOARCHS_MAC); do \
		zip -q -r dist/$(PROJECT_NAME)-mac-$$GOARCH.zip dist/mac_$$GOARCH/; \
	done

	ARCH_RELEASE_DIRS=$$(find dist -type d -name "*_*"); \
	for ARCH_RELEASE_DIR in $$ARCH_RELEASE_DIRS; do \
		cp conf/config.default.toml $$ARCH_RELEASE_DIR/config.toml; \
		rm -rfd $$ARCH_RELEASE_DIR; \
	done

test:
	go test -v ./...

# Usage: make run APP=greet -- -h
run:
	go run cmd/main.go

clean:
	rm -rfd dist

.PHONY: all, default, clean