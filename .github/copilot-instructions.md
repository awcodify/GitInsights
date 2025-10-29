# Copilot Instructions for GitInsights

## Repository Overview

GitInsights is a Go-based command-line tool that analyzes GitHub profiles and generates detailed statistics about repository languages, commit patterns, and productivity metrics. The tool updates a README.md file with visual insights about a user's GitHub activity.

## Technology Stack

- **Language**: Go 1.21
- **Key Dependencies**:
  - `github.com/google/go-github/v38` - GitHub API client
  - `golang.org/x/oauth2` - OAuth2 authentication
- **Testing**: Standard Go testing framework
- **CI/CD**: GitHub Actions

## Project Structure

```
.
├── main.go           # Main application logic
├── main_test.go      # Unit tests with mocked GitHub API
├── go.mod            # Go module dependencies
├── go.sum            # Dependency checksums
├── README.md         # Project documentation with embedded statistics
└── .github/
    └── workflows/    # GitHub Actions workflows
```

## Key Components

### Core Functionality

1. **GitHub API Integration** (`GitHubClient` interface):
   - `GetUser()` - Retrieve user details
   - `ListRepositories()` - Fetch user's repositories
   - `ListLanguages()` - Get language statistics per repository
   - `ListAllCommits()` - Retrieve commit history across repositories

2. **Statistics Generation**:
   - `summarizeGitHubProfile()` - Aggregates language usage across all repositories
   - `calculateMostProductiveDay()` - Analyzes commit patterns by day of week
   - `calculateMostProductiveTime()` - Analyzes commit patterns by hour
   - `generateMarkdown()` - Creates formatted output with progress bars

3. **README Update**:
   - `updateReadme()` - Replaces content between `<!--START_SECTION:GitInsights-->` and `<!--END_SECTION:GitInsights-->` markers

## Development Workflow

### Building

```bash
go build ./...
```

### Running Tests

```bash
go test -v ./...
```

### Linting

```bash
go vet ./...
gofmt -l .
```

### Running the Application

Requires `GITHUB_TOKEN` environment variable:

```bash
export GITHUB_TOKEN=$(gh auth token)
go run main.go
```

## Code Style and Conventions

1. **Interface-Based Design**: Use the `GitHubClient` interface for GitHub operations to enable mocking in tests
2. **Concurrency**: Use goroutines and channels for parallel API calls (see repository and commit processing)
3. **Error Handling**: Always return and check errors; log non-fatal errors with `log.Printf()`
4. **Testing**: Mock the `GitHubClient` interface using `MockGitHubAPI` for unit tests
5. **Formatting**: Code must pass `gofmt` checks (enforced in CI)

## Testing Guidelines

- All tests use mock implementations of `GitHubClient` interface
- Mock data should be realistic and cover edge cases
- Test functions follow the pattern `Test<FunctionName>`
- Use the `MockGitHubAPI` struct defined in `main_test.go`
- Tests should verify both success cases and error handling

## Important Notes

1. **README Markers**: The application expects `<!--START_SECTION:GitInsights-->` and `<!--END_SECTION:GitInsights-->` markers in README.md
2. **Language Aggregation**: Languages below 5% threshold are combined into "Other" category
3. **Concurrency**: Repository and commit processing uses goroutines for performance
4. **Authentication**: Requires valid GitHub personal access token with appropriate scopes

## GitHub Actions

Three workflows are configured:

1. **test.yaml** - Runs on push/PR: tests, linting, formatting checks, build
2. **build_and_release.yaml** - Creates releases with compiled binaries
3. **update_readme.yaml** - Scheduled/manual execution to update README with latest statistics

## Making Changes

When modifying code:

1. Maintain the interface-based architecture for testability
2. Add corresponding unit tests with mocked dependencies
3. Ensure concurrent operations use proper synchronization (WaitGroups, channels)
4. Run full test suite before committing: `go test -v ./...`
5. Verify formatting: `gofmt -l .` (should return no files)
6. Run linter: `go vet ./...`
7. Test builds successfully: `go build ./...`

## Common Tasks

### Adding New Statistics

1. Create a new function following the pattern of `calculateMostProductiveDay()`
2. Add a method to `GitHubClient` interface if new API calls are needed
3. Implement the method in both `GitHubAPI` and `MockGitHubAPI`
4. Update `generateMarkdown()` to include the new statistic
5. Add comprehensive unit tests

### Modifying Output Format

1. Update `generateMarkdown()` function
2. Maintain the `<!--START_SECTION:GitInsights-->` and `<!--END_SECTION:GitInsights-->` markers
3. Test with `updateReadme()` function to ensure proper section replacement

### Adding API Endpoints

1. Add method to `GitHubClient` interface
2. Implement in `GitHubAPI` struct using `g.Client` (go-github client)
3. Add mock implementation in `MockGitHubAPI` for testing
4. Write unit tests using the mock
