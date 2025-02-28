# Dockerfile Optimizer

A Dagger module that helps optimize your Dockerfiles using AI assistance. This tool analyzes your Dockerfile and suggests improvements for better efficiency, security, and best practices. Once the analysis is complete, it automatically creates a pull request with the suggested optimizations.

This project serves as an example implementation of a simple AI agent using Dagger, demonstrating how to integrate OpenAI's capabilities with GitHub automation in a containerized environment.

## Prerequisites

Before using this module, make sure you have the following:

1. OpenAI API Token
2. GitHub Token (for repository access)

## Environment Setup

Set up the required environment variables:

```bash
export OPENAI_API_KEY="your-openai-api-key"
export GITHUB_TOKEN="your-github-token"
```

## Usage

To use the Dockerfile optimizer, you can analyze Dockerfiles directly from GitHub repositories:

```bash
dagger shell -c "optimize-dockerfile $GITHUB_TOKEN https://github.com/username/repository"
```

The module will:
1. Clone the specified GitHub repository
2. Locate the Dockerfile
3. Analyze it using AI
4. Apply optimization suggestions
5. Create a new pull request with the improvements
6. Return the URL of the created pull request

## Example

```bash
# Example: Optimizing a Dockerfile from a GitHub repository
dagger shell -c "optimize-dockerfile $GITHUB_TOKEN https://github.com/samalba/demo-app"
```

The tool will analyze the Dockerfile and create a pull request with improvements for:
- Security improvements
- Size reduction opportunities
- Best practices recommendations
- Performance optimizations

After execution, the module will output the URL of the newly created pull request, where you can review all suggested changes.
