# Pull Request Reviewer 🤖

## What is this?

A Dagger module for reviewing Pull Requests and suggesting improvements.

## How to use it?

Start a dev Dagger Engine with LLM support using: https://docs.dagger.io/ai-agents#initial-setup

Trigger a review from the Dagger Shell:

```shell
export GITHUB_TOKEN="your-github-token"
dagger shell -c "review-pr GITHUB_TOKEN <PR_NUMBER> <REPO_URL>"
# Example:
dagger shell -c "review-pr GITHUB_TOKEN 19 https://github.com/samalba/demo-app"
```

*Note: the PR number can be a full PR URL or just the number.*

## Integration with Github action

Ideally, we want to the PR reviewer to run anytime there is a new PR that is ready for review.

Create a github workflow in the corresponding repository at `.github/workflows/pr-reviewer.yaml` with the following content:

```yaml
name: PR Review

on:
  pull_request:
    types: [opened, ready_for_review]

jobs:
  dagger:
    name: Run Dagger Pipeline
    runs-on: ubuntu-latest
    # Skip running on draft PRs
    if: github.event.pull_request.draft == false
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Call Dagger Function to review PR
        uses: dagger/dagger-for-github@8.0.0
        with:
          version: "latest"
          verb: call
          module: github.com/samalba/agents/pr-reviewer
          args: --allow-llm=all review-pr --github-token=env://GH_TOKEN --query=${{ github.event.pull_request.number }} --repo-url=${{ github.event.repository.html_url }}
          cloud-token: ${{ secrets.DAGGER_CLOUD_TOKEN }} # optional
        env:
          GH_TOKEN: ${{ secrets.GH_TOKEN }} # this is not the default provided token, need access to add comments
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
```

*Note: `GH_TOKEN` is not the default provided token, the token need access to add comments to a PR.*
