package stringunpack

import (
	"fmt"
	"strings"
)

const (
	symTypeUndefined = 0
	symTypeChar      = 1
	symTypeNum       = 2
	symTypeBackslash = 3
	symZero          = 48
	symNine          = symZero + 9
)

// Unpack convert symbols*number to symbols
func Unpack(input string) (string, error) {
	var builder strings.Builder
	var prevSymbol, nextPrintSymbol rune
	prevSymbolType := symTypeUndefined
	repeatSymCount := 0
	// adding fake symbol at the end
	input += string(0)
	fmt.Printf("Processing string: %s \n\n", input)
	for currentPosition, currentSymbol := range input {
		fmt.Printf("Processing symbol: %c - %d\n", currentSymbol, currentSymbol)
		switch {
		case currentSymbol >= symZero && currentSymbol <= symNine:
			// this is num
			switch prevSymbolType {
			case symTypeUndefined:
				return "", fmt.Errorf("No previous symbol before number")
			case symTypeChar, symTypeNum:
				repeatSymCount = repeatSymCount*10 + int(currentSymbol) - symZero
				fmt.Printf("Current repeat time: %d\n", repeatSymCount)
				prevSymbolType = symTypeNum
			case symTypeBackslash:
				// this is not int - this is number, interpreted as char
				nextPrintSymbol = currentSymbol
				prevSymbolType = symTypeChar
			default:
				return "", fmt.Errorf("Unknown previous symbol: %d", currentSymbol)
			}
		case currentSymbol == '\\':
			switch prevSymbolType {
			case symTypeNum:
				for i := 0; i < repeatSymCount; i++ {
					builder.WriteRune(nextPrintSymbol)
				}
				repeatSymCount = 0
				prevSymbolType = symTypeBackslash
			case symTypeChar:
				builder.WriteRune(prevSymbol)
				prevSymbolType = symTypeBackslash
			case symTypeBackslash:
				nextPrintSymbol = currentSymbol
				prevSymbolType = symTypeChar
			}
		default:
			// this is char
			switch prevSymbolType {
			case symTypeUndefined:
			case symTypeChar:
				builder.WriteRune(prevSymbol)
			case symTypeNum:
				for i := 0; i < repeatSymCount; i++ {
					builder.WriteRune(nextPrintSymbol)
				}
				repeatSymCount = 0
			case symTypeBackslash:
				return "", fmt.Errorf("Backslash before char: %c at position: %d", currentSymbol, currentPosition)
			default:
				return "", fmt.Errorf("Unknown previous symbol: %d at position: %d", currentSymbol, currentPosition)
			}
			prevSymbolType = symTypeChar
			nextPrintSymbol = currentSymbol
		}
		prevSymbol = currentSymbol
	}
	result := builder.String()
	return result, nil
}
