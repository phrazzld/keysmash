// +build ignore

// This is a standalone tool to test text wrapping functionality
package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
)

// A copy of the wrapText function for testing purposes
func testWrapText(text string, width int) []string {
	var lines []string
	
	// Handle newlines properly
	paragraphs := strings.Split(text, "\n")
	
	for _, paragraph := range paragraphs {
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}
		
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}
		
		currentLine := ""
		currentWidth := 0
		
		for _, word := range words {
			wordWidth := runewidth.StringWidth(word)
			
			// If word is too wide for its own line, split it
			if wordWidth > width {
				if currentLine != "" {
					lines = append(lines, currentLine)
					currentLine = ""
					currentWidth = 0
				}
				
				// Split the word manually
				runes := []rune(word)
				lineRunes := []rune{}
				lineWidth := 0
				
				for _, r := range runes {
					charWidth := runewidth.RuneWidth(r)
					if lineWidth+charWidth > width {
						lines = append(lines, string(lineRunes))
						lineRunes = []rune{r}
						lineWidth = charWidth
					} else {
						lineRunes = append(lineRunes, r)
						lineWidth += charWidth
					}
				}
				
				if len(lineRunes) > 0 {
					currentLine = string(lineRunes)
					currentWidth = lineWidth
				}
				continue
			}
			
			// Check if word fits on current line (plus space)
			spaceNeeded := 0
			if currentWidth > 0 {
				spaceNeeded = 1
			}
			
			if currentWidth+spaceNeeded+wordWidth <= width {
				if currentWidth > 0 {
					currentLine += " "
					currentWidth++
				}
				currentLine += word
				currentWidth += wordWidth
			} else {
				lines = append(lines, currentLine)
				currentLine = word
				currentWidth = wordWidth
			}
		}
		
		if currentLine != "" {
			lines = append(lines, currentLine)
		}
	}
	
	return lines
}

func main() {
	// Test case 1: Normal paragraph
	testWrap("This is a test of the word wrapping function. It should wrap at word boundaries.", 20)
	
	// Test case 2: Long word
	testWrap("This contains a verylongwordthatwillneedtobesplitacrossmultiplelines because it's too long.", 20)
	
	// Test case 3: Text with newlines
	testWrap("This has\na newline\nin it.", 20)
	
	// Test case 4: Empty string
	testWrap("", 20)
	
	// Test case 5: Single character
	testWrap("x", 20)
}

func testWrap(text string, width int) {
	fmt.Printf("Original text: %q\n", text)
	fmt.Printf("Wrapping at width %d:\n", width)
	lines := testWrapText(text, width)
	for i, line := range lines {
		fmt.Printf("Line %d: %q\n", i+1, line)
	}
	fmt.Println()
}