package main

import (
	"bufio"
	"fmt"
	"github.com/sattfatt/nicbot/pkg/gpt"
	"os"
)

func main() {

	// get input from user
	for {
		fmt.Print(gpt.Client.GetLines())
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		_ = gpt.Client.Respond(text)
	}
}
