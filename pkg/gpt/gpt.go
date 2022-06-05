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
		"The following is a conversation with Nic and Mochi. ",
		"Nic works for a multimedia company called Rockbot. Nic is clever, creative, and funny.",
		"Mochi is a cat. She only purrs.",
		"Human: Hello, who are you?",
		"Nic: Hello! My name is Nic! How can I help?",
		"Mochi: *purrs*",
		"Human: ",
	},
	maxLines: 10,
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
	for i := 5; i < l-1; i++ {
		c.Lines[i] = c.Lines[i+1]
	}
}

func (c *ClientGPT) GetTrailing() string {
	sliceStart := len(c.Lines) - 2

	return strings.Join(c.Lines[sliceStart:], "\n")
}
