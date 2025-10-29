# GitInsight - Clean Architecture

A GitHub profile statistics generator that follows Clean Architecture principles.

## Architecture Overview

The project is structured into distinct layers, each with a specific responsibility:

```
GitInsight/
├── domain/              # Core business entities and interfaces
│   ├── entities.go      # Business entities (ProfileStats, LanguageStats, Commit)
│   └── repository.go    # Repository interfaces (contracts)
├── usecase/             # Business logic orchestration
│   ├── profile_stats.go # Profile statistics use case
│   └── profile_stats_test.go
├── infrastructure/      # External dependencies & implementations
│   ├── github_client.go # GitHub API implementation
│   └── file_manager.go  # File operations implementation
├── presentation/        # Output formatting
│   └── markdown_generator.go
└── main.go             # Application entry point & dependency wiring
```

## Layer Responsibilities

### 1. Domain Layer (`domain/`)
- **Purpose**: Core business entities and repository interfaces
- **Dependencies**: None (completely independent)
- **Contents**:
  - `entities.go`: Business entities like `ProfileStats`, `LanguageStats`, `Commit`
  - `repository.go`: Repository interfaces that define contracts for data access

### 2. Use Case Layer (`usecase/`)
- **Purpose**: Business logic orchestration
- **Dependencies**: Only depends on `domain` layer
- **Contents**:
  - `profile_stats.go`: Orchestrates fetching data, calculating statistics, and preparing results
  - Contains pure business logic: sorting, filtering, aggregating data
  - Independent of external frameworks and UI

### 3. Infrastructure Layer (`infrastructure/`)
- **Purpose**: Implementations of external dependencies
- **Dependencies**: `domain` layer interfaces, external packages (go-github, oauth2)
- **Contents**:
  - `github_client.go`: Implements GitHub API calls
  - `file_manager.go`: Handles file I/O operations
  - These are concrete implementations of domain repository interfaces

### 4. Presentation Layer (`presentation/`)
- **Purpose**: Output formatting and display logic
- **Dependencies**: `domain` layer entities
- **Contents**:
  - `markdown_generator.go`: Converts domain entities to markdown format
  - Responsible for visual representation only

### 5. Main (`main.go`)
- **Purpose**: Application entry point and dependency injection
- **Dependencies**: All layers
- **Responsibilities**:
  - Creates concrete implementations
  - Wires dependencies together
  - Orchestrates the application flow

## Benefits of This Architecture

1. **Separation of Concerns**: Each layer has a single, well-defined responsibility
2. **Testability**: Easy to mock dependencies and test each layer independently
3. **Maintainability**: Changes in one layer don't cascade to others
4. **Flexibility**: Easy to swap implementations (e.g., switch from GitHub to GitLab)
5. **Clean Dependencies**: Dependencies flow inward (infrastructure → usecase → domain)

## Dependency Flow

```
main.go
  ↓
infrastructure (implements) → domain (interfaces)
  ↓                              ↑
usecase (uses) ←─────────────────┘
  ↓
presentation (formats)
```

## Testing Strategy

- **Domain**: Test pure business entities (minimal testing needed)
- **Use Case**: Test with mocked repositories (unit tests)
- **Infrastructure**: Test with real APIs (integration tests) or mocks
- **Presentation**: Test output formatting with sample data

## Running the Application

```bash
# Set your GitHub token
export GITHUB_TOKEN="your_token_here"

# Run the application
go run main.go
```

## Example: Adding a New Data Source

To add support for GitLab:

1. Create `infrastructure/gitlab_client.go` implementing `domain.GitHubRepository`
2. In `main.go`, swap `infrastructure.NewGitHubClient()` with `infrastructure.NewGitLabClient()`
3. No changes needed in use case or presentation layers!

This demonstrates the power of clean architecture - you can swap entire data sources without touching business logic.
