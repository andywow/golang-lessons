package stringunpack

import (
	"fmt"
	"strings"
)

const (
	symTypeUndefined = iota
	symTypeChar
	symTypeNum
	symTypeBackslash
)

const (
	symZero = 48
	symNine = symZero + 9
)

var (
	builder                                         strings.Builder
	currentSymbol, nextPrintSymbol, prevSymbol      rune
	currentPosition, prevSymbolType, repeatSymCount int
)

func processSymbol() error {
	fmt.Printf("Processing symbol: %c - %d\n", currentSymbol, currentSymbol)
	switch {
	case currentSymbol >= symZero && currentSymbol <= symNine:
		// this is num
		switch prevSymbolType {
		case symTypeUndefined:
			return fmt.Errorf("No previous symbol before number")
		case symTypeChar, symTypeNum:
			repeatSymCount = repeatSymCount*10 + int(currentSymbol) - symZero
			fmt.Printf("Current repeat time: %d\n", repeatSymCount)
			prevSymbolType = symTypeNum
		case symTypeBackslash:
			// this is not int - this is number, interpreted as char
			nextPrintSymbol = currentSymbol
			prevSymbolType = symTypeChar
		default:
			return fmt.Errorf("Unknown previous symbol: %d", currentSymbol)
		}
	case currentSymbol == '\\':
		switch prevSymbolType {
		case symTypeUndefined:
			prevSymbolType = symTypeBackslash
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
			return fmt.Errorf("Backslash before char: %c at position: %d", currentSymbol, currentPosition)
		default:
			return fmt.Errorf("Unknown previous symbol: %d at position: %d", currentSymbol, currentPosition)
		}
		prevSymbolType = symTypeChar
		nextPrintSymbol = currentSymbol
	}
	prevSymbol = currentSymbol
	return nil
}

// Unpack convert symbols*number to symbols
func Unpack(input string) (string, error) {
	builder.Reset()
	prevSymbolType = symTypeUndefined
	repeatSymCount = 0
	// adding fake symbol at the end
	fmt.Printf("Processing string: %s \n\n", input)
	for currentPosition, currentSymbol = range input {
		if err := processSymbol(); err != nil {
			return "", err
		}
	}
	currentSymbol = 0
	processSymbol()
	return builder.String(), nil
}
