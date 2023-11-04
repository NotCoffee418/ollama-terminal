package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		helpText := `Usage: ollama-terminal [command]

Commands:
  run       Runs the program with the default model. Usage: 'ollama-terminal run [prompt]'
  install   Installs the application for use.
  set       Sets the default model to use. Usage: 'ollama-terminal set model llama2'
            alias        Sets alias for ollama-terminal. Usage: 'ollama-terminal set alias llama'
            model        Sets the default model to use.
  help      Shows this menu.
`
		fmt.Print(helpText)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}
	switch os.Args[1] {
	case "run":
		run(os.Args[2:])
	case "install":
		install()
	case "set":
		if len(os.Args) < 4 {
			panic("Not enough arguments for set command")
		}
		set(os.Args[2], os.Args[3])
	case "help":
		flag.Usage()
	default:
		panic(fmt.Sprintf("Unknown command: %s", os.Args[1]))
	}
}

func install() {
	s := LoadSettings()
	if !s.isUndefined {
		fmt.Println("ollama-terminal is already installed.")
		return
	}
	s.isUndefined = false
	s.model = "llama2"
	s.alias = "llama"
	s.Save()
}

func run(params []string) {
	// Ensure environment is set up correctly
	sudoCheck()
	if !execCheck("ollama") {
		panic("ollama is not installed. Please install it first. (https://ollama.ai/download)")
	}
	installCheck()

	// Parse any prompt
	model := os.Getenv("OLLAMA_TERMINAL_MODEL")
	if model == "" {
		panic("OLLAMA_TERMINAL_MODEL is not set. Please run `ollama-terminal install`")
	}
	if strings.HasSuffix(model, ":text") {
		// Dont use text models
		model = model[:len(model)-5]
		fmt.Println("Warning: Text models are not supported. Using the non-text version of the model.")
	}
	initialPrompt := strings.Join(params, " ")

	// Run ollama with the API
	StartOllama(model, initialPrompt)
}

func set(key string, value string) {
	s := LoadSettings()
	switch key {
	case "alias":
		if value == "ollama" || value == "ollama-terminal" {
			panic("Cannot set alias to ollama or ollama-terminal. It will break stuff.")
		}
		s.alias = value
	case "model":
		s.model = value
	default:
		panic(fmt.Sprintf("Unknown setting: %s", key))
	}
	s.Save()
}
