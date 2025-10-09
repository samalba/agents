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
	"time"
)

type Workspace struct {
	// The workspace's container state.
	// +internal-use-only
	Container *dagger.Container
	// FIXME: make this configurable
	RootPath string
}

func New(ctx context.Context, container *dagger.Container, rootPath string) *Workspace {
	return &Workspace{
		Container: container,
		RootPath:  rootPath,
	}
}

// Add (install) packages using apk
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
func (w *Workspace) WriteFile(ctx context.Context, filename string, content string) (*Workspace, error) {
	ctr, err := w.Container.
		WithExec([]string{"sh", "-c", fmt.Sprintf("mount -t virtiofs mount0 /mnt && cd '%s' && cat > '%s'", w.RootPath, filename)}, dagger.ContainerWithExecOpts{
			Stdin:                    content,
			InsecureRootCapabilities: true,
		}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}
	w.Container = ctr
	return w, nil
}

// Run shell command in the container and return the output (do not use this to write files or add packages)
func (w *Workspace) RunShellCommand(ctx context.Context, cmd string) (string, error) {
	ctr := w.Container.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{"bash", "-c", fmt.Sprintf("mount -t virtiofs mount0 /mnt && cd '%s' && %s", w.RootPath, cmd)},
			dagger.ContainerWithExecOpts{
				InsecureRootCapabilities: true,
			})
	output, err := ctr.Stdout(ctx)
	if err != nil {
		return "", err
	}
	w.Container = ctr
	return output, nil
}
