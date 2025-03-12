package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
	"testing"
)

// testWrapText is a separate implementation for testing purposes
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

func TestWrapText(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		width    int
		expected []string
	}{
		{
			name:     "Normal paragraph",
			text:     "This is a test of the word wrapping function. It should wrap at word boundaries.",
			width:    20,
			expected: []string{"This is a test of", "the word wrapping", "function. It should", "wrap at word", "boundaries."},
		},
		{
			name:     "Long word",
			text:     "This contains a verylongwordthatwillneedtobesplitacrossmultiplelines because it's too long.",
			width:    20,
			expected: []string{"This contains a", "verylongwordthatwil", "lneedtobesplitacros", "smultiplelines", "because it's too", "long."},
		},
		{
			name:     "Text with newlines",
			text:     "This has\na newline\nin it.",
			width:    20,
			expected: []string{"This has", "a newline", "in it."},
		},
		{
			name:  "Empty string",
			text:  "",
			width: 20,
			expected: []string{},
		},
		{
			name:     "Single character",
			text:     "x",
			width:    20,
			expected: []string{"x"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := wrapText(tc.text, tc.width)
			
			// Compare results
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d lines, got %d lines", len(tc.expected), len(result))
				return
			}
			
			for i := range result {
				if result[i] != tc.expected[i] {
					t.Errorf("Line %d: expected %q, got %q", i+1, tc.expected[i], result[i])
				}
			}
		})
	}
}

// TestWrapTextSpecial is a more manual test for debugging purposes
func TestWrapTextSpecial(t *testing.T) {
	// These tests just print the output for visual inspection
	testCases := []struct {
		text  string
		width int
	}{
		{"This is a test of the word wrapping function. It should wrap at word boundaries.", 20},
		{"This contains a verylongwordthatwillneedtobesplitacrossmultiplelines because it's too long.", 20},
		{"This has\na newline\nin it.", 20},
		{"", 20},
		{"x", 20},
	}

	for _, tc := range testCases {
		fmt.Printf("Original text: %q\n", tc.text)
		fmt.Printf("Wrapping at width %d:\n", tc.width)
		lines := wrapText(tc.text, tc.width)
		for i, line := range lines {
			fmt.Printf("Line %d: %q\n", i+1, line)
		}
		fmt.Println()
	}
}