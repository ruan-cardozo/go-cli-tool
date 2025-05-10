package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CountDependenciesAnalyzerImpl implements the logic to count external dependencies
type CountDependenciesAnalyzerImpl struct{}

type DependencyResult struct {
    TotalDependencies int      `json:"total_dependencies"`
    Dependencies      []string `json:"dependencies"`
    NativeModules     []string `json:"native_modules"`
}

// List of regexes to capture different dependency patterns
var dependencyRegexes = []*regexp.Regexp{
    // import defaultExport from 'module-name';
    regexp.MustCompile(`(?m)^import\s+\w+\s+from\s+['"]([^'"]+)['"]`),

    // import * as name from 'module-name';
    regexp.MustCompile(`(?m)^import\s+\*\s+as\s+\w+\s+from\s+['"]([^'"]+)['"]`),

    // import { something } from 'module-name';
    regexp.MustCompile(`(?m)^import\s+{[^}]+}\s+from\s+['"]([^'"]+)['"]`),

    // import 'module-name';
    regexp.MustCompile(`(?m)^import\s+['"]([^'"]+)['"]`),

    // require('module-name') or require('module-name').something()
    regexp.MustCompile(`(?m)require\(\s*['"]([^'"]+)['"]\s*\)`),
}

func (a *CountDependenciesAnalyzerImpl) CountDependenciesByFilePath(filePath string) (map[string]interface{}, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    externalDependencies := make(map[string]struct{})
    nativeModules := make(map[string]struct{})

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())

        // Check each regex for a match
        for _, r := range dependencyRegexes {
            if matches := r.FindStringSubmatch(line); matches != nil {
                dependency := matches[1]
                if dependency != "" {
                    normalizedDependency := normalizeModuleName(dependency)
                    if isNativeModule(normalizedDependency) {
                        nativeModules[normalizedDependency] = struct{}{}
                    } else if isExternalDependency(normalizedDependency) {
                        externalDependencies[normalizedDependency] = struct{}{}
                    }
                    break
                }
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "total_dependencies": len(externalDependencies),
        "dependencies":       mapKeysToSlice(externalDependencies),
        "native_modules":     mapKeysToSlice(nativeModules),
    }, nil
}

func normalizeModuleName(moduleName string) string {
    // Remove o prefixo "node:" para normalizar os módulos nativos
    if strings.HasPrefix(moduleName, "node:") {
        return strings.TrimPrefix(moduleName, "node:")
    }
    return moduleName
}

// CountDependenciesByDirectory analyzes all JavaScript files in a directory for external and native dependencies
func (a *CountDependenciesAnalyzerImpl) CountDependenciesByDirectory(directoryPath string) (map[string]interface{}, error) {
    results := make(map[string]interface{})

    err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && strings.HasSuffix(path, ".js") {
            result, err := a.CountDependenciesByFilePath(path)
            if err != nil {
                return err
            }
            results[path] = result
        }
        return nil
    })

    return results, err
}

// isExternalDependency checks if a dependency is external (not relative)
func isExternalDependency(dependency string) bool {
    // External dependencies do not start with "./", "../", or "/"
    return !strings.HasPrefix(dependency, "./") &&
        !strings.HasPrefix(dependency, "../") &&
        !strings.HasPrefix(dependency, "/")
}

// isNativeModule checks if a dependency is a native Node.js module
func isNativeModule(dependency string) bool {
    nativeModules := map[string]bool{
        "assert":          true,
        "async_hooks":     true,
        "buffer":          true,
        "child_process":   true,
        "cluster":         true,
        "console":         true,
        "constants":       true,
        "crypto":          true,
        "dgram":           true,
        "diagnostics_channel": true,
        "dns":             true,
        "domain":          true,
        "events":          true,
        "fs":              true,
        "http":            true,
        "http2":           true,
        "https":           true,
        "inspector":       true,
        "module":          true,
        "net":             true,
        "os":              true,
        "path":            true,
        "perf_hooks":      true,
        "process":         true,
        "punycode":        true,
        "querystring":     true,
        "readline":        true,
        "repl":            true,
        "stream":          true,
        "string_decoder":  true,
        "sys":             true,
        "timers":          true,
        "timers/promises": true,
        "tls":             true,
        "trace_events":    true,
        "tty":             true,
        "url":             true,
        "util":            true,
        "v8":              true,
        "vm":              true,
        "worker_threads":  true,
        "zlib":            true,
    }
    return nativeModules[dependency]
}

// appendIfMissing appends a dependency to the list if it is not already present
func appendIfMissing(slice []string, item string) []string {
    for _, existing := range slice {
        if existing == item {
            return slice
        }
    }
    return append(slice, item)
}

func mapKeysToSlice(m map[string]struct{}) []string {
    keys := make([]string, 0, len(m))
    for key := range m {
        keys = append(keys, key)
    }
    return keys
}