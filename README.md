# GitInsights

Git Insights is a tool that provides a summary of your GitHub profile, including language usage in repositories.


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

### Command-Line Options

By default, GitInsights **excludes forked repositories** from analysis. To include forks:

```bash
./GitInsights --include-forks
```

Or with `go run`:
```bash
go run main.go --include-forks
```

You can limit the number of languages displayed in the language statistics (default is 10):

```bash
./GitInsights --max-visible-language 5
```

Or with `go run`:
```bash
go run main.go --max-visible-language 5
```

Combine multiple options:
```bash
./GitInsights --include-forks --max-visible-language 15
```

### Authentication

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

<div align="center">

# 📊 Git Insights

</div>

<div align="center">

![Account Age](https://img.shields.io/badge/Account_Age-9_years_6_months-blue?style=for-the-badge&logo=github)
![Current Streak](https://img.shields.io/badge/Current_Streak-1_days-orange?style=for-the-badge&logo=fire)
![Longest Streak](https://img.shields.io/badge/Longest_Streak-5_days-red?style=for-the-badge&logo=trophy)

![Most Productive Day](https://img.shields.io/badge/Most_Productive_Day-Tuesday-green?style=for-the-badge&logo=calendar)
![Most Productive Hour](https://img.shields.io/badge/Most_Productive_Hour-08:00_--_09:00-purple?style=for-the-badge&logo=clock)

</div>

## 📊 Weekly Commit Distribution

```text
📅 Monday    [█████████████████░░░░░░░░░░░░░] 116 commits
📅 Tuesday   [██████████████████████████████] 201 commits
📅 Wednesday [█████████████████████░░░░░░░░░] 142 commits
📅 Thursday  [█████████████████████████████░] 200 commits
📅 Friday    [█████████████████░░░░░░░░░░░░░] 116 commits
🎉 Saturday  [███████████░░░░░░░░░░░░░░░░░░░] 80 commits
🎉 Sunday    [██████░░░░░░░░░░░░░░░░░░░░░░░░] 45 commits
```

## 💻 Language Statistics

```text
🥇 JavaScript [██████████████████████████░░░░] 87.64%
🥈 TypeScript [█░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  5.10%
🥉 SCSS       [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  1.96%
   Go         [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  1.65%
   HTML       [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  1.24%
   Ruby       [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  0.82%
   Elixir     [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  0.68%
   CSS        [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  0.37%
   Python     [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  0.35%
   Vim Script [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  0.10%
   Other      [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  0.07%
```

---

<div align="center">

⏰ _Last updated: 2025-10-29 14:43:18_

**✨ Generated with [GitInsights](https://github.com/awcodify/GitInsights) ✨**

</div>

<!--END_SECTION:GitInsights-->