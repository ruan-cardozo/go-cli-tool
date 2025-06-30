package send_metrics

import (
	"encoding/json"
	"fmt"
	"go-cli-tool/internal/api"
	"go-cli-tool/internal/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
    apiURL      string
    metricsFile string
    projectID   string
)

var SendMetricsCmd = &cobra.Command{
    Use:   "send-metrics",
    Short: "Send analysis metrics to an API",
    Run: func(cmd *cobra.Command, args []string) {
        if metricsFile == "" {
            fmt.Fprintf(cmd.OutOrStdout(), "%sInvalid input. Please provide a metrics file (-m).%s\n", utils.RED, utils.RESET_COLOR)
            return
        }

        if projectID == "" {
            fmt.Fprintf(cmd.OutOrStdout(), "%sInvalid input. Please provide a project ID (-p).%s\n", utils.RED, utils.RESET_COLOR)
            return
        }

        finalAPIURL := strings.Replace(apiURL, "<id>", projectID, 1)

        jsonData, err := os.ReadFile(metricsFile)
        if err != nil {
            fmt.Fprintf(cmd.OutOrStdout(), "%sError reading metrics file: %v%s\n", utils.RED, err, utils.RESET_COLOR)
            return
        }

        var rawData map[string]interface{}
        if err := json.Unmarshal(jsonData, &rawData); err != nil {
            fmt.Fprintf(cmd.OutOrStdout(), "%sError parsing metrics JSON: %v%s\n", utils.RED, err, utils.RESET_COLOR)
            return
        }

        formatter := &api.MetricsFormatter{}
        payload := formatter.FormatMetrics(rawData)

        client := api.NewAPIClient(finalAPIURL, "")

        fmt.Fprintf(cmd.OutOrStdout(), "%sSending metrics to API: %s%s\n", utils.BLUE, finalAPIURL, utils.RESET_COLOR)
        if err := client.SendMetrics(payload); err != nil {
            fmt.Fprintf(cmd.OutOrStdout(), "%sError sending metrics to API: %v%s\n", utils.RED, err, utils.RESET_COLOR)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "%sMetrics successfully sent to API!%s\n", utils.GREEN, utils.RESET_COLOR)
    },
}

func init() {
    SendMetricsCmd.Flags().StringVarP(&metricsFile, "metrics-file", "m", "", "Path to the JSON metrics file")
    SendMetricsCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "Project ID for the API")
    SendMetricsCmd.Flags().StringVar(&apiURL, "api-url", "http://localhost:8030/api/v1/project/metric/<id>", "API URL (use <id> as placeholder for project ID)")
    
    SendMetricsCmd.MarkFlagRequired("metrics-file")
    SendMetricsCmd.MarkFlagRequired("project-id")
}