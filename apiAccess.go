package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jmorganca/ollama/api"
)

func checkModelExists(modelName string) bool {
	// If no tag is specified, use :latest
	re := regexp.MustCompile(`(.+):(\w+)`)
	if !re.MatchString(modelName) {
		modelName += ":latest"
	}

	localModels, err := Get[api.ListResponse]("api/tags")
	if err != nil {
		panic(err)
	}
	for _, model := range localModels.Models {
		if model.Name == modelName {
			return true
		}
	}
	return false
}

func streamedPrompt(model string, prompt string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	inputObj := api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}
	inputJson, err := json.Marshal(inputObj)
	if err != nil {
		log.Fatal("Error marshalling json: ", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(inputJson))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error occurred during request. Error: ", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	s.Stop()
	for {
		var sr api.GenerateResponse
		if err := decoder.Decode(&sr); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Error decoding JSON: ", err)
		}

		fmt.Printf(sr.Response)

		if sr.Done {
			break
		}
	}
	fmt.Println()
}

func ollamaUpCheck() {
	var quickClient = &http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := quickClient.Get("http://localhost:11434")
	if err != nil {
		panic("Ollama service is not running. Please ensure ollama it is installed and it's service running")
	}

	rawResp, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || err != nil || string(rawResp) != "Ollama is running" {
		panic("Ollama service is not running. Please ensure ollama it is installed and it's service running")
	}
}
