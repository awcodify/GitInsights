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

[![Profile Stats](https://img.shields.io/badge/Git-Insights-blueviolet?style=for-the-badge&logo=github)](https://github.com/awcodify/GitInsights)

</div>

---

<div align="center">

## 🎯 Quick Stats

</div>

<table align="center">
<tr>
<td align="center" width="200">
<img src="https://img.icons8.com/fluency/96/000000/resume.png" width="48"/>
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
<br>💚 <code>Thursday</code>
</td>
<td align="center">
<img src="https://img.icons8.com/fluency/96/000000/clock.png" width="40"/>
<br><strong>Peak Hours</strong>
<br>⏰ <code>08:00 - 09:00</code>
</td>
</tr>
</table>

<div align="center">

## 📊 Weekly Progress

</div>

<table align="center">
<tr>
<td align="center" width="200">
<img src="https://img.icons8.com/fluency/96/000000/calendar-7.png" width="48"/>
<br><strong>This Week</strong>
<br><code>29 commits</code>
</td>
<td align="center" width="200">
<img src="https://img.icons8.com/fluency/96/000000/calendar-6.png" width="48"/>
<br><strong>Last Week</strong>
<br><code>6 commits</code>
</td>
<td align="center" width="200">
<img src="https://img.icons8.com/fluency/96/000000/arrow-up.png" width="48"/>
<br><strong>Growth</strong>
<br><code>+23 (+383.3%)</code>
</td>
</tr>
</table>

<div align="center">

## 📈 Weekly Activity

</div>

```text
🌙 Monday     ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░  112 commits
🔥 Tuesday    ███████████████████████████░░░  188 commits
💎 Wednesday  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░  130 commits
💚 Thursday   ██████████████████████████████  205 commits
🎉 Friday     ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░  111 commits
🌟 Saturday   ▒▒▒▒▒▒▒▒▒▒▒░░░░░░░░░░░░░░░░░░░   80 commits
☀️ Sunday     ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   41 commits
```

<div align="center">

## 📅 Overtime Activity (Last 6 Months)

</div>

```text
May 2025   ██████████████████████████████  203 commits
Jun 2025   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    0 commits
Jul 2025   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    0 commits
Aug 2025   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   46 commits
Sep 2025   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    4 commits
Oct 2025   ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   35 commits
```

<div align="center">

## 💻 Language Distribution

</div>

<div align="center">

![TypeScript](https://img.shields.io/badge/TypeScript-34.1%25-blue?style=flat-square&logo=typescript) ![JavaScript](https://img.shields.io/badge/JavaScript-21.0%25-blue?style=flat-square&logo=javascript) ![SCSS](https://img.shields.io/badge/SCSS-13.1%25-blue?style=flat-square&logo=sass) ![Go](https://img.shields.io/badge/Go-12.2%25-blue?style=flat-square&logo=go) ![HTML](https://img.shields.io/badge/HTML-6.9%25-blue?style=flat-square&logo=html5)

</div>

<details>
<summary><b>📊 Detailed Breakdown</b></summary>

```text
🔷 TypeScript ▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░  34.09%
🟨 JavaScript ▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  21.00%
🎨 SCSS       ▒▒▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  13.14%
🔵 Go         ▒▒▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  12.19%
🌐 HTML       ▒▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   6.85%
💧 Elixir     ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   4.57%
💎 Ruby       ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   2.52%
🐍 Python     ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   2.31%
🎨 CSS        ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   2.17%
🟢 Vim Script ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.70%
💻 Other      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   0.44%
```

</details>

---

<div align="center">

<sub>📅 Last updated: Thursday, October 30, 2025 at 10:31 AM</sub>

<sub>⚡ Generated with [GitInsights](https://github.com/awcodify/GitInsights)</sub>

</div>

<!--END_SECTION:GitInsights-->