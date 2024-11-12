package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	r.POST("", hello)

	r.Run() // listen and serve on 0.0.0.0:8080
}

type chatMessage struct {
	Text    string `json:"text"`
	Context []int  `json:"context"`
}

func hello(c *gin.Context) {
	var newChatMessage chatMessage

	if err := c.BindJSON(&newChatMessage); err != nil {
		fmt.Println(err)
		return
	}

	res, err := getAiResponse(newChatMessage.Text, newChatMessage.Context)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}

	log.Printf("response: %s", res["response"])

	c.JSON(200, gin.H{"response": res["response"], "context": res["context"]})
}

func getAiResponse(text string, context []int) (map[string]interface{}, error) {
	values := map[string]interface{}{
		"model": "llama3.2",
		"system": `
Your task is to help users find new and unique coding projects. Use your creativity!
You have 10 Questions to find a good project idee.
Start with very broad questions and get more precise as you go. 
After 10 Questions the 11 Question should be what project idee the user wants to choose.

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
