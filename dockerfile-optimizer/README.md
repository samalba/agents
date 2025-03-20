# Dockerfile Optimizer 🤖

## What is this?

A Dagger module for optimizing Dockerfiles using AI assistance. This tool analyzes your Dockerfile and suggests improvements for better efficiency, security, and best practices. Once the analysis is complete, it automatically creates a pull request with the suggested optimizations.

## How to use it?

Prerequisites:

1. OpenAI API Token
2. GitHub Token (for repository access)

*Note: the Github token can be generated from https://github.com/settings/personal-access-tokens*

Start a dev Dagger Engine with LLM support using: https://docs.dagger.io/ai-agents#initial-setup

Run an optimization from the Dagger Shell:

```shell
export OPENAI_API_KEY="your-openai-api-key"
export GITHUB_TOKEN="your-github-token"
dagger shell -c "optimize-dockerfile GITHUB_TOKEN <REPO_URL>"
# Example:
dagger shell -c "optimize-dockerfile GITHUB_TOKEN https://github.com/samalba/demo-app"
```

The module will:
1. Clone the specified GitHub repository
2. Locate the Dockerfile
3. Analyze it using AI
4. Apply optimization suggestions
5. Create a new pull request with the improvements
6. Return the URL of the created pull request

## Bonus point

Also enable [the PR reviewer agent](../pr-reviewer) so the PR can be reviewed by another Agent.
