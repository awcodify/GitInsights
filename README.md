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
      with:
        go-version: 1.17

    - name: Install gh CLI
      run: |
        sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-key C99B11DEB97541F0
        sudo apt-add-repository -y https://cli.github.com/packages
        sudo apt-get update
        sudo apt-get install gh

    - name: Set up Git Insight
      run: |
        git clone https://github.com/awcodify/gitinsights.git $HOME/gitinsights
        go run $HOME/gitinsights/main.go
      env:
        GH_TOKEN: ${{ secrets.GH_TOKEN }}
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
```
## Metrics

* Language Usage: Provides a breakdown of the languages used across your repositories.
* (Under development, and need your contribution!)

## Sample Result

<!--START_SECTION:GitInsights-->
### Git Insight

Language Statistics:
Ruby       ████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 32.50%
Go         █████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 23.72%
HTML       █████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 14.01%
CSS        █████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 12.84%
JavaScript ████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 11.91%
Other      ██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  5.03%
<!--END_SECTION:GitInsights-->