package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func StartOllama(modelName string, initialPrompt string) {
	ollamaUpCheck()
	ollamaModelPrep(modelName)
}

func ollamaUpCheck() {
	resp, err := httpClient.Get("http://localhost:11434")
	if err != nil {
		panic("Ollama service is not running. Please ensure ollama it is installed and it's service running")
	}

	rawResp, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || err != nil || string(rawResp) != "Ollama is running" {
		panic("Ollama service is not running. Please ensure ollama it is installed and it's service running")
	}
}

func ollamaModelPrep(modelName string) {
	if checkModelExists(modelName) {
		return
	}
	cmd := exec.Command("ollama", "pull", modelName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(fmt.Sprintf("Error executing command: %s\n", err))
	}
}
