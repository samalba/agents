// A generated module for RobotDeveloper functions
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
	"dagger/robot-developer/internal/dagger"
)

type RobotDeveloper struct{}

// Returns a container that echoes whatever string argument is provided
func (m *RobotDeveloper) Ask(ctx context.Context, prompt string) *dagger.Container {
	return dag.LLM().
		WithWorkspace(dag.Workspace()).
		WithPromptVar("prompt", prompt).
		WithPrompt(`You are a Robot Developer with deep knowledge of programming and software development best practices.
You have acess to workspace which you can use to install system packages, write files, run commands, etc.
The assignment must result in a final workspace that contains the result of your work.

Assignment: $prompt
`).
		Workspace().
		Container()
}
