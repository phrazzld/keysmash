# NERV Typing Simulation

A sleek CLI typing test inspired by Neon Genesis Evangelion's aesthetic—neon colors and retro futurism meet typing practice.

## Features

- **Immersive Terminal Experience**: Terminal-based UI with Eva-inspired style
- **Dynamic Test Selection**: Random quotes from local text files
- **Real-time Feedback**: Color-coded input (green for correct, red for errors)
- **Performance Metrics**: WPM calculation and error tracking
- **Cross-platform**: Works on macOS, Linux, and Windows terminals

## Quick Start

```bash
# Clone and build
git clone https://github.com/yourusername/typing-eva.git
cd typing-eva
go mod tidy
go build -o typing-eva main.go

# Run
./typing-eva
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

Built with Go using [tcell](https://github.com/gdamore/tcell) for terminal rendering. The aesthetic draws inspiration from NERV terminals in Evangelion, utilizing a neon color palette and minimalist design.

## License

MIT