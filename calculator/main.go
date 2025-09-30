package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Node struct {
	value string
	left  *Node
	right *Node
}

func tokenize(expression string) []string {
	var i int = 0

	var tokens []string
	var length int = len(expression)

	for i < length {
		c := expression[i]
		switch {
		case c >= '0' && c <= '9':
			{
				number := []byte{c}

				for i+1 < length && expression[i+1] >= '0' && expression[i+1] <= '9' {
					i++
					number = append(number, expression[i])
				}

				tokens = append(tokens, string(number))
			}
		case slices.Contains([]byte{'*', '/', '+', '-', '(', ')'}, c):
			{
				tokens = append(tokens, string(c))
			}
		}
		i++
	}

	return tokens
}

func isAddSub(s string) bool {
	if s == "+" || s == "-" {
		return true
	}
	return false
}

func isMulDiv(s string) bool {
	if s == "*" || s == "/" {
		return true
	}
	return false
}

func isValue(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func ConvertString(s string) (number float32) {
	number = 0
	for i := 0; i < len(s); i++ {
		number *= 10
		numberFromString := s[i]

		if numberFromString >= '0' && numberFromString <= '9' {
			number += float32(numberFromString - '0')
		} else {
			fmt.Printf("Invalid character: %c\n", numberFromString+'0')
			return 0
		}
	}
	return
}

func ParseTree(n *Node) string {
	if n == nil {
		return ""
	}

	if n.left == nil && n.right == nil {
		return fmt.Sprint(n.value)
	}
	return fmt.Sprintf("{ %s, %s, %s }", n.value, ParseTree(n.left), ParseTree(n.right))
}

func ParseAddSub(tokens []string) *Node {
	if len(tokens) == 1 && isValue(tokens[0]) {
		return &Node{value: tokens[0]}
	}

	var accumulated []string
	var parenthesesCounter int = 0
	for i := range tokens {
		if tokens[i] == "(" {
			parenthesesCounter++
		}
		if tokens[i] == ")" {
			parenthesesCounter--
		}
		if isAddSub(tokens[i]) && parenthesesCounter == 0 {
			return &Node{
				value: tokens[i],
				left:  ParseMulDiv(accumulated),
				right: ParseAddSub(tokens[i+1:]),
			}
		}
		accumulated = append(accumulated, tokens[i])
	}

	return ParseMulDiv(tokens)
}

func ParseMulDiv(tokens []string) *Node {
	if len(tokens) == 1 && isValue(tokens[0]) {
		return &Node{value: tokens[0]}
	}

	var accumulated []string
	var parenthesesCounter int = 0
	for i := range tokens {
		if tokens[i] == "(" {
			parenthesesCounter++
		}
		if tokens[i] == ")" {
			parenthesesCounter--
		}
		if isMulDiv(tokens[i]) && parenthesesCounter == 0 {
			return &Node{
				value: tokens[i],
				left:  ParseAddSub(accumulated),
				right: ParseAddSub(tokens[i+1:]),
			}
		}
		accumulated = append(accumulated, tokens[i])
	}

	return ParseParen(tokens)
}

func ParseParen(tokens []string) *Node {
	if tokens[0] == "(" && tokens[len(tokens)-1] == ")" {
		tokens = tokens[1 : len(tokens)-1]
	}

	return ParseAddSub(tokens)
}

func CalculateParseTree(node *Node) float32 {
	if node.value == "+" {
		return CalculateParseTree(node.left) + CalculateParseTree(node.right)
	}
	if node.value == "-" {
		return CalculateParseTree(node.left) - CalculateParseTree(node.right)
	}
	if node.value == "*" {
		return CalculateParseTree(node.left) * CalculateParseTree(node.right)
	}
	if node.value == "/" {
		return CalculateParseTree(node.left) / CalculateParseTree(node.right)
	}

	num := ConvertString(node.value)
	return num
}

func main() {
	fmt.Print(">> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	equation := scanner.Text()

	tokens := tokenize(equation)

	ast := ParseAddSub(tokens)
	tree := ParseTree(ast)
	fmt.Println("Parsed:", tree)

	fmt.Println("Result:", CalculateParseTree(ast))
}
