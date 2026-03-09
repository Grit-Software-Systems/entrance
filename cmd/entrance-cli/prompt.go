package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const (
	backspaceCharacter byte = 8
	deleteCharacter    byte = 127
	maskingCharacter        = "*"
)

func readChoice(label string, minimumChoice int, maximumChoice int) int {
	var reader *bufio.Reader = bufio.NewReader(os.Stdin)
	for {
		fmt.Print(label)
		var line string
		line, _ = reader.ReadString('\n')
		var trimmedLine string = strings.TrimSpace(line)
		var parsedValue int
		var parseError error
		parsedValue, parseError = strconv.Atoi(trimmedLine)
		if parseError != nil {
			fmt.Printf("Please enter a number between %d and %d.\n", minimumChoice, maximumChoice)
			continue
		}
		if parsedValue < minimumChoice || parsedValue > maximumChoice {
			fmt.Printf("Please enter a number between %d and %d.\n", minimumChoice, maximumChoice)
			continue
		}
		result := parsedValue
		return result
	}
}

func readLine(label string) string {
	var reader *bufio.Reader = bufio.NewReader(os.Stdin)
	fmt.Print(label)
	var line string
	line, _ = reader.ReadString('\n')
	result := strings.TrimSpace(line)
	return result
}

func readPassword(label string) string {
	fmt.Print(label)
	var fileDescriptor int = int(os.Stdin.Fd())
	var previousState *term.State
	var makeRawError error
	previousState, makeRawError = term.MakeRaw(fileDescriptor)
	if makeRawError != nil {
		fmt.Println()
		fmt.Println("Error reading password.")
		result := ""
		return result
	}
	defer term.Restore(fileDescriptor, previousState)
	var passwordBytes []byte
	for {
		var singleByte [1]byte
		var readCount int
		var readError error
		readCount, readError = os.Stdin.Read(singleByte[:])
		if readError != nil || readCount == 0 {
			break
		}
		var inputByte byte = singleByte[0]
		if inputByte == '\r' || inputByte == '\n' {
			fmt.Print("\r\n")
			break
		}
		if inputByte == backspaceCharacter || inputByte == deleteCharacter {
			if len(passwordBytes) > 0 {
				passwordBytes = passwordBytes[:len(passwordBytes)-1]
				fmt.Print("\b \b")
			}
			continue
		}
		passwordBytes = append(passwordBytes, inputByte)
		fmt.Print(maskingCharacter)
	}
	var passwordString string = string(passwordBytes)
	result := strings.TrimSpace(passwordString)
	return result
}
