package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type TestState struct {
	referenceText string
	userInput     string
	errors        int
	startTime     time.Time
	endTime       time.Time
	testStarted   bool
	testComplete  bool
	testFile      string
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Initialize screen
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating screen: %v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	// Set default style
	defStyle := tcell.StyleDefault
	screen.SetStyle(defStyle)

	// Check if tests directory exists
	if _, err := os.Stat("tests"); os.IsNotExist(err) {
		drawError(screen, "Tests directory not found. Please create 'tests/' directory with text files.")
		waitForKey(screen)
		return
	}

	// Main application loop
	for {
		// Show welcome screen
		showWelcomeScreen(screen)
		if !waitForKey(screen) {
			// User pressed Escape, exit the program
			return
		}

		// Select and load a test
		state, err := selectRandomTest()
		if err != nil {
			drawError(screen, fmt.Sprintf("Error loading test: %v", err))
			if !waitForKey(screen) {
				return // User pressed Escape to quit
			}
			continue
		}

		// Run the typing test
		testResult := runTypingTest(screen, &state)

		// Handle post-test options
		if !handlePostTest(screen, testResult, &state) {
			break // User chose to quit
		}
	}
}

func showWelcomeScreen(screen tcell.Screen) {
	screen.Clear()
	width, height := screen.Size()

	// Draw basic welcome information
	title := "KEYSMASH"
	drawCenteredText(screen, width/2, height/2-3, tcell.StyleDefault, title)
	
	subtitle := "TYPING TEST"
	drawCenteredText(screen, width/2, height/2-1, tcell.StyleDefault, subtitle)
	
	prompt := "Press any key to start, ESC to quit"
	drawCenteredText(screen, width/2, height/2+3, tcell.StyleDefault, prompt)

	screen.Show()
}

func selectRandomTest() (TestState, error) {
	// Read test files
	files, err := os.ReadDir("tests")
	if err != nil {
		return TestState{}, err
	}

	// Filter for .txt files
	var textFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {
			textFiles = append(textFiles, file)
		}
	}

	if len(textFiles) == 0 {
		return TestState{}, fmt.Errorf("no .txt files found in tests/ directory")
	}

	// Select random file
	randomFile := textFiles[rand.Intn(len(textFiles))]

	// Read file content
	content, err := os.ReadFile("tests/" + randomFile.Name())
	if err != nil {
		return TestState{}, err
	}

	return TestState{
		referenceText: strings.TrimSpace(string(content)),
		userInput:     "",
		errors:        0,
		testStarted:   false,
		testComplete:  false,
		testFile:      randomFile.Name(),
	}, nil
}

func runTypingTest(screen tcell.Screen, state *TestState) TestState {
	width, _ := screen.Size()

	for {
		// Render current state
		renderScreen(screen, state, width)

		// Poll for events
		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
			width, _ = screen.Size()
		case *tcell.EventKey:
			// Handle key event
			if ev.Key() == tcell.KeyEscape {
				// Exit test
				return *state
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				// Handle backspace
				if len(state.userInput) > 0 {
					state.userInput = state.userInput[:len(state.userInput)-1]
				}
			} else if ev.Key() == tcell.KeyEnter {
				// Handle enter key - add a newline if reference text has one
				if len(state.userInput) < len(state.referenceText) &&
					len(state.referenceText) > len(state.userInput) &&
					state.referenceText[len(state.userInput)] == '\n' {
					state.userInput += "\n"
				}
			} else if r := ev.Rune(); r != 0 {
				// Handle character input
				if !state.testStarted {
					state.testStarted = true
					state.startTime = time.Now()
				}

				state.userInput += string(r)

				// Check for error
				if len(state.userInput) <= len(state.referenceText) {
					// Check if character matches
					if state.userInput[len(state.userInput)-1] != state.referenceText[len(state.userInput)-1] {
						state.errors++
					}
				} else {
					// Extra character is an error
					state.errors++
				}

				// Check if test is complete
				if len(state.userInput) == len(state.referenceText) && state.userInput == state.referenceText {
					state.testComplete = true
					state.endTime = time.Now()
					return *state
				}
			}
		}
	}
}

func renderScreen(screen tcell.Screen, state *TestState, width int) {
	screen.Clear()

	// Get screen dimensions
	width, screenHeight := screen.Size()
	
	// Set horizontal padding
	hPadding := 4
	
	// Calculate content width for wrapping
	contentWidth := width - (hPadding * 2)
	
	// Draw header with more padding
	headerText := "KEYSMASH - TYPING TEST"
	drawCenteredText(screen, width/2, 1, tcell.StyleDefault, headerText)
	
	// Show file name
	sourceText := fmt.Sprintf("Source: %s", state.testFile)
	drawCenteredText(screen, width/2, 3, tcell.StyleDefault, sourceText)
	
	// Draw stats if test started
	if state.testStarted {
		elapsed := time.Since(state.startTime).Seconds()
		
		// Calculate stats
		wpm := float64(len(state.userInput)/5) / (elapsed / 60.0)
		if wpm < 0 || elapsed < 1 {
			wpm = 0
		}
		
		// Display stats
		statsText := fmt.Sprintf("Time: %.1fs | WPM: %.1f | Errors: %d", 
			elapsed, wpm, state.errors)
		drawCenteredText(screen, width/2, 5, tcell.StyleDefault, statsText)
		
		// Display progress percentage
		completionPct := float64(len(state.userInput)) / float64(len(state.referenceText))
		if completionPct > 1.0 {
			completionPct = 1.0
		}
		
		pctText := fmt.Sprintf("Progress: %d%%", int(completionPct*100))
		drawText(screen, hPadding, 7, tcell.StyleDefault, pctText)
	}
	
	// Draw divider with more padding
	drawText(screen, 0, 9, tcell.StyleDefault, strings.Repeat("-", width))

	// Determine maximum content height to prevent overflow
	maxContentHeight := screenHeight - 18 // More space for headers and footers
	
	// Draw reference text title with more padding
	drawText(screen, hPadding, 11, tcell.StyleDefault, "Text to type:")
	
	// Wrap and draw reference text
	refLines := wrapText(state.referenceText, contentWidth)
	
	// Limit to half of available space
	maxRefLines := maxContentHeight / 2
	if len(refLines) > maxRefLines {
		refLines = refLines[:maxRefLines]
	}
	
	// Draw reference text with more padding between lines
	for i, line := range refLines {
		drawText(screen, hPadding, 13+i, tcell.StyleDefault, line)
	}
	
	// Draw separator between reference and input with more space
	separatorY := 13 + len(refLines) + 2
	drawText(screen, 0, separatorY, tcell.StyleDefault, strings.Repeat("-", width))
	
	// Draw input area label with more padding
	inputLabelY := separatorY + 2
	drawText(screen, hPadding, inputLabelY, tcell.StyleDefault, "Your typing:")
	
	// Draw user input with more padding
	userInputY := inputLabelY + 2
	
	// Fixed cursor position calculation to handle multibyte characters correctly
	cursorPos := 0
	cursorLine := 0
	
	if len(state.userInput) > 0 {
		// Wrap user input for display
		inputLines := wrapText(state.userInput, contentWidth)
		
		// Draw user input lines
		for i, line := range inputLines {
			drawText(screen, hPadding, userInputY+i, tcell.StyleDefault, line)
		}
		
		// Calculate cursor position
		if len(inputLines) > 0 {
			// Last line length gives cursor position
			lastLine := inputLines[len(inputLines)-1]
			cursorPos = runewidth.StringWidth(lastLine)
			cursorLine = len(inputLines) - 1
		}
	}
	
	// Draw blinking cursor at end of input
	cursorX := hPadding + cursorPos
	cursorY := userInputY + cursorLine
	
	if time.Now().UnixNano()/4e7%10 >= 5 {
		screen.SetContent(cursorX, cursorY, ' ', nil, tcell.StyleDefault.Reverse(true))
	} else {
		screen.SetContent(cursorX, cursorY, '_', nil, tcell.StyleDefault)
	}
	
	// Draw progress indicator with more padding
	progressBarY := screenHeight - 3
	progress := 0
	if len(state.referenceText) > 0 {
		progress = len(state.userInput) * 100 / len(state.referenceText)
	}
	
	// Create a wider progress bar
	progressBarWidth := 60
	filledWidth := progressBarWidth * progress / 100
	
	progressBar := fmt.Sprintf("[%s%s] %d%%", 
		strings.Repeat("=", filledWidth), 
		strings.Repeat(" ", progressBarWidth-filledWidth),
		progress)
	drawText(screen, hPadding, progressBarY, tcell.StyleDefault, progressBar)
	
	// Draw help text
	drawText(screen, hPadding, screenHeight-1, tcell.StyleDefault, "ESC to quit")

	screen.Show()
}

func wrapText(text string, width int) []string {
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

func handlePostTest(screen tcell.Screen, state TestState, originalState *TestState) bool {
	if !state.testComplete {
		return true // Test was interrupted, continue with a new test
	}

	screen.Clear()
	width, height := screen.Size()
	
	// Calculate test metrics
	duration := state.endTime.Sub(state.startTime).Minutes()
	wpm := float64(len(state.referenceText)/5) / duration
	accuracy := 100.0
	if len(state.userInput) > 0 {
		accuracy = 100.0 * (1.0 - float64(state.errors)/float64(len(state.userInput)))
	}
	if accuracy < 0 {
		accuracy = 0
	}
	
	// Display results with more spacing
	drawCenteredText(screen, width/2, height/2-8, tcell.StyleDefault, "TEST COMPLETE")
	
	// Show source
	drawCenteredText(screen, width/2, height/2-6, tcell.StyleDefault, fmt.Sprintf("Source: %s", state.testFile))
	
	// Draw results with more spacing
	drawCenteredText(screen, width/2, height/2-3, tcell.StyleDefault, fmt.Sprintf("WPM: %.1f", wpm))
	drawCenteredText(screen, width/2, height/2-1, tcell.StyleDefault, fmt.Sprintf("Accuracy: %.1f%%", accuracy))
	drawCenteredText(screen, width/2, height/2+1, tcell.StyleDefault, fmt.Sprintf("Time: %.1fs", state.endTime.Sub(state.startTime).Seconds()))
	drawCenteredText(screen, width/2, height/2+3, tcell.StyleDefault, fmt.Sprintf("Characters: %d (Errors: %d)", len(state.userInput), state.errors))
	
	// Draw options with more spacing
	drawCenteredText(screen, width/2, height/2+6, tcell.StyleDefault, "R: Retry  N: New Test  Q: Quit")
	
	screen.Show()
	
	// Wait for user choice
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch unicode := ev.Rune(); unicode {
				case 'R', 'r':
					// Retry the same test
					originalState.userInput = ""
					originalState.errors = 0
					originalState.testStarted = false
					originalState.testComplete = false
					return true
				case 'N', 'n':
					// New test
					return true
				case 'Q', 'q':
					// Quit
					return false
				}
			case tcell.KeyEscape:
				return false
			}
		}
	}
}

func drawError(screen tcell.Screen, message string) {
	screen.Clear()
	width, height := screen.Size()
	
	// Display error message with more spacing
	drawCenteredText(screen, width/2, height/2-4, tcell.StyleDefault, "ERROR")
	drawCenteredText(screen, width/2, height/2, tcell.StyleDefault, message)
	drawCenteredText(screen, width/2, height/2+4, tcell.StyleDefault, "Press any key to retry, ESC to quit")
	
	screen.Show()
}

// Helper function to draw text at a specific position
func drawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

// Helper function to draw centered text
func drawCenteredText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	textWidth := runewidth.StringWidth(text)
	startX := x - textWidth/2
	drawText(screen, startX, y, style, text)
}

func waitForKey(screen tcell.Screen) bool {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Return false if user pressed Escape (indicating quit)
			if ev.Key() == tcell.KeyEscape {
				return false
			}
			return true // Any other key continues
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}