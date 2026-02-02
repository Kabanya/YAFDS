package utils

import (
	"bufio"
	"os"
	"strings"
)

const PathToEnv = "config/.env"

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		sep := ""
		if strings.Contains(line, ":=") {
			sep = ":="
		} else if strings.Contains(line, "=") {
			sep = "="
		}
		if sep == "" {
			continue
		}

		parts := strings.SplitN(line, sep, 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Expand variables like $(VAR)
			for {
				start := strings.Index(value, "$(")
				if start == -1 {
					break
				}
				end := strings.Index(value[start:], ")")
				if end == -1 {
					break
				}
				end += start
				varName := value[start+2 : end]
				varValue := os.Getenv(varName)
				value = value[:start] + varValue + value[end+1:]
			}

			if existing, ok := os.LookupEnv(key); ok && existing != "" {
				continue
			}

			os.Setenv(key, value)
		}
	}
	return scanner.Err()
}
