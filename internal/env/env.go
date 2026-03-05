package env

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ParseEnvLines reads VAR=value lines from r, skipping blank lines and # comments.
func ParseEnvLines(r io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result, scanner.Err()
}

// ProcessEnv returns the current process environment as a map.
func ProcessEnv() map[string]string {
	result := make(map[string]string)
	for _, pair := range os.Environ() {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// Merge returns a new map combining all input maps; later maps overwrite earlier ones.
func Merge(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
