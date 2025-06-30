package api

import (
	"fmt"
	"time"
)

type MetricsFormatter struct{}

func (f *MetricsFormatter) FormatMetrics(rawData map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})

    summary, ok := rawData["summary"].(map[string]interface{})
    if !ok {
        result["recorded_at"] = time.Now().UTC().Format(time.RFC3339)
        return result
    }

    for k, v := range summary {
        result[k] = v
    }

    intFields := []string{"lines", "functions", "classes", "comments"}
    stringFields := []string{"comment_percentage"}

    for _, field := range intFields {
        if v, ok := result[field]; ok {
            switch val := v.(type) {
            case int:
                result[field] = val
            case float64:
                result[field] = int(val)
            case string:
                var intVal int
                fmt.Sscanf(val, "%d", &intVal)
                result[field] = intVal
            default:
                result[field] = 0
            }
        } else {
            result[field] = 0
        }
    }

    for _, field := range stringFields {
        if v, ok := result[field]; ok {
            result[field] = fmt.Sprintf("%v", v)
        } else {
            result[field] = ""
        }
    }

    result["recorded_at"] = time.Now().UTC().Format(time.RFC3339)
    return result
}