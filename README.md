# keysmash

A sleek CLI typing test.

## Features

- **Immersive Terminal Experience**: Terminal-based UI with Eva-inspired style
- **Dynamic Test Selection**: Random quotes from local text files
- **Real-time Feedback**: Color-coded input (green for correct, red for errors)
- **Performance Metrics**: WPM calculation and error tracking
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

Once running, the interface is intuitive:
- Type the displayed text exactly as shown
- Watch real-time feedback as you type
- View your WPM and error count upon completion
- `R`: Retry test | `N`: New test | `Q`: Quit

## About

Built with Go using [tcell](https://github.com/gdamore/tcell) for terminal rendering.

## License

MIT
