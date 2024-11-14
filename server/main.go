package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	r := gin.Default()

	r.POST("", hello)

	r.Run() // listen and serve on 0.0.0.0:8080
}

type request struct {
	Text    string `json:"text"`
	Context []int  `json:"context"`
}

type response struct {
	Context  []int    `json:"context"`
	Thoughts string   `json:"thoughts"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

func hello(c *gin.Context) {
	var newChatMessage request

	if err := c.BindJSON(&newChatMessage); err != nil {
		fmt.Println(err)
		return
	}

	for {
		aiResponse, err := getResponse(newChatMessage)
		if err != nil {
			fmt.Println(err)
			continue
		}

		c.JSON(200, gin.H{"response": aiResponse})
		return
	}
}

func getResponse(newChatMessage request) (*response, error) {
	res, err := getAiResponse(newChatMessage.Text, newChatMessage.Context)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}

	log.Printf("response: %s", res["response"])

	aiResponse, err := parseAiResponse(res)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil, err
	}

	return aiResponse, nil
}

func parseAiResponse(data map[string]interface{}) (*response, error) {
	re := regexp.MustCompile("\\*\\*Thoughts\\*\\*|\\*\\*Question\\*\\*|\\*\\*Answers\\*\\*")
	split := re.Split(data["response"].(string), -1)

	if len(split) != 4 {
		return nil, errors.New("invalid response")
	}

	fmt.Println(len(split))

	rawContext := data["context"].([]interface{})
	thoughts := strings.TrimSpace(split[1])
	question := strings.TrimSpace(split[2])
	answersRaw := strings.TrimSpace(split[3])

	// Initialize an int slice to hold the converted values
	context := make([]int, len(rawContext))

	// Convert each element of []interface{} to int
	for i, v := range rawContext {
		// Assert that each value is a float64 (common for JSON numbers) or int
		switch v := v.(type) {
		case int:
			context[i] = v
		case float64:
			context[i] = int(v)
		default:
			fmt.Println("Error: non-numeric value in rawContext")
			return nil, errors.New("invalid rawContext")
		}
	}

	answers := strings.Split(answersRaw, "*")
	var filteredAnswers []string
	for _, answer := range answers {
		if strings.TrimSpace(answer) != "" {
			filteredAnswers = append(filteredAnswers, answer)
		}
	}

	return &response{Context: context, Thoughts: thoughts, Question: question, Answers: filteredAnswers}, nil
}

func getAiResponse(text string, context []int) (map[string]interface{}, error) {
	values := map[string]interface{}{
		"model": "llama3.2",
		"system": `
Your task is to help users find new and unique coding projects. Use your creativity!
You have as many Questions as you need to find a good project idee.
Once the user asks for a project idea you have to provide one!
simple put the project ideas as a answer to the question "what project idea do you like?"

start with broad question and than get more precise 

Provide answers in a specific structure, separated into three sections, with each section clearly separated by a Header.
The headers should be "**Thoughts**", "**Question**" and "**Answers**"
Do Not add any sort of Chapter Title do only use "**Thoughts**", "**Question**" and "**Answers**"

    First section: Provide your thoughts.
    Second section: Pose a question to the user to clarify their interests or objectives.
    Third section: List the possible answers to your question, each as a single line with a * as a element marker.
		`,
		"prompt":  text,
		"context": context,
		"stream":  false,
	}
	jsonData, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	return res, nil
}
