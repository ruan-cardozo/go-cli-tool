# Guia do Usuário - Go CLI Tool

## Introdução ao Produto

O **Go CLI Tool** é uma ferramenta de linha de comando desenvolvida em Go para realizar análises de código JavaScript. Ele é projetado para desenvolvedores e equipes que desejam obter insights detalhados sobre seus projetos, como contagem de linhas, análise de comentários, classes, funções, dependências externas e nativas, além de análise de indentação. Este guia fornece uma visão geral do produto, instruções de instalação, recursos, casos de uso e dicas para solução de problemas.

---

## Instalação e Configuração

### Pré-requisitos
- **Go** instalado na máquina (versão 1.19 ou superior).
- **Git** para clonar o repositório.

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

### Configuração
- Para configurar o ambiente, certifique-se de que o diretório do projeto tenha permissões adequadas para leitura e escrita.
- Caso deseje salvar relatórios em um local específico, use a flag `-o` para especificar o caminho de saída.

---

## Recursos e Casos de Uso do Produto

### Recursos Principais
- **Contador de Linhas:** Conta o número total de linhas de código em arquivos ou diretórios.
- **Contador de Classes e Funções:** Identifica e conta classes e funções declaradas no código.
- **Contador de Comentários:** Analisa e contabiliza comentários no código.
- **Analisador de Indentação:** Verifica o uso de tabs ou espaços e calcula os níveis de indentação.
- **Analisador de Dependências:** Classifica dependências externas e nativas em projetos JavaScript/Node.js.
- **Proporção Código/Comentários:** Calcula a proporção entre linhas de código e comentários.

### Casos de Uso
- **Desenvolvedores Individuais:** Obter insights sobre a qualidade e estrutura do código.
- **Equipes de Desenvolvimento:** Garantir consistência e documentação adequada em projetos colaborativos.
- **Revisores de Código:** Identificar rapidamente áreas de melhoria no código.

---

## Como Operar o Produto

### Comando de ajuda
Para visualizar os comandos disponíveis:
```bash
./go-cli-tool --help
```

### Analisar um único arquivo
```bash
./go-cli-tool analyze -f caminho/para/arquivo.js -o .
```

### Analisar um diretório inteiro
```bash
./go-cli-tool analyze -d caminho/para/diretorio -o .
```

### Gerar saída detalhada
Para gerar uma saída detalhada com informações por arquivo:
```bash
./go-cli-tool analyze -d caminho/para/diretorio --detailed -o .
```

---

## Dicas de Solução de Problemas

### Problema: Comando não encontrado
- Certifique-se de que o binário foi compilado corretamente.
- Verifique se o binário está no diretório atual ou no PATH do sistema.

### Problema: Permissão negada ao executar o binário
- Use o comando `chmod` para garantir permissões de execução:
  ```bash
  chmod +x go-cli-tool
  ```

### Problema: Relatórios não são gerados
- Verifique se o caminho de saída especificado com `-o` é válido e possui permissões de escrita.

---

## Perguntas Frequentes (FAQs)

### 1. O que é o Go CLI Tool?
É uma ferramenta de linha de comando para análise de código JavaScript, desenvolvida em Go.

### 2. Quais linguagens de programação são suportadas?
Atualmente, o Go CLI Tool é projetado para analisar projetos JavaScript/Node.js.

### 3. Posso usar a ferramenta em sistemas Windows?
Sim, a ferramenta é compatível com Windows, Linux e macOS.

### 4. Como salvar os resultados em um arquivo JSON?
Use a flag `-o` para especificar o caminho de saída:
```bash
./go-cli-tool analyze -d caminho/para/diretorio -o ./resultado.json
```

---

## Glossário

- **CLI (Command Line Interface):** Interface de linha de comando usada para interagir com o programa.
- **Dependências Externas:** Bibliotecas ou pacotes de terceiros usados no projeto.
- **Indentação:** Espaçamento usado para organizar o código, como tabs ou espaços.
- **LOC (Lines of Code):** Número total de linhas de código em um arquivo ou projeto.
- **JSON:** Formato de dados usado para saída estruturada.

---

Se precisar de mais informações ou suporte, consulte o arquivo `README.md` ou entre em contato com os desenvolvedores do projeto.
