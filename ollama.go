package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func StartOllama(modelName string, initialPrompt string) {
	ollamaUpCheck()
	ollamaModelPrep(modelName)

	fmt.Println("Type /bye to exit")

	initialPromptHandled := initialPrompt == ""
	for {
		// Handle initial prompt or prompt input
		var promptStr string
		if !initialPromptHandled {
			promptStr = initialPrompt
			initialPromptHandled = true
		} else {
			promptStr = getInputPrompt()
		}

		// Handle special commands
		if promptStr == "/bye" {
			return
		}

		// Prompt the model
		streamedPrompt(modelName, promptStr)
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

func getInputPrompt() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPrompt: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from input:", err)
	}
	text = strings.TrimSpace(text)
	return text
}
