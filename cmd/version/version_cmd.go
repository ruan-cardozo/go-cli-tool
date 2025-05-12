package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var Version = "dev"

type GithubRelease struct {
    TagName    string `json:"tag_name"`
    Name       string `json:"name"`
    Published  string `json:"published_at"`
    Prerelease bool   `json:"prerelease"`
}

func VersionCommand() *cobra.Command {
    versionCmd := &cobra.Command{
        Use:   "version",
        Short: "Verifica a versão atual do CLI",
        Long:  "Exibe a versão instalada e verifica se há atualizações disponíveis no GitHub",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Versão instalada: %s\n", Version)

            latest, err := getLatestRelease()
            if err != nil {
                fmt.Printf("Erro ao verificar atualizações: %s\n", err)
                return
            }
            
            latestVersion := strings.TrimPrefix(latest.TagName, "v")
            fmt.Printf("Versão mais recente: %s\n", latestVersion)

            if Version != latestVersion && Version != "dev" {
                fmt.Println("Uma nova versão está disponível!")
                fmt.Printf("Você pode atualizar baixando a nova versão em:\n")
                fmt.Printf("https://github.com/ruan-cardozo/go-cli-tool/releases/tag/%s\n", latest.TagName)
            } else if Version == latestVersion {
                fmt.Println("Você está usando a versão mais recente!")
            }
        },
    }
    
    return versionCmd
}

func getLatestRelease() (*GithubRelease, error) {

    apiURL := "https://api.github.com/repos/ruan-cardozo/go-cli-tool/releases/latest"

    resp, err := http.Get(apiURL)
    if err != nil {
        return nil, fmt.Errorf("falha ao conectar com a API do GitHub: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API retornou status inválido: %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("erro ao ler resposta da API: %w", err)
    }
    
    var release GithubRelease
    if err := json.Unmarshal(body, &release); err != nil {
        return nil, fmt.Errorf("erro ao decodificar resposta JSON: %w", err)
    }
    
    return &release, nil
}