package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var httpClient = &http.Client{}

type ErrorType struct {
	Error string `json:"error"`
}

func Post[T interface{}, U interface{}](path string, inputData T) (U, error) {
	url := fmt.Sprintf("http://localhost:11434/%s", path)
	inputJson, err := json.Marshal(inputData)
	if err != nil {
		panic(fmt.Sprintf("Error marshalling json: %s", err))
	}
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(inputJson))
	if err != nil {
		log.Fatalf("Error occurred during request. Error: %s", err.Error())
	}
	return decodeHttpResponse[U](resp)
}

func Get[T interface{}](path string) (T, error) {
	url := fmt.Sprintf("http://localhost:11434/%s", path)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Fatalf("Error occurred during request. Error: %s", err.Error())
	}
	return decodeHttpResponse[T](resp)

}

func decodeHttpResponse[T interface{}](resp *http.Response) (T, error) {
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error occurred during reading the response Body. Error: %s", err.Error())
	}

	// Try to get proper response
	var response T
	err = json.Unmarshal(body, &response)
	if err != nil {
		var errorResponse ErrorType
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			log.Fatalf("Error occurred during unmarshalling the response. Error: %s", err.Error())
		}
		return response, errors.New(errorResponse.Error)
	}
	return response, nil
}
