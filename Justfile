# Enable .env file support for local configuration
set dotenv-load

# Use bash with strict error checking
set shell := ["bash", "-uc"]

# Allow passing arguments to recipes
set positional-arguments

# Common command aliases for convenience
alias t := test
alias b := build
alias r := run
alias help := default

project_name := "go-scule"

# Build configuration
# Tags for conditional compilation
build_tags := ""
extra_tags := ""
all_tags := build_tags + " " + extra_tags

# Test configuration
# Settings for test execution and coverage
test_timeout := "5m"
coverage_threshold := "80"
bench_time := "2s"

# Go settings
# Core Go environment variables and configuration
export GOPATH := env_var_or_default("GOPATH", `go env GOPATH`)
export GOOS := env_var_or_default("GOOS", `go env GOOS`)
export GOARCH := env_var_or_default("GOARCH", `go env GOARCH`)
export CGO_ENABLED := env_var_or_default("CGO_ENABLED", "1")
go := env_var_or_default("GO", "go")
gobin := GOPATH + "/bin"

# Automatically detect version information from git
# Falls back to timestamp if not in a git repository
version := if `git rev-parse --git-dir 2>/dev/null; echo $?` == "0" {
    `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
} else {
    `date -u '+%Y%m%d-%H%M%S'`
}
git_commit := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`
git_branch := `git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown"`
build_time := `date -u '+%Y-%m-%d_%H:%M:%S'`
build_by := `whoami`

# Directories
# Project directory structure
root_dir := justfile_directory()
bin_dir := root_dir + "/bin"
dist_dir := root_dir + "/dist"
docs_dir := root_dir + "/docs"

# Build flags
# Linker flags for embedding version information
ld_flags := "-s -w \
    -X '$(go list -m)/pkg/version.Version=" + version + "' \
    -X '$(go list -m)/pkg/version.Commit=" + git_commit + "' \
    -X '$(go list -m)/pkg/version.Branch=" + git_branch + "' \
    -X '$(go list -m)/pkg/version.BuildTime=" + build_time + "' \
    -X '$(go list -m)/pkg/version.BuildBy=" + build_by + "'"

# Show available recipes with their descriptions
@default:
    just --list

# Initialize a new project with a basic structure and configuration
init:
    #!/usr/bin/env bash
    if [ ! -f "go.mod" ]; then
        {{go}} mod init "$(basename "$(pwd)")"
    fi
    if [ ! -f ".gitignore" ]; then
        curl -sL https://www.gitignore.io/api/go > .gitignore
    fi
    mkdir -p \
        main \
        testdata \
        .github/workflows
    if [ ! -f "main/main.go" ]; then
        mkdir -p main
        printf '%s\n' \
            'package main' \
            '' \
            'import "fmt"' \
            '' \
            'func main() {' \
            '    fmt.Println("Hello, World!")' \
            '}' \
            > main/main.go
    fi

# Build the project
build:
    mkdir -p {{bin_dir}}
    {{go}} build \
        -ldflags '{{ld_flags}}' \
        -o {{bin_dir}}/{{project_name}} \
        ./main

# Run the application
run: build
    {{bin_dir}}/{{project_name}}

# Install the application
install: build
    {{go}} install -tags '{{all_tags}}' -ldflags '{{ld_flags}}' ./main

# Generate code
generate:
    {{go}} generate ./...

# Run tests
test:
    {{go}} test -v -race -cover ./...

# Run tests with coverage
test-coverage:
    {{go}} test -v -race -coverprofile=coverage.out ./...
    {{go}} tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
    {{go}} test -bench=. -benchmem -run=^$ -benchtime={{bench_time}} ./...

# Format code
fmt:
    {{go}} fmt ./...

# Run linters
lint:
    {{gobin}}/golangci-lint run --fix

# Run go vet
vet:
    {{go}} vet ./...

# Cross-compile for all platforms
build-all:
    #!/usr/bin/env sh
    mkdir -p {{dist_dir}}
    for platform in \
        "linux/amd64/-" \
        "linux/arm64/-" \
        "linux/arm/7" \
        "darwin/amd64/-" \
        "darwin/arm64/-" \
        "windows/amd64/-" \
        "windows/arm64/-"; do
        os=$(echo $platform | cut -d/ -f1)
        arch=$(echo $platform | cut -d/ -f2)
        arm=$(echo $platform | cut -d/ -f3)
        output="{{dist_dir}}/{{project_name}}-${os}-${arch}$([ "$os" = "windows" ] && echo ".exe")"

        GOOS=$os GOARCH=$arch $([ "$arm" != "-" ] && echo "GOARM=$arm") \
        CGO_ENABLED={{CGO_ENABLED}} {{go}} build \
            -tags '{{all_tags}}' \
            -ldflags '{{ld_flags}}' \
            -o "$output" \
            ./main

        tar czf "$output.tar.gz" "$output"
        rm -f "$output"
    done

# Generate documentation
docs:
    mkdir -p {{docs_dir}}
    {{go}} doc -all > {{docs_dir}}/API.md

# Show version information
version:
    @echo "Version:    {{version}}"
    @echo "Commit:     {{git_commit}}"
    @echo "Branch:     {{git_branch}}"
    @echo "Built:      {{build_time}}"
    @echo "Built by:   {{build_by}}"
    @echo "Go version: $({{go}} version)"

