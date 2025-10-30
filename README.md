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

# 📊 Git Insights Dashboard

[![Profile Stats](https://img.shields.io/badge/Profile-Statistics-blueviolet?style=for-the-badge&logo=github)](https://github.com)

</div>

---

<div align="center">

## 🎯 Quick Stats

</div>

<table align="center">
<tr>
<td align="center" width="200">
<img src="https://img.icons8.com/fluency/96/000000/user.png" width="48"/>
<br><strong>Account Age</strong>
<br><code>9 years 6 months</code>
</td>
<td align="center" width="200">
<img src="https://img.icons8.com/fluency/96/000000/fire-element.png" width="48"/>
<br><strong>Current Streak</strong>
<br><code>2 days</code>
</td>
<td align="center" width="200">
<img src="https://img.icons8.com/fluency/96/000000/trophy.png" width="48"/>
<br><strong>Longest Streak</strong>
<br><code>5 days</code>
</td>
</tr>
</table>

<div align="center">

## ⚡ Productivity Insights

</div>

<table align="center">
<tr>
<td align="center">
<img src="https://img.icons8.com/fluency/96/000000/calendar.png" width="40"/>
<br><strong>Most Productive Day</strong>
<br>🔥 <code>Tuesday</code>
</td>
<td align="center">
<img src="https://img.icons8.com/fluency/96/000000/clock.png" width="40"/>
<br><strong>Peak Hours</strong>
<br>⏰ <code>08:00 - 09:00</code>
</td>
</tr>
</table>

<div align="center">

## 📈 Weekly Activity

</div>

```text
🌙 Monday    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░  116 commits
🔥 Tuesday   ██████████████████████████████  201 commits
💎 Wednesday ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░  142 commits
⚡ Thursday  ██████████████████████████████  201 commits
🎉 Friday    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░  116 commits
🌟 Saturday  ▒▒▒▒▒▒▒▒▒▒▒░░░░░░░░░░░░░░░░░░░   80 commits
☀️ Sunday    ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   45 commits
```

<div align="center">

## 💻 Language Distribution

</div>

<div align="center">

![JavaScript](https://img.shields.io/badge/JavaScript-87.6%25-blue?style=flat-square&logo=javascript&logoColor=white) 
![TypeScript](https://img.shields.io/badge/TypeScript-5.1%25-blue?style=flat-square&logo=typescript&logoColor=white) 
![SCSS](https://img.shields.io/badge/SCSS-2.0%25-blue?style=flat-square&logo=sass&logoColor=white)

![Go](https://img.shields.io/badge/Go-1.7%25-blue?style=flat-square&logo=go&logoColor=white) 
![HTML](https://img.shields.io/badge/HTML-1.2%25-blue?style=flat-square&logo=html5&logoColor=white) 

</div>

<details>
<summary><b>📊 Detailed Breakdown</b></summary>

```text
🟨 JavaScript ███████████████████████████████████░░░░░  87.64%
🔷 TypeScript ▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   5.10%
🎨 SCSS       ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   1.96%
🔵 Go         ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   1.65%
🌐 HTML       ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   1.24%
💎 Ruby       ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.82%
💧 Elixir     ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.68%
🎨 CSS        ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.37%
🐍 Python     ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.35%
🟢 Vim Script ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.10%
💻 Other      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.07%
```

</details>

---

<div align="center">

<sub>📅 Last updated: Thursday, October 30, 2025 at 9:09 AM</sub>

<sub>⚡ Generated with [GitInsights](https://github.com/awcodify/GitInsights)</sub>

</div>

<!--END_SECTION:GitInsights-->