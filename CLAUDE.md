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

## Commit Standards
- **Conventional Commits**: Use structured prefixes (`feat:`, `fix:`, `docs:`, `chore:`)
- **Atomic Changes**: Each commit should represent exactly one logical change
- **Meaningful Messages**: Clearly communicate intent and purpose

## Code Principles
- **Structure**: Clean, modular functions with single responsibility
- **Imports**: Standard library first, then external (alphabetically)
- **Naming**: `CamelCase` for exports, `camelCase` for internals
- **Errors**: Explicit handling with contextual messages
- **Types**: Custom types (like `TestState`) to organize related data
- **UI**: Consistent style chains (background → foreground → attributes)
- **Security**: Validate all inputs, use secure defaults
- **Performance**: Monitor critical metrics, establish benchmarks

## Architecture & Design
- **Modularity**: Embrace loose coupling for maintainability
- **Separation**: Keep business logic distinct from infrastructure
- **Resilience**: Design for graceful degradation and recovery
- **Documentation**: Document the "why" behind design decisions

## Testing Standards
- **Coverage**: Maintain high test coverage (unit, integration, end-to-end)
- **TDD**: Adopt Test-Driven Development where feasible
- **Deterministic**: Ensure tests are repeatable and reliable
- **Automation**: Integrate testing into development workflow

## Logging & Observability
- **Structured Logging**: Consistent format across environments
- **Detailed Logs**: Generate comprehensive logs during development
- **Metrics**: Establish meaningful metrics aligned with user experience

## Visual Language
- **Palette**: Fuchsia (titles), Green (correct), Red (errors), Aqua (pending)
- **Layouts**: Box-drawing characters for frames, centered content
- **Content**: EVA/NERV-inspired terminology in user prompts

## File Organization
- `main.go`: Core application logic and UI rendering
- `test-wrap.go`: Text wrapping utilities and tests
- `tests/*.txt`: Text files for typing challenges
- `go.mod/go.sum`: Dependency management

## Continuous Improvement
- **Retrospectives**: Regularly review and identify improvement opportunities
- **Technical Debt**: Actively manage and reduce technical debt
- **Iterative Delivery**: Focus on incremental, deployable changes

## Current TODOs
- Implement word-based line breaking
- Add test history tracking and analytics
- Show quote sources on results screen