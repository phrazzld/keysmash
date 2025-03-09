# KEYSMASH: Developer's Reference

## Essential Commands
```bash
go build -o keysmash main.go     # Build executable
./keysmash                       # Run application
go mod tidy                      # Manage dependencies
go fmt ./...                     # Format code
golangci-lint run                # Lint codebase
go run main.go                   # Run without building
go run test-wrap.go              # Run text wrapping tests
```

## Code Principles
- **Structure**: Clean, modular functions with single responsibility
- **Imports**: Standard library first, then external (alphabetically)
- **Naming**: `CamelCase` for exports, `camelCase` for internals
- **Errors**: Explicit handling with contextual messages
- **Types**: Custom types (like `TestState`) to organize related data
- **UI**: Consistent style chains (background → foreground → attributes)

## Visual Language
- **Palette**: Fuchsia (titles), Green (correct), Red (errors), Aqua (pending)
- **Layouts**: Box-drawing characters for frames, centered content
- **Content**: EVA/NERV-inspired terminology in user prompts

## File Organization
- `main.go`: Core application logic and UI rendering
- `test-wrap.go`: Text wrapping utilities and tests
- `tests/*.txt`: Text files for typing challenges
- `go.mod/go.sum`: Dependency management

## Current TODOs
- Implement word-based line breaking
- Add test history tracking and analytics
- Show quote sources on results screen