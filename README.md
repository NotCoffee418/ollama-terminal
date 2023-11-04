# ollama-terminal

## Installation
ollama-terminal is a terminal application that uses [ollama](https://ollama.ai) to generate a prompt through the terminal, allowing you to use non-text models in the terminal.
You can specify an alias for the application, and the model to use for the prompt.

**Features on request**
- [ ] Customizable preprompt
- [ ] Chat finetuning

## Install and update
```bash
curl -s https://raw.githubusercontent.com/NotCoffee418/ollama-terminal/main/install.sh | sudo bash
```

## Usage
```bash
llama [optional initial prompt}
```
See Configuration section for changing `llama` to a different alias.

## Configuration
The alias is the command you run to activate the application. By default the alias is `llama`.
To change the alias run the following command:
```bash
/local/usr/bin/ollama-terminal set alias [new alias]
```

To set the model ollama-terminal uses to generate the prompt run the following command:
```bash
/local/usr/bin/ollama-terminal set model [model name]
```
You can find a list of usable models on the [ollama website](https://ollama.ai/library).
