package gpt

import "C"
import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"os"
	"strings"
)

var Client = ClientGPT{
	Client:  gpt3.NewClient(os.Getenv("OPENAI_API_KEY"), gpt3.WithDefaultEngine("curie")),
	Context: context.Background(),
	Lines: []string{
		"The following is a conversation with Nic.",
		"Nic works for a multimedia company called Rockbot. Nic is busy, creative, and funny.",
		"Nic likes to repeat what he says many times.",
		"Human: Hello, who are you?",
		"Nic: Hello! My name is Nic! How can I help?",
		"Human: ",
	},
	maxLines: 15,
}

type ClientGPT struct {
	Client   gpt3.Client
	Context  context.Context
	Lines    []string
	maxLines int
}

func (c *ClientGPT) Respond(prompt string) error {
	completionRequest := gpt3.CompletionRequest{
		Prompt: []string{
			strings.Join(c.Lines, "\n") + prompt + "\nNic:",
		},
		MaxTokens:        gpt3.IntPtr(50),
		Temperature:      gpt3.Float32Ptr(0.7),
		TopP:             gpt3.Float32Ptr(1),
		N:                gpt3.IntPtr(1),
		Echo:             false,
		Stop:             []string{"Human:"},
		PresencePenalty:  0.4,
		FrequencyPenalty: 0.5,
		Stream:           false,
	}
	completion, err := c.Client.Completion(c.Context, completionRequest)
	if err != nil {
		return err
	}
	_, last := c.GetLast()
	c.Lines[last] = strings.Trim(c.Lines[last]+prompt, "\n")
	trimmed := strings.Trim(completion.Choices[0].Text, "\n")
	if len(c.Lines) > c.maxLines {
		c.shift()
		c.shift()
		c.SetSecondToLast("Nic:" + trimmed)
		c.SetLast("Human: ")
	} else {
		c.Lines = append(c.Lines, "Nic:"+trimmed, "Human: ")
	}
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

func (c *ClientGPT) SetSecondToLast(s string) {
	slast := len(c.Lines) - 2
	c.Lines[slast] = s
}

func (c *ClientGPT) SetLast(s string) {
	last := len(c.Lines) - 1
	c.Lines[last] = s
}

func (c *ClientGPT) shift() {
	l := len(c.Lines)
	for i := 6; i < l-1; i++ {
		c.Lines[i] = c.Lines[i+1]
	}
}

func (c *ClientGPT) GetTrailing() string {
	sliceStart := len(c.Lines) - 2

	return strings.Join(c.Lines[sliceStart:], "\n")
}
