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
	"path/filepath"
	"strings"

	"dagger/workspace/internal/dagger"
)

type Workspace struct {
	// The workspace's directory state
	// +internal-use-only
	Workdir *dagger.Directory
	// Repository URL
	// +internal-use-only
	RepoURL string
	// GitHub token
	// +internal-use-only
	GitHubToken *dagger.Secret
}

func New(githubToken *dagger.Secret, repoURL string) Workspace {
	if !strings.HasSuffix(repoURL, ".git") {
		repoURL = repoURL + ".git"
	}

	// Clone the repository
	repo := dag.Git(repoURL).Head().Tree()

	return Workspace{
		// Build a base container optimized for Go development
		Workdir:     repo,
		RepoURL:     repoURL,
		GitHubToken: githubToken,
	}
}

// Read a file at the given path
func (w *Workspace) Read(ctx context.Context, path string) (string, error) {
	return w.Workdir.File(path).Contents(ctx)
}

// Write a file at the given path with the given content
func (w Workspace) Write(path, content string) Workspace {
	w.Workdir = w.Workdir.WithNewFile(path, content)
	return w
}

// Build the container from the Dockerfile at the given path
func (w *Workspace) Build(ctx context.Context, path string) error {
	// Split directory and filename from path
	dirname, filename := filepath.Split(path)
	_, err := dag.Container().
		Build(w.Workdir.Directory(dirname), dagger.ContainerBuildOpts{Dockerfile: filename}).
		Sync(ctx)
	return err
}
