// A generated module for DockerfileOptimizer functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"strings"

	"dagger/dockerfile-optimizer/internal/dagger"
)

type DockerfileOptimizer struct{}

// Optimize a Dockerfile
func (m *DockerfileOptimizer) OptimizeDockerfile(ctx context.Context, githubToken *dagger.Secret, repoURL string) (string, error) {
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// Create a new workspace, using third-party module
	ws := dag.Workspace(githubToken, repoURL)
	// Run the agent loop in the workspace
	after := dag.Llm().
		WithWorkspace(ws).
		WithPrompt(`
You are a Platform Engineer with deep knowledge of Dockerfiles. You have access to a workspace.
Use the read, write, build, find, create-pr tools to complete the following assignment.

1. Look for a Dockerfile in the workspace (use the find tool with the "*Dockerfile*" pattern).
2. Read the Dockerfile and optimize it for reducing its size, number of layers,
and build time. And if possible, increasing the security level of the image by implementing best practices.
3. Write the optimized Dockerfile to the workspace, in the same directory, replacing the original one.
4. Build the container from the optimized Dockerfile to ensure it builds.
5. Create a Pull Request with the changes in the workspace, with a useful title, the body should include the whole explanation of the changes made to the Dockerfile.
6. Return the Pull Request URL as the last message.

If the container build on step 4, read the build error and try to fix it in the Dockerfile until the container builds.
`)
	// Return the last message which containers the Pull Request URL
	return after.LastReply(ctx)
}
