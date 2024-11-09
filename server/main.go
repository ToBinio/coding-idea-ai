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

	values := map[string]interface{}{
		"model":   "llama3.2",
		"system":  "your job is to give short answers. Keep your answers as short as possible!",
		"prompt":  newChatMessage.Text,
		"context": newChatMessage.Context,
		"stream":  false,
	}
	jsonData, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["response"])

	c.JSON(200, gin.H{"response": res["response"], "context": res["context"]})
}
