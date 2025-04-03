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
func (m *RobotDeveloper) Ask(ctx context.Context, assignment string) *dagger.Container {
	workspace := dag.Workspace()
	env := dag.Env().
		WithWorkspaceInput("before", workspace, "These are the tools to complete the assignment").
		WithStringInput("assignment", assignment, "This describes the assignment to complete").
		WithWorkspaceOutput("after", "Final state of the workspace after completing the assignment")

	llm := dag.LLM().
		WithEnv(env).
		WithPrompt(`You are a Robot Developer with deep knowledge of programming and software development best practices.
Use the tools in the workspace to complete the assignment. You can install packages, run shell commands and write files.
Do not rely on system libraries dependencies and instead implement the dependencies yourself whenever it's possible.
Do not stop until the code builds with no errors.`).
		Env().Output("after").AsWorkspace()

	return llm.Container()
}
