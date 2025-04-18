// A generated module for HostSync functions
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
	"dagger/host-sync/internal/dagger"
)

type CodeCompanion struct {
	BaseContainer *dagger.Container
}

func New() *CodeCompanion {
	base := dag.Container().
		From("cgr.dev/chainguard/wolfi-base").
		WithExec([]string{"apk", "add", "bash"}).
		WithMountedCache("/src", dag.CacheVolume("MyCode"), dagger.ContainerWithMountedCacheOpts{
			Sharing: dagger.CacheSharingModeShared,
		}).
		WithWorkdir("/src")

	return &CodeCompanion{
		BaseContainer: base,
	}
}

func (h *CodeCompanion) Ask(ctx context.Context, assignment string) (*dagger.Container, error) {
	workspace := dag.Workspace(h.BaseContainer)

	env := dag.Env().
		WithWorkspaceInput("before", workspace, "These are the tools to complete the assignment").
		WithStringInput("assignment", assignment, "This describes the assignment to complete").
		WithWorkspaceOutput("after", "Final state of the workspace after completing the assignment")

	llm := dag.LLM().
		WithEnv(env).
		WithPrompt(`You are a Robot Developer with deep knowledge of programming and software development best practices.
Use the tools in the workspace to complete the assignment. You can install packages, run shell commands and write files.
When you cannot use a specific tool, use the shell to run commands.
Do not rely on system libraries dependencies and instead implement the dependencies yourself whenever it's possible.
Always use relative paths to never escape the workspace working directory.
Do not stop until the code builds with no errors.`).
		Env().Output("after").AsWorkspace()

	return llm.Container().Sync(ctx)
}

// Reset the workspace by removing all files (does not work)
// func (h *CodeCompanion) ResetWorkspace(ctx context.Context) error {
// 	// Empty the shared volume
// 	_, err := h.BaseContainer.
// 		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
// 		WithExec([]string{"rm", "-rf", "*"}).Sync(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
