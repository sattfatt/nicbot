package main

import (
	"bufio"
	"fmt"
	"github.com/sattfatt/nicbot/pkg/gpt"
	"os"
)

func main() {

	// get input from user
	fmt.Print(gpt.Client.GetLines())
	for {
		last := gpt.Client.GetTrailing()
		fmt.Print(last)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		_ = gpt.Client.Respond(text)
	}
}
