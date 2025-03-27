# Dockerfile Optimizer 🤖

## Overview

Dockerfile Optimizer demonstrates how to automatically improve Dockerfiles using LLM-powered agents and the Dagger platform. This demo shows an end-to-end AI workflow: analyzing a Dockerfile, suggesting improvements, and opening a pull request—all without human intervention.

Built with [Dagger](https://dagger.io), the open platform for agentic software.

## Demo

[![Watch the demo](https://img.youtube.com/vi/WN9IBSD55Kk/hqdefault.jpg)](https://youtu.be/WN9IBSD55Kk)

What this shows:
An AI agent uses the Dagger API to find and optimize a Dockerfile in a GitHub repository. Once optimized, it opens a pull request. A second Dagger-powered agent picks up the PR and reviews it—showing how multiple agents can collaborate in an automated DevOps pipeline.

## Highlights

LLM-Driven Automation: The agent uses OpenAI to interpret and optimize Dockerfiles, relying on prompts and structured tools defined in a Dagger workspace.

Composable Modules: Built as a reusable Dagger module, the agent leverages Dagger’s container API for tasks like cloning repos, reading files, and writing PRs.

Agent Collaboration: Demonstrates a multi-agent workflow, where one agent submits code changes and another reviews them - enabling scalable, intelligent automation in software delivery.

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
