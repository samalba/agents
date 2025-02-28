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
	"fmt"
	"path/filepath"
	"strings"

	"dagger/workspace/internal/dagger"

	"github.com/google/uuid"
)

type Workspace struct {
	// The workspace's container state
	// +internal-use-only
	Container *dagger.Container
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
		Container: dag.Container().
			From("cgr.dev/chainguard/wolfi-base").
			WithMountedDirectory("/src", repo).
			WithWorkdir("/src"),
		RepoURL:     repoURL,
		GitHubToken: githubToken,
	}
}

// func (w *Workspace) Ctr(ctx context.Context) *dagger.Container {
// 	return w.Container
// }

// Ensure the path is relative to the git clone directory
func translatePath(path string) string {
	if strings.HasPrefix(path, "/src") {
		return path
	}
	return filepath.Join("/src", path)
}

// Read a file at the given path
func (w *Workspace) Read(ctx context.Context, path string) (string, error) {
	path = translatePath(path)
	return w.Container.File(path).Contents(ctx)
}

// Write a file at the given path with the given content
func (w Workspace) Write(path, content string) Workspace {
	path = translatePath(path)
	w.Container = w.Container.WithNewFile(path, content)
	return w
}

// Build the container from the Dockerfile at the given path
func (w *Workspace) Build(ctx context.Context, path string) error {
	// Split directory and filename from path
	dirname, filename := filepath.Split(path)
	path = translatePath(path)
	_, err := w.Container.Build(w.Container.Directory(dirname), dagger.ContainerBuildOpts{Dockerfile: filename}).Sync(ctx)
	return err
}

// Find files that match the given glob pattern
func (w *Workspace) Find(ctx context.Context,
	pattern string,
) ([]string, error) {
	return w.Container.Directory("/src").Glob(ctx, pattern)
}

// Create a new PullRequest with the changes in the workspace, the given title and body, returns the PR URL
func (w *Workspace) CreatePR(ctx context.Context, title, body string) (string, error) {
	// generate a random branch name
	branchName := "dockerfile-improvements-" + uuid.New().String()[:8]
	// The changeset needs to contain only the Dockerfile otherwise the diff will fail (FIXME?)
	changeset := dag.Directory().WithFile("Dockerfile", w.Container.File("Dockerfile"))
	// Create a new feature branch
	featureBranch := dag.FeatureBranch(w.GitHubToken, w.RepoURL, branchName).
		WithChanges(changeset)

	// Make sure changes have been made to the workspace
	diff, err := featureBranch.Diff(ctx)
	if err != nil {
		return "", err
	}

	if diff == "" {
		return "", fmt.Errorf("got empty diff on feature branch (llm did not make any changes)")
	}

	return featureBranch.PullRequest(ctx, title, body)
}
