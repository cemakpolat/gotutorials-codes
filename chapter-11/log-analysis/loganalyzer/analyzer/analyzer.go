package analyzer

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func AnalyzeLogs(filename string, filter string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	severityCounts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	regex := regexp.MustCompile(filter)

	for scanner.Scan() {
		line := scanner.Text()
		if filter != "" && !regex.MatchString(line) {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		firstPart := parts[0]
		if strings.Contains(firstPart, "[") && strings.Contains(firstPart, "]") {
			severity := strings.ReplaceAll(firstPart, "[", "")
			severity = strings.ReplaceAll(severity, "]", "")
			severityCounts[severity]++
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return severityCounts, nil
}
