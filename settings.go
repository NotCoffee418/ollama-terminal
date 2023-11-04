package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Exports struct {
	model       string
	alias       string
	isUndefined bool
}

const settingsStartString = "# START: ollama-terminal"
const settingsEndString = "# END: ollama-terminal"

func (e *Exports) Save() {
	// Load the current bashrc for saving
	contentStr := readBashrc()
	startIndex, endIndex, found := getSettingsPosition(contentStr)
	var preText string
	var postText string
	if found {
		preText = contentStr[:startIndex]
		postText = contentStr[endIndex:]
	} else {
		preText = contentStr + "\n"
		postText = ""
	}

	// Write the new settings contents
	sb := strings.Builder{}
	sb.WriteString(preText)
	sb.WriteString(settingsStartString)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("export OLLAMA_TERMINAL_MODEL=%s\n", e.model))
	sb.WriteString(fmt.Sprintf("%s() { /usr/local/bin/ollama-terminal run \"$@\"; }\n", e.alias))
	sb.WriteString(settingsEndString)
	sb.WriteString("\n")
	sb.WriteString(postText)

	// Write the new bashrc
	bashrcPath := GetRcPath()
	err := os.WriteFile(bashrcPath, []byte(sb.String()), 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing to ~/.bashrc: %s", err))
	}
	sourceBashrc()
	fmt.Println("Updated")
}

func LoadSettings() *Exports {
	contentStr := readBashrc()

	// First run, expect no settings
	startIndex, endIndex, found := getSettingsPosition(contentStr)
	if !found {
		return &Exports{isUndefined: true}
	}
	exports := &Exports{
		isUndefined: false,
	}

	// Extract lines without the comments
	lines := strings.Split(contentStr[startIndex:endIndex], "\n")
	if len(lines) < 2 {
		panic(fmt.Sprintf("Corrupt bashrc. Expect %s and %s to be on separate lines", settingsStartString, settingsEndString))
	}
	lines = lines[1 : len(lines)-1]

	aliasRx := regexp.MustCompile(`(\w+)\(\) { \/usr\/local\/bin\/ollama-terminal run "\$@"; }`)
	for _, line := range lines {
		if strings.HasPrefix(line, "export OLLAMA_TERMINAL_MODEL=") {
			exports.model = strings.TrimPrefix(line, "export OLLAMA_TERMINAL_MODEL=")
		} else if aliasRx.MatchString(line) {
			exports.alias = aliasRx.FindStringSubmatch(line)[1]
		}
	}
	return exports
}

// Returns start and end of settings related to this app, false if not found
func getSettingsPosition(contents string) (int, int, bool) {
	startIndex := -1
	endIndex := -1
	startIndex = strings.Index(contents, settingsStartString)
	if startIndex == -1 {
		return -1, -1, false
	}
	endIndex = strings.Index(contents, settingsEndString)
	if endIndex == -1 {
		return -1, -1, false
	}
	endIndex += len(settingsEndString)

	if startIndex > endIndex {
		panic(fmt.Sprintf(".bashrc is corrupted, %s comes after %s", settingsStartString, settingsEndString))
	}
	return startIndex, endIndex, true
}

func GetRcPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting the home directory: %s", err))
	}
	bashrcPath := filepath.Join(homeDir, ".bashrc")

	if _, err := os.Stat(bashrcPath); err != nil {
		panic(fmt.Sprintf("Error accessing ~/.bashrc: %s", err))
	}
	return bashrcPath
}

func readBashrc() string {
	bashrcPath := GetRcPath()
	contentRaw, err := os.ReadFile(bashrcPath)
	if err != nil {
		panic(fmt.Sprintf("Error reading ~/.bashrc: %s", err))
	}
	return string(contentRaw)
}
