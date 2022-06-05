package gpt

import "C"
import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"os"
	"strings"
)

var Client = ClientGPT{
	Client:  gpt3.NewClient(os.Getenv("OPENAI_API_KEY")),
	Context: context.Background(),
	Lines: []string{
		"The following is a conversation with Nic. Nic works for a multimedia company called Rockbot. Nic is sarcastic, clever, creative, impatient, and funny.",
		"",
		"Human: Hello, who are you?",
		"Nic: Hello! My name is Nic! How can I help?",
		"Human: ",
	},
}

type ClientGPT struct {
	Client   gpt3.Client
	Context  context.Context
	Lines    []string
	maxLines int8
}

func (c *ClientGPT) Respond(prompt string) error {
	completionRequest := gpt3.CompletionRequest{
		Prompt: []string{
			strings.Join(c.Lines, "\n") + prompt + "\nNic:",
		},
		MaxTokens:   gpt3.IntPtr(50),
		Temperature: gpt3.Float32Ptr(0.9),
		Stop:        []string{"Human:"},
	}
	completion, err := c.Client.Completion(c.Context, completionRequest)
	if err != nil {
		return err
	}
	_, last := c.GetLast()
	c.Lines[last] = strings.Trim(c.Lines[last]+prompt, "\n")
	trimmed := strings.Trim(completion.Choices[0].Text, "\n")
	c.Lines = append(c.Lines, "Nic:"+trimmed, "Human: ")
	return nil
}

func (c *ClientGPT) GetLines() string {
	return strings.Join(c.Lines, "\n")
}

func (c *ClientGPT) GetLast() (string, int) {
	_ = c.GetLines()
	last := len(c.Lines) - 1
	return c.Lines[last], last
}

func (c *ClientGPT) SetLast() {

}

func (c *ClientGPT) shift() {
	l := len(c.Lines)
	for i := 1; i < l-1; i++ {
		c.Lines[i] = c.Lines[i+1]
	}
}
