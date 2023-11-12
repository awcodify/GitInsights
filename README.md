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
    - name: Checkout code
      uses: actions/checkout@v2
      
    - name: Download and run GitInsight
      run: |
        wget https://github.com/awcodify/GitInsights/releases/download/v0.1.0/GitInsights -O GitInsights
        chmod +x GitInsights
        ./GitInsights

        # Commit and push changes
        git config --local user.email "awcodify@gmail.com"
        git config --local user.name "awcodify"
        git add .
        git commit -m "Update README.md from GitInsights"
        git push
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
* Most productive day and time
* (Under development, and need your contribution!)

## Sample Result

<!--START_SECTION:GitInsights-->
### Git Insight

Language Statistics:
```
Ruby       [█████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░] 32.97%
Go         [█████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 24.38%
HTML       [█████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 14.20%
CSS        [█████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 12.65%
JavaScript [████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 11.71%
Other      [█░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  4.09%
```

📅 Most Productive Day: Monday

⌚️ Most Productive Hour: 14:00 - 15:00
<!--END_SECTION:GitInsights-->