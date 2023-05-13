package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	scanner := bufio.NewScanner(os.Stdin)
	var history = []openai.ChatCompletionMessage{}
	for {

		scanner.Scan()

		if err != nil {
			print(err.Error())
			return
		}
		item := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: scanner.Text(),
		}
		history = append(history, item)
		resp, done := call_completion(err, client, history)
		if done {
			return
		}
		item = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: resp,
		}
		history = append(history, item)
		fmt.Println(resp)
	}
}

func call_completion(err error, client *openai.Client, history []openai.ChatCompletionMessage) (string, bool) {

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: history,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return openai.ChatCompletionResponse{}.Choices[0].Message.Content, true
	}
	return resp.Choices[0].Message.Content, false
}
