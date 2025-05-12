#!/bin/bash

# Nome do binário
BINARY_NAME="go-cli-tool"

# Diretório de destino para o binário
INSTALL_DIR="/usr/local/bin"

export PATH=$PATH:/usr/local/go/bin

# Função para compilar o projeto
build() {
    echo "Building the project..."

    VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")

    go build -ldflags "-X go-cli-tool/cmd/version.Version=${VERSION#v}" -o $BINARY_NAME

    if [ $? -ne 0 ]; then
        echo "Build failed!"
        exit 1
    fi
    echo "Build successful!"
}

# Função para instalar o binário
install() {
    echo "Installing the binary to $INSTALL_DIR..."
    sudo cp $BINARY_NAME $INSTALL_DIR
    if [ $? -ne 0 ]; then
        echo "Installation failed!"
        exit 1
    fi
    echo "Installation successful!"
}

# Função para limpar os arquivos gerados
clean() {
    echo "Cleaning up..."
    rm -f $BINARY_NAME
    echo "Clean up successful!"
}

# Função para rodar o projeto
run() {
    build
    echo "Running the project..."
    ./$BINARY_NAME
}

# Verifica o argumento passado para o script
case "$1" in
    build)
        build
        ;;
    install)
        build
        install
        ;;
    clean)
        clean
        ;;
    run)
        run
        ;;
    all)
        clean
        build
        install
        ;;
    *)
        echo "Usage: $0 {build|install|clean|run|all}"
        exit 1
        ;;
esac