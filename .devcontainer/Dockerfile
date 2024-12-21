FROM mcr.microsoft.com/devcontainers/go:1.23-bullseye

RUN sudo apt update && sudo apt install -y \
    curl \
    zsh \
    sqlite3

USER vscode

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.2 && \
    go install go.uber.org/mock/mockgen@latest && \
    go install github.com/rakyll/gotest@latest && \
    go install go.uber.org/mock/mockgen@latest

RUN curl -L -o /tmp/sqlc.tgz https://downloads.sqlc.dev/sqlc_1.27.0_linux_amd64.tar.gz && \
    cd /usr/local/bin && sudo tar xf /tmp/sqlc.tgz