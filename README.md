<h1 align="center"> Go Cli Tool </h1>

## 📜 Descrição do Projeto

Este projeto é uma ferramenta de análise de código JavaScript que realiza múltiplas análises, como contagem de linhas, comentários, classes, funções, análise de indentação e dependências externas.

---

## 🚀 Funcionalidades

<!-- Liste as funcionalidades principais do seu projeto. -->

- `Contador de linhas`: Conta o número total de linhas de código em um arquivo ou diretório, desconsiderando linhas em branco.
- `Contador de Classes e Funções`: Utiliza expressões regulares para identificar e contar a quantidade de classes e funções declaradas em arquivos ou diretórios.
- `Contador de Comentários`: Identifica e contabiliza comentários presentes em arquivos ou diretórios com o uso de expressões regulares.
- `Analisador de Identação`: Analisa a identação de arquivos ou diretório e retorna informações se uso tabs ou espaços e os levels de identação presente no arquivo
- `Analisador de Percentual de Código vs. Comentários`: Calcula a proporção entre linhas de comentários e o total de linhas de código, fornecendo uma visão geral da documentação no projeto.
- `Analisador de Dependências Externas e Nativas`: Analisa declarações de `import` e `require` em projetos JavaScript/Node.js, classificando as dependências como externas (de pacotes) ou nativas (do Node.js).

---

## 🛠️ Tecnologias Utilizadas

<!-- Liste as tecnologias, linguagens ou frameworks usados no projeto. -->

- Go (Golang)
- Cobra CLI
- JSON
- Vscode
- Shell Script
- Git e Github

---

## 📂 Estrutura do Projeto

- `cmd/`: Contém os comandos principais da CLI, organizados por funcionalidades específicas. Exemplos:
  - `count-class-and-functions/`: Comando para contar classes e funções.
  - `count-comments/`: Comando para contar comentários.
  - `count-lines/`: Comando para contar linhas de código.
  - `dependencies/`: Comando para analisar dependências externas e nativas.
  - `run-all-commands/`: Comando para executar todas as análises de uma vez.

- `internal/`: Contém os módulos internos que implementam a lógica principal do projeto. Subdiretórios incluem:
  - `analyzer/`: Implementações dos analisadores, como contagem de linhas, classes, funções e análise de dependências.
  - `policies/`: Regras e políticas usadas pelos analisadores.
  - `utils/`: Funções auxiliares e constantes reutilizáveis em todo o projeto.

- `templates/`: Contém arquivos de template usados para gerar relatórios, como:
  - `report.html`: Template HTML para relatórios de análise.
  - `template.go`: Código Go para manipulação de templates.

- `javascript-tests/`: Scripts de teste em JavaScript para validar funcionalidades específicas.

- Arquivos adicionais:
  - `go.mod` e `go.sum`: Gerenciamento de dependências do Go.
  - `build_and_install.sh`: Script para compilar e instalar o projeto.
  - `README.md`: Documentação do projeto.
  - `LICENSE`: Licença do projeto.

## 📦 Como Instalar e Rodar o Projeto

### Pré-requisitos

<!-- Liste os pré-requisitos necessários para rodar o projeto. -->

- Go instalado na máquina.
- Git para clonar o repositório.

### Passos para Instalação

1. Clone o repositório:
   ```bash
   git clone https://github.com/ruan-cardozo/go-cli-tool
   ```
2. Navegue até o diretório do projeto:
   ```bash
   cd go-cli-tool
   ```
3. Compile o projeto:
   ```bash
   go build -o go-cli-tool
   ```
4. Execute o binário:
   ```bash
   ./go-cli-tool
   ```

---

## 📊 Exemplos de Uso

<!-- Mostre exemplos de como usar o projeto, incluindo comandos e saídas esperadas. -->

```bash
# Comando de help do cli para facilitar a utilização 
./go-cli-tool --help

# Analisar um único arquivo
./go-cli-tool analyze -f caminho/para/arquivo.js -o .

# Analisar um diretório inteiro
./go-cli-tool analyze -d caminho/para/diretorio -o .
```

---

## 🧑‍💻 Contribuidores

<!-- Liste os contribuidores do projeto. -->

| [<img loading="lazy" src="https://github.com/ruan-cardozo.png" width=115><br><sub>Ruan Cardozo</sub>](https://github.com/ruan-cardozo) |  [<img loading="lazy" src="https://github.com/guimachado1.png" width=115><br><sub>Guilherme Machado</sub>](https://github.com/guimachado1) |  [<img loading="lazy" src="https://github.com/guilherme-kopsch.png" width=115><br><sub>Guilherme Kopsch</sub>](https://github.com/guilherme-kopsch) |
| :---: | :---: | :---: |

---

## 📜 Licença

Este projeto está licenciado sob a licença MIT. Consulte o arquivo `LICENSE` para mais detalhes.
