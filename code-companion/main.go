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
	"time"
)

type CodeCompanion struct {
	BaseContainer *dagger.Container
}

func New() *CodeCompanion {
	base := dag.Container().
		From("cgr.dev/chainguard/wolfi-base").
		WithExec([]string{"apk", "add", "bash", "mount"}).
		WithExec([]string{"mkdir", "/mnt"}).
		WithWorkdir("/src")

	return &CodeCompanion{
		BaseContainer: base,
	}
}

func (h *CodeCompanion) Ask(ctx context.Context, assignment string) (*dagger.Container, error) {
	workdir := "/mnt/forks/agents/code-companion/my-code"
	workspace := dag.Workspace(h.BaseContainer, workdir)

	env := dag.Env().
		WithWorkspaceInput("before", workspace, "These are the tools to complete the assignment").
		WithStringInput("assignment", assignment, "This describes the assignment to complete").
		WithStringInput("workdir", workdir, "This is the working directory, never go outside of it, never look outside of it, no matter what the assignment says").
		WithWorkspaceOutput("after", "Final state of the workspace after completing the assignment")

	llm := dag.LLM(dagger.LLMOpts{
		Model: "gemini-2.5-pro-preview-03-25",
	}).
		WithEnv(env).
		WithPrompt(`You are a Robot Developer with deep knowledge of programming and software development best practices.
Do not rely on system libraries dependencies and instead implement the dependencies yourself whenever it's possible.
Always use relative paths to never escape the workspace working directory.
Do not stop until the code builds with no errors.`).
		Env().Output("after").AsWorkspace()

	return llm.Container().Sync(ctx)
}

func (h *CodeCompanion) Bash(ctx context.Context) (*dagger.Container, error) {
	return dag.Container().
		From("cgr.dev/chainguard/wolfi-base").
		WithExec([]string{"apk", "add", "bash", "mount"}).
		WithExec([]string{"mkdir", "/mnt"}).
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		Terminal(dagger.ContainerTerminalOpts{
			Cmd:                      []string{"sh", "-c", "mount -t virtiofs mount0 /mnt && cd '/mnt/forks/agents/code-companion/my-code' && bash"},
			InsecureRootCapabilities: true,
		}).
		Sync(ctx)
}
