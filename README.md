# Go Cli Tool

Este é o projeto Go Cli Tool, desenvolvido para analisar e processar dados de código Javascript.

## Instalação

Para instalar o Go (Golang) e configurar o ambiente, siga os passos abaixo:

### Passo 1: Baixar o Go

Visite a página de [downloads do Go](https://golang.org/dl/) e baixe a versão 1.21 para o seu sistema operacional.

### Passo 2: Instalar o Go

Siga as instruções de instalação específicas para o seu sistema operacional:

- **Windows**: Execute o instalador baixado e siga as instruções na tela.
- **macOS**: Abra o pacote `.pkg` baixado e siga as instruções na tela.
- **Linux**: Extraia o arquivo tar.gz baixado e mova-o para `/usr/local`:

```sh
tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
```

### Passo 3: Configurar o PATH

Adicione o diretório do Go ao seu PATH. Adicione as seguintes linhas ao seu arquivo de perfil (`~/.profile`, `~/.bashrc`, `~/.zshrc`, etc.):

```sh
export PATH=$PATH:/usr/local/go/bin
```

### Passo 4: Verificar a Instalação

Verifique se o Go foi instalado corretamente executando o comando:

```sh
go version
```

Você deve ver a versão 1.21 do Go instalada.

## Uso

Para usar o CLI, siga os passos abaixo:

### Passo 1: Clonar o Repositório

Clone o repositório do Code Mind Analyzer:

```sh
git clone https://github.com/ruan-cardozo/go-cli-tool.git
cd go-cli-tool
```

### Passo 2: Construir e Instalar

Execute o script `build_and_install.sh` com os seguintes comandos:

```sh
./build_and_install.sh build
./build_and_install.sh install
./build_and_install.sh clean
```