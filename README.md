# KEYSMASH

A minimalist terminal-based typing test application.

## Features

- **Simple Terminal UI**: Clean, distraction-free interface
- **Dynamic Test Selection**: Random quotes from various works of literature and pop culture
- **Real-time Feedback**: Immediate typing feedback
- **Performance Metrics**: WPM calculation and accuracy tracking
- **Progress Visualization**: Live progress bar and completion percentage
- **Cross-platform**: Works on macOS, Linux, and Windows terminals

## Quick Start

```bash
# Clone and build
git clone https://github.com/phrazzld/keysmash.git
cd keysmash
go mod tidy
go build -o keysmash main.go

# Run
./keysmash
```

## Adding Custom Tests

Place plain text files in the `tests/` directory:

```
tests/
├── your-quote.txt
├── coding-snippet.txt
└── practice-text.txt
```

## Usage

The interface is straightforward:
- Type the displayed text exactly as shown
- Watch your progress with real-time WPM and accuracy stats
- View your performance metrics upon completion
- Commands: `R`: Retry test | `N`: New test | `Q`: Quit

## About

KEYSMASH was built with Go using [tcell](https://github.com/gdamore/tcell) for terminal rendering. The application focuses on providing a clean, distraction-free typing experience that helps users practice and improve their typing speed and accuracy.

The minimalist design prioritizes functionality while maintaining readability and ease of use, making it suitable for regular typing practice sessions.

## License

MIT
