package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	msg := make(chan string)

	go func() {
		for {
			message := <-msg
			fmt.Println(message)
		}
	}()

	for {
		var input strings.Builder
		for {
			scanner.Scan()
			line := scanner.Text()
			if len(line) == 0 {
				break
			}
			input.WriteString(line)
			input.WriteString("\n")
		}
		msg <- input.String()
	}
}
