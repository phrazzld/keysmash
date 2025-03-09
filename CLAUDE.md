# TYPING-EVA: Developer's Reference

## Essential Commands
```bash
go build -o typing-eva main.go    # Build executable
./typing-eva                      # Run application
go mod tidy                       # Manage dependencies
go fmt ./...                      # Format code
golangci-lint run                 # Lint codebase
```

## Code Principles
- **Structure**: Clean, modular functions with single responsibility
- **Imports**: Standard library first, then external (alphabetically)
- **Naming**: `CamelCase` for exports, `camelCase` for internals
- **Errors**: Explicit handling with contextual messages
- **UI**: Consistent style chains (background → foreground → attributes)

## Visual Language
- **Palette**: Fuchsia (titles), Green (correct), Red (errors), Aqua (pending)
- **Layouts**: Box-drawing characters for frames, centered content
- **Content**: EVA-inspired terminology in user prompts

## File Organization
- `main.go`: Core application logic
- `tests/*.txt`: Text files for typing challenges
- `go.mod/go.sum`: Dependency management