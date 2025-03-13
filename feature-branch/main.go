// A generated module for FeatureBranch functions
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
	"dagger/feature-branch/internal/dagger"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type FeatureBranch struct {
	// +internal-use-only
	Ctr *dagger.Container
	// +internal-use-only
	BranchName string
	// +internal-use-only
	Changes *dagger.Directory
}

// Initialize a new feature branch
func New(ctx context.Context, githubToken *dagger.Secret, repoURL string, branchPrefix string) *FeatureBranch {
	return &FeatureBranch{
		Ctr: dag.Container().
			From("cgr.dev/chainguard/wolfi-base:latest").
			WithExec([]string{"apk", "add", "git", "gh", "rsync"}).
			WithMountedSecret("/secrets/github-token", githubToken).
			WithExec([]string{"gh", "auth", "login", "--with-token", "/secrets/github-token"}).
			WithExec([]string{"gh", "repo", "clone", repoURL, "/src"}).
			WithWorkdir("/src"),
		BranchName: branchPrefix + "-" + uuid.New().String()[:8],
	}
}

// Set changeset of the feature branch
func (m *FeatureBranch) WithChanges(changes *dagger.Directory) *FeatureBranch {
	m.Changes = changes
	return m
}

func applyChanges(ctx context.Context, baseImage *dagger.Container, changes *dagger.Directory) *dagger.Container {
	return baseImage.
		WithMountedDirectory("/changes", changes).
		WithExec([]string{"rsync", "-a", "/changes", "/src"})
}

// Diff the changeset of the feature branch
func (m *FeatureBranch) Diff(ctx context.Context, namesOnly bool) (string, error) {
	if m.Changes == nil {
		return "", errors.New("no changes to diff")
	}

	diffArgs := []string{"git", "diff"}
	if namesOnly {
		diffArgs = append(diffArgs, "--name-only")
	}

	return applyChanges(ctx, m.Ctr, m.Changes).
		WithExec(diffArgs).
		Stdout(ctx)
}

// Commit the changes
func (m *FeatureBranch) Commit(ctx context.Context, message string) error {
	if m.Changes == nil {
		return errors.New("no changes to commit")
	}

	m.Ctr = applyChanges(ctx, m.Ctr, m.Changes).
		WithExec([]string{"git", "checkout", "-b", m.BranchName}).
		WithExec([]string{"git", "add", "."}).
		WithExec([]string{"git", "commit", "-m", message})

	_, err := m.Ctr.Sync(ctx)
	return err
}

// Push the changes to the remote branch
func (m *FeatureBranch) Push(ctx context.Context) error {
	_, err := m.Ctr.WithExec([]string{"git", "push", "origin", m.BranchName}).Sync(ctx)
	return err
}

// Opens a Pull Request on GitHub
func (m *FeatureBranch) CreatePullRequest(ctx context.Context, title string, body string, draft bool) (string, error) {
	prArgs := []string{"gh", "pr", "create"}
	if title == "" || body == "" {
		prArgs = append(prArgs, "--fill")
	} else {
		prArgs = append(prArgs, "--title", title, "--body", body)
	}
	if draft {
		prArgs = append(prArgs, "--draft")
	}

	prURL, err := m.Ctr.WithExec(prArgs).Stdout(ctx)
	if err != nil {
		return "", err
	}

	prURL = strings.TrimSpace(prURL)
	return prURL, nil
}

// Opens a Pull Request on GitHub, with the help of an LLM
func (m *FeatureBranch) CreatePullRequestWithLLM(ctx context.Context, additionalContext string) (string, error) {
	// First create a draft PR with filled title and body
	prURL, err := m.CreatePullRequest(ctx, "", "", true)
	if err != nil {
		return "", err
	}
	// Get the diff of the PR
	diff, err := m.Ctr.WithExec([]string{"gh", "pr", "diff"}).Stdout(ctx)
	if err != nil {
		return "", err
	}
	// Augment the PR with the LLM
	llm := dag.Llm().
		WithPromptVar("diff", diff).
		WithPromptVar("additionalContext", additionalContext).
		WithPrompt(`Generate a detailed description of the changes in the PR.
		Include the following information:
		- The changes made to the code
		- The rationale for the changes
		- Any potential risks or considerations
		- Any other relevant details

		Take into account the following additional context:
		$additionalContext

		And take into account the following diff, this is the output of the git diff command:
		$diff

		Only output the description, nothing else.
		`)
	generatedDescription, err := llm.LastReply(ctx)
	if err != nil {
		return "", err
	}
	// Update the PR with the LLM's description
	_, err = m.Ctr.
		WithMountedDirectory("/input", dag.Directory().WithNewFile("body-file.txt", generatedDescription)).
		WithExec([]string{"gh", "pr", "edit", "--body-file", "/input/body-file.txt"}).
		Sync(ctx)

	return prURL, nil
}
