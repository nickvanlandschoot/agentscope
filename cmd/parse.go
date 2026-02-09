package cmd

import (
	"bufio"
	"slices"
	"strings"
)

func Parse(content string) map[string]string {
	options := make(map[string]string)
	prefixes := []string{
		"name",
		"position",
		"enabledByDefault",
	}
	scanner := bufio.NewScanner(strings.NewReader(content))

	started := false
	for scanner.Scan() {
		line := scanner.Text()
		if started == true && strings.TrimSpace(line) == "---" {
			break
		}
		if strings.TrimSpace(line) == "---" {
			started = true
		}

		if started == false {
			continue
		}

		if strings.Contains(line, ":") {
			key, val, _ := strings.Cut(line, ":")
			if slices.Contains(prefixes, key) {
				options[key] = strings.TrimSpace(val)
			}
		}
	}

	return options
}
