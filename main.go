package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Global variable to store the path to the tests directory
var testsDir string

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

// findTestsDir tries to locate the tests directory in various locations
func findTestsDir() string {
	// Try current directory first
	if _, err := os.Stat("tests"); err == nil {
		return "tests"
	}

	// Try executable directory
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		testsInExecDir := filepath.Join(execDir, "tests")
		if _, err := os.Stat(testsInExecDir); err == nil {
			return testsInExecDir
		}
		
		// Check one level up (for GOPATH/bin scenario)
		parentDir := filepath.Dir(execDir)
		testsInParentDir := filepath.Join(parentDir, "tests")
		if _, err := os.Stat(testsInParentDir); err == nil {
			return testsInParentDir
		}
	}
	
	// Try the source directory where keysmash was built
	sourceDir := "/Users/phaedrus/Development/keysmash/tests"
	if _, err := os.Stat(sourceDir); err == nil {
		return sourceDir
	}

	// Not found
	return ""
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

	// Find tests directory
	testsDir = findTestsDir()
	if testsDir == "" {
		drawError(screen, "Tests directory not found. Please create a 'tests' directory with text files.")
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
	// Read test files from the identified tests directory
	files, err := os.ReadDir(testsDir)
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
		return TestState{}, fmt.Errorf("no .txt files found in %s directory", testsDir)
	}

	// Select random file
	randomFile := textFiles[rand.Intn(len(textFiles))]

	// Read file content using the full path
	content, err := os.ReadFile(filepath.Join(testsDir, randomFile.Name()))
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
				// Always allow Enter key to add a newline
				if !state.testStarted {
					state.testStarted = true
					state.startTime = time.Now()
				}
				
				// Add the newline
				state.userInput += "\n"
				
				// Check if the newline matches the reference text
				if len(state.userInput) <= len(state.referenceText) {
					if state.referenceText[len(state.userInput)-1] != '\n' {
						state.errors++
					}
				} else {
					// Extra character is an error
					state.errors++
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

// renderScreen handles the UI drawing with adaptive layout
func renderScreen(screen tcell.Screen, state *TestState, width int) {
	screen.Clear()

	// Get screen dimensions
	width, screenHeight := screen.Size()
	
	// Check for minimum screen size
	minWidth := 40
	minHeight := 15
	
	if width < minWidth || screenHeight < minHeight {
		// Screen is too small, render minimal UI with error message
		renderMinimalScreen(screen, state, width, screenHeight)
		return
	}
	
	// Set horizontal padding (adaptive based on screen width)
	hPadding := min(4, width/10)
	
	// Calculate content width for wrapping
	contentWidth := max(20, width - (hPadding * 2))
	
	// Draw header (adaptive based on space)
	if screenHeight >= 18 {
		headerText := "KEYSMASH - TYPING TEST"
		drawCenteredText(screen, width/2, 1, tcell.StyleDefault, headerText)
		
		// Show file name
		sourceText := fmt.Sprintf("Source: %s", state.testFile)
		drawCenteredText(screen, width/2, 3, tcell.StyleDefault, sourceText)
	} else {
		// For smaller screens, just show a compact header
		headerText := "KEYSMASH"
		drawCenteredText(screen, width/2, 0, tcell.StyleDefault, headerText)
	}
	
	// Wrap all text first
	refLines := wrapText(state.referenceText, contentWidth)
	inputLines := []string{}
	if len(state.userInput) > 0 {
		inputLines = wrapText(state.userInput, contentWidth)
	}
	
	// Calculate cursor position
	cursorPos := 0
	cursorLine := 0
	if len(inputLines) > 0 {
		lastLine := inputLines[len(inputLines)-1]
		cursorPos = runewidth.StringWidth(lastLine)
		cursorLine = len(inputLines) - 1
	}
	
	// Calculate dynamic UI layout
	var topMargin, statsHeight, refHeaderHeight, refSectionHeight int
	var inputHeaderHeight, inputSectionHeight, bottomMargin int
	
	// Adaptive layout based on screen size
	if screenHeight >= 24 {
		// Full featured layout for large screens
		topMargin = 4
		statsHeight = 3
		refHeaderHeight = 2
		bottomMargin = 3
		inputHeaderHeight = 2
	} else if screenHeight >= 18 {
		// Medium layout
		topMargin = 2
		statsHeight = 2
		refHeaderHeight = 1
		bottomMargin = 2
		inputHeaderHeight = 1
	} else {
		// Minimal layout
		topMargin = 1
		statsHeight = 1
		refHeaderHeight = 1
		bottomMargin = 2
		inputHeaderHeight = 1
	}
	
	// Draw stats if test started
	statsY := topMargin
	if state.testStarted {
		elapsed := time.Since(state.startTime).Seconds()
		
		// Calculate stats
		wpm := float64(len(state.userInput)/5) / (elapsed / 60.0)
		if wpm < 0 || elapsed < 1 {
			wpm = 0
		}
		
		// Display stats (adaptive based on space)
		if screenHeight >= 18 {
			statsText := fmt.Sprintf("Time: %.1fs | WPM: %.1f | Errors: %d", 
				elapsed, wpm, state.errors)
			drawCenteredText(screen, width/2, statsY, tcell.StyleDefault, statsText)
			
			// Display progress percentage
			completionPct := float64(len(state.userInput)) / float64(len(state.referenceText))
			if completionPct > 1.0 {
				completionPct = 1.0
			}
			
			pctText := fmt.Sprintf("Progress: %d%%", int(completionPct*100))
			drawText(screen, hPadding, statsY+1, tcell.StyleDefault, pctText)
		} else {
			// Compact stats for smaller screens
			statsText := fmt.Sprintf("WPM: %.1f | Err: %d", wpm, state.errors)
			drawCenteredText(screen, width/2, statsY, tcell.StyleDefault, statsText)
		}
	}
	
	// Calculate main content area boundaries
	contentStartY := topMargin + statsHeight + 1
	contentEndY := screenHeight - bottomMargin
	contentHeight := contentEndY - contentStartY
	
	// Safety check - ensure we have minimum content space
	if contentHeight < 4 {
		// Screen is too small, render minimal UI with error message
		renderMinimalScreen(screen, state, width, screenHeight)
		return
	}
	
	// Dynamic space allocation - reference gets 1/3, input gets 2/3
	// but ensure at least 2 lines for each section
	refSectionHeight = max(2, contentHeight / 3)
	inputSectionHeight = max(2, contentHeight - refSectionHeight - refHeaderHeight - inputHeaderHeight - 1) // -1 for separator
	
	// Reference text section
	refTextTitleY := contentStartY
	refTextStartY := refTextTitleY + refHeaderHeight
	
	// Draw divider between stats and content
	drawText(screen, 0, contentStartY-1, tcell.StyleDefault, strings.Repeat("-", width))
	
	// Draw reference text title
	drawText(screen, hPadding, refTextTitleY, tcell.StyleDefault, "Text to type:")
	
	// Ensure we have at least one line to display reference text
	if refSectionHeight > 0 {
		// Handle case when reference text is longer than available space
		if len(refLines) > refSectionHeight {
			// Calculate which portion to display based on typing progress
			refProgress := 0.0
			if len(state.referenceText) > 0 {
				refProgress = float64(len(state.userInput)) / float64(len(state.referenceText))
			}
			refMidpoint := int(refProgress * float64(len(refLines)))
			
			// Calculate start/end lines with bounds checking
			refStartLine := max(0, refMidpoint-(refSectionHeight/2))
			refEndLine := min(len(refLines), refStartLine+refSectionHeight)
			
			// Adjust if we're near the end
			if refEndLine >= len(refLines) {
				refStartLine = max(0, len(refLines)-refSectionHeight)
				refEndLine = len(refLines)
			}
			
			// Safety check for array bounds
			if refStartLine < refEndLine && refStartLine >= 0 && refEndLine <= len(refLines) {
				// Draw only the visible portion
				for i, line := range refLines[refStartLine:refEndLine] {
					drawText(screen, hPadding, refTextStartY+i, tcell.StyleDefault, line)
				}
				
				// Add scroll indicators if needed (if we have room)
				if refStartLine > 0 && width > 20 {
					drawText(screen, width-6, refTextStartY, tcell.StyleDefault, "↑")
				}
				if refEndLine < len(refLines) && width > 20 {
					drawText(screen, width-6, refTextStartY+refSectionHeight-1, tcell.StyleDefault, "↓")
				}
			}
		} else if len(refLines) > 0 {
			// Draw all reference text if it fits
			for i, line := range refLines {
				if i < refSectionHeight { // Bounds check
					drawText(screen, hPadding, refTextStartY+i, tcell.StyleDefault, line)
				}
			}
		}
	}
	
	// Calculate input section position
	separatorY := refTextStartY + refSectionHeight
	inputLabelY := separatorY + 1
	inputStartY := inputLabelY + inputHeaderHeight
	
	// Draw separator between reference and input
	if separatorY < screenHeight-1 {
		drawText(screen, 0, separatorY, tcell.StyleDefault, strings.Repeat("-", width))
	}
	
	// Draw input area label
	if inputLabelY < screenHeight-1 {
		drawText(screen, hPadding, inputLabelY, tcell.StyleDefault, "Your typing:")
	}
	
	// Draw user input if we have space
	if inputSectionHeight > 0 && inputStartY < screenHeight-1 {
		if len(inputLines) > 0 {
			// Calculate how many lines we can display
			inputStartLine := 0
			
			// If cursor would be beyond visible area, scroll to show it
			if cursorLine >= inputSectionHeight {
				// Keep cursor a few lines from the bottom for context
				inputStartLine = max(0, cursorLine-(inputSectionHeight-1))
			}
			
			// Calculate the end line (capped by available lines or content)
			inputEndLine := min(len(inputLines), inputStartLine+inputSectionHeight)
			
			// Safety check for array bounds
			if inputStartLine < inputEndLine && inputStartLine >= 0 && inputEndLine <= len(inputLines) {
				// Draw visible input lines
				for i, line := range inputLines[inputStartLine:inputEndLine] {
					if inputStartY+i < screenHeight-1 { // Bounds check
						drawText(screen, hPadding, inputStartY+i, tcell.StyleDefault, line)
					}
				}
				
				// Add scroll indicators if needed (if we have room)
				if inputStartLine > 0 && width > 20 {
					drawText(screen, width-6, inputStartY, tcell.StyleDefault, "↑")
				}
				if inputEndLine < len(inputLines) && width > 20 && inputStartY+inputSectionHeight-1 < screenHeight-1 {
					drawText(screen, width-6, inputStartY+inputSectionHeight-1, tcell.StyleDefault, "↓")
				}
			}
			
			// Position cursor (with bounds checking)
			if cursorLine >= inputStartLine {
				cursorY := inputStartY + (cursorLine - inputStartLine)
				cursorX := hPadding + cursorPos
				
				if cursorX < width && cursorY < screenHeight-1 {
					// Draw blinking cursor at end of input
					if time.Now().UnixNano()/4e7%10 >= 5 {
						screen.SetContent(cursorX, cursorY, ' ', nil, tcell.StyleDefault.Reverse(true))
					} else {
						screen.SetContent(cursorX, cursorY, '_', nil, tcell.StyleDefault)
					}
				}
			}
		} else {
			// No input yet, just show cursor at start position
			cursorX := hPadding
			cursorY := inputStartY
			
			if cursorX < width && cursorY < screenHeight-1 {
				if time.Now().UnixNano()/4e7%10 >= 5 {
					screen.SetContent(cursorX, cursorY, ' ', nil, tcell.StyleDefault.Reverse(true))
				} else {
					screen.SetContent(cursorX, cursorY, '_', nil, tcell.StyleDefault)
				}
			}
		}
	}
	
	// Draw progress bar at bottom
	progressBarY := screenHeight - 2
	if progressBarY > 0 {
		progress := 0
		if len(state.referenceText) > 0 {
			progress = len(state.userInput) * 100 / len(state.referenceText)
		}
		
		// Adaptive progress bar width
		progressBarWidth := min(60, width - (2 * hPadding))
		if progressBarWidth < 10 {
			// Just show percentage for very narrow screens
			progressText := fmt.Sprintf("%d%%", progress)
			drawCenteredText(screen, width/2, progressBarY, tcell.StyleDefault, progressText)
		} else {
			// Draw progress bar
			filledWidth := progressBarWidth * progress / 100
			
			progressBar := fmt.Sprintf("[%s%s] %d%%", 
				strings.Repeat("=", filledWidth), 
				strings.Repeat(" ", progressBarWidth-filledWidth),
				progress)
			drawText(screen, hPadding, progressBarY, tcell.StyleDefault, progressBar)
		}
		
		// Draw help text at very bottom
		if screenHeight > 2 {
			drawText(screen, hPadding, screenHeight-1, tcell.StyleDefault, "ESC to quit")
		}
	}

	screen.Show()
}

// renderMinimalScreen is a simplified UI for very small terminal windows
func renderMinimalScreen(screen tcell.Screen, state *TestState, width, height int) {
	// Show compact header
	if height > 0 {
		title := "KEYSMASH"
		if width < len(title) {
			title = title[:width]
		}
		drawCenteredText(screen, width/2, 0, tcell.StyleDefault, title)
	}
	
	// Show an error message about screen size
	if height > 2 && width > 15 {
		msg := "Window too small"
		drawCenteredText(screen, width/2, 2, tcell.StyleDefault, msg)
	}
	
	// Show minimal stats if we have room
	if height > 4 && state.testStarted {
		elapsed := time.Since(state.startTime).Seconds()
		wpm := float64(len(state.userInput)/5) / (elapsed / 60.0)
		if wpm < 0 || elapsed < 1 {
			wpm = 0
		}
		
		statsText := fmt.Sprintf("WPM:%.1f", wpm)
		if width > len(statsText)+2 {
			drawCenteredText(screen, width/2, 4, tcell.StyleDefault, statsText)
		}
	}
	
	// Show help if we have room
	if height > 6 && width > 15 {
		helpText := "ESC to quit"
		drawCenteredText(screen, width/2, 6, tcell.StyleDefault, helpText)
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