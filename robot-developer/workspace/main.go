// A generated module for Workspace functions
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
	"dagger/workspace/internal/dagger"
	"fmt"
)

type Workspace struct {
	// The workspace's container state.
	// +internal-use-only
	Container *dagger.Container
}

func New() *Workspace {
	return &Workspace{
		Container: dag.Container().
			From("cgr.dev/chainguard/wolfi-base").
			WithExec([]string{"apk", "add", "bash"}).
			WithWorkdir("/src"),
	}
}

// Add packages to the container, using apk
func (w *Workspace) AddPackages(ctx context.Context, pkgs string) (*Workspace, error) {
	ctr := w.Container.WithExec([]string{"sh", "-c", fmt.Sprintf("apk add %s", pkgs)})
	// Check the packages were added before updating the workspace
	_, err := ctr.Sync(ctx)
	if err != nil {
		return nil, err
	}
	w.Container = ctr
	return w, nil
}

// Search package by name
func (w *Workspace) SearchPackage(ctx context.Context, name string) (string, error) {
	w.Container = w.Container.WithExec([]string{"apk", "search", name})
	return w.Container.Stdout(ctx)
}

// Write a file to the container, takes the filename and content
func (w *Workspace) WriteFile(ctx context.Context, filename string, content string) *Workspace {
	w.Container = w.Container.WithNewFile(filename, content)
	return w
}

// Run shell command in the container and return the output
func (w *Workspace) RunShellCommand(ctx context.Context, cmd string) (*Workspace, error) {
	w.Container = w.Container.WithExec([]string{"bash", "-c", cmd})
	// Check the command was run before updating the workspace
	_, err := w.Container.Sync(ctx)
	if err != nil {
		return nil, err
	}
	return w, nil
}
