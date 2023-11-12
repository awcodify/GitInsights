# GitInsights

Git Insights is a tool that provides a summary of your GitHub profile, including language usage in repositories.

## Usage

### GitHub Action

Automatically update your README.md with Git Insights using the following GitHub Action workflow:

```yaml
# .github/workflows/update-readme.yml

name: Update Readme

on:
  schedule:
    - cron: '0 0 * * *' # Run daily at midnight (UTC)
  workflow_dispatch: # Trigger manually if needed

jobs:
  update_readme:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2

    - name: Set up Git Insight
      run: |
        git clone https://github.com/awcodify/gitinsights.git $HOME/gitinsights
        go run $HOME/gitinsights/main.go
      env:
        GITHUB_TOKEN: ${{ secrets.GH_TOKEN }}
```
Ensure you have added the GH_TOKEN secret with the necessary permissions.

## Manual Execution
You can also run Git Insights manually. Clone the repository and execute the following command:

```bash
go run main.go
```
Make sure you already logged in to Github with:
```bash
gh auth login

export GITHUB_TOKEN=$(gh auth token)
```
## Metrics

* Language Usage: Provides a breakdown of the languages used across your repositories.
* (Under development, and need your contribution!)

## Sample Result

<!--START_SECTION:GitInsights-->
### Git Insight

Language Statistics:
```
Ruby       [█████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░] 33.01%
Go         [█████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 24.28%
HTML       [█████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 14.22%
CSS        [█████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 12.66%
JavaScript [████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 11.72%
Other      [█░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  4.10%
```

📅 Most Productive Day: Monday

⌚️ Most Productive Hour: 14:00 - 15:00
<!--END_SECTION:GitInsights-->t