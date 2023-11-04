package main

import (
	"regexp"

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
