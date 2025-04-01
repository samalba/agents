# Robot Developer 🤖

## Overview

Robot Developer demonstrates how to build and manage robotic systems using LLM-powered agents and the Dagger platform. This demo showcases an intelligent agent that can assist in developing, debugging, and optimizing robot code—streamlining the development process for robotic applications.

Built with [Dagger](https://dagger.io), the open platform for agentic software.

## Highlights

- **LLM-Driven Development**: The agent leverages OpenAI to understand, generate, and optimize robot code, using structured tools defined in a Dagger workspace.

- **Composable Modules**: Built as a reusable Dagger module, the agent provides capabilities for robot code development, testing, and deployment.

- **Intelligent Assistance**: Offers advanced support for robotics-specific tasks, including motion planning, sensor integration, and control system development.

## How to use it?

Prerequisites:

1. OpenAI API Token
2. GitHub Token (for repository access)

*Note: the Github token can be generated from https://github.com/settings/personal-access-tokens*

Start a dev Dagger Engine with LLM support using: https://docs.dagger.io/ai-agents#initial-setup

Run the robot developer from the Dagger Shell:

```shell
export OPENAI_API_KEY="your-openai-api-key"
export GITHUB_TOKEN="your-github-token"
dagger shell
robot-developer ⋈ ask "Write a curl-like program in your compiled language of choice." | terminal
```

The module will:
1. Provide an interactive development environment for robot code
2. Assist with code generation and optimization
3. Help debug and test robotic systems
4. Support integration with common robotics frameworks
