package main

import (
	"flag"
	"fmt"
	"os"
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

}

func set(key string, value string) {
	s := LoadSettings()
	switch key {
	case "alias":
		s.alias = value
	case "model":
		s.model = value
	default:
		panic(fmt.Sprintf("Unknown setting: %s", key))
	}
	s.Save()
}
