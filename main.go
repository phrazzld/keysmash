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
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorDefault)
	screen.SetStyle(defStyle)
	
	// Make defStyle available to other functions
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
		waitForKey(screen)

		// Select and load a test
		state, err := selectRandomTest()
		if err != nil {
			drawError(screen, fmt.Sprintf("Error loading test: %v", err))
			waitForKey(screen)
			return
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

	// Define colors
	titleStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorFuchsia).
		Bold(true)

	subtitleStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorAqua)
		
	defaultStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorDefault)

	// ASCII Art title
	title := []string{
		"╔══════════════════════════════╗",
		"║       NERV TYPING TEST       ║",
		"║    SYNCHRONIZATION SYSTEM    ║",
		"╚══════════════════════════════╝",
	}

	// Draw title centered
	for i, line := range title {
		drawText(screen, (width-runewidth.StringWidth(line))/2, height/4+i, titleStyle, line)
	}

	// Draw subtitle
	subtitle := "Initializing NERV Typing Simulation. Prepare for input."
	drawText(screen, (width-runewidth.StringWidth(subtitle))/2, height/4+len(title)+2, subtitleStyle, subtitle)

	// Draw instructions
	instructions := "Press any key to begin..."
	drawText(screen, (width-runewidth.StringWidth(instructions))/2, height/2+4, defaultStyle, instructions)

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
				if state.userInput == state.referenceText {
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
	
	// Define styles
	titleStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorFuchsia)
	
	correctStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorGreen)
	
	incorrectStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorRed)
	
	pendingStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorAqua)
	
	// Draw title
	drawText(screen, 2, 1, titleStyle, "NERV TYPING TEST - Synchronization in progress...")
	
	// Wrap reference text to fit screen width
	refLines := wrapText(state.referenceText, width-4)
	_ = wrapText(state.userInput, width-4) // Just to avoid unused variable
	
	// Draw reference text
	drawText(screen, 2, 3, titleStyle, "Reference Text:")
	for i, line := range refLines {
		drawText(screen, 2, 4+i, pendingStyle, line)
	}
	
	// Draw horizontal line
	for x := 0; x < width; x++ {
		screen.SetContent(x, 4+len(refLines)+1, tcell.RuneHLine, nil, titleStyle)
	}
	
	// Draw user input text with character-by-character styling
	drawText(screen, 2, 4+len(refLines)+3, titleStyle, "Your Input:")
	
	// Create a wrapped version of user input for display
	userInputLines := wrapText(state.userInput, width-4)
	inputY := 4 + len(refLines) + 4
	
	// Get character-by-character styling info first
	userInputRunes := []rune(state.userInput)
	refTextRunes := []rune(state.referenceText)
	
	// Create an array of styles for each character in the user input
	runeStyles := make([]tcell.Style, len(userInputRunes))
	for i, r := range userInputRunes {
		if i < len(refTextRunes) {
			if r == refTextRunes[i] {
				runeStyles[i] = correctStyle
			} else {
				runeStyles[i] = incorrectStyle
			}
		} else {
			runeStyles[i] = incorrectStyle // Extra characters
		}
	}
	
	// Now render each line with proper word wrapping
	runeOffset := 0 // Track which rune we're at in the original input
	
	for _, line := range userInputLines {
		inputX := 2
		
		// Process each character in the line with its pre-calculated style
		for _, r := range line {
			// Use the pre-calculated style for this rune
			if runeOffset < len(runeStyles) {
				style := runeStyles[runeOffset]
				screen.SetContent(inputX, inputY, r, nil, style)
			} else {
				// Fallback (shouldn't happen)
				screen.SetContent(inputX, inputY, r, nil, incorrectStyle)
			}
			
			inputX += runewidth.RuneWidth(r)
			runeOffset++
		}
		
		// Move to next line
		inputY++
	}
	
	// Draw cursor position (blinking underscore) at the end of the last line
	if time.Now().UnixNano()/1e8%10 < 5 {
		// If we have wrapped lines, put cursor at end of last line
		if len(userInputLines) > 0 {
			lastLine := userInputLines[len(userInputLines)-1]
			cursorX := 2 + runewidth.StringWidth(lastLine)
			cursorY := 4 + len(refLines) + 4 + len(userInputLines) - 1
			screen.SetContent(cursorX, cursorY, '_', nil, pendingStyle)
		} else {
			// No text yet, cursor at beginning
			screen.SetContent(2, 4 + len(refLines) + 4, '_', nil, pendingStyle)
		}
	}
	
	// Show timer and error count if test has started
	if state.testStarted {
		elapsed := time.Since(state.startTime).Seconds()
		timerText := fmt.Sprintf("Time: %.1fs | Errors: %d", elapsed, state.errors)
		drawText(screen, width-len(timerText)-2, 1, pendingStyle, timerText)
	}
	
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
	
	// Define styles
	titleStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorFuchsia).
		Bold(true)
		
	resultStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorGreen)
		
	defaultStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorDefault)
	
	// Calculate test metrics
	duration := state.endTime.Sub(state.startTime).Minutes()
	wpm := float64(len(state.referenceText)/5) / duration
	
	// Draw results frame
	frame := []string{
		"╔══════════════════════════════╗",
		"║     SIMULATION COMPLETE      ║",
		"║   ANALYZING PERFORMANCE...   ║",
		"╚══════════════════════════════╝",
	}
	
	for i, line := range frame {
		drawText(screen, (width-runewidth.StringWidth(line))/2, height/4+i, titleStyle, line)
	}
	
	// Draw results
	results := []string{
		fmt.Sprintf("WPM: %.1f", wpm),
		fmt.Sprintf("Errors: %d", state.errors),
		fmt.Sprintf("Time: %.1fs", state.endTime.Sub(state.startTime).Seconds()),
	}
	
	for i, line := range results {
		drawText(screen, (width-runewidth.StringWidth(line))/2, height/4+len(frame)+2+i, resultStyle, line)
	}
	
	// Draw options
	options := "Press [R] to retry, [N] for new test, [Q] to quit"
	drawText(screen, (width-runewidth.StringWidth(options))/2, height/2+6, defaultStyle, options)
	
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
	
	errorStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorRed).
		Bold(true)
		
	defaultStyle := tcell.StyleDefault.
		Background(tcell.ColorDefault).
		Foreground(tcell.ColorDefault)
		
	title := "ERROR"
	drawText(screen, (width-len(title))/2, height/3, errorStyle, title)
	
	drawText(screen, (width-len(message))/2, height/3+2, defaultStyle, message)
	
	instruction := "Press any key to exit"
	drawText(screen, (width-len(instruction))/2, height/3+4, defaultStyle, instruction)
	
	screen.Show()
}

func drawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func waitForKey(screen tcell.Screen) {
	for {
		ev := screen.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}