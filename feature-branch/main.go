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
	"encoding/json"
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
func New(ctx context.Context, githubToken *dagger.Secret, repoURL string, branchName string) *FeatureBranch {
	repoURL = strings.TrimSuffix(repoURL, ".git")

	return &FeatureBranch{
		Ctr: dag.Container().
			From("cgr.dev/chainguard/wolfi-base:latest").
			WithExec([]string{"apk", "add", "git", "gh", "rsync"}).
			WithSecretVariable("GITHUB_TOKEN", githubToken).
			WithExec([]string{"git", "config", "--global", "user.email", "sam+module-feature-branch@dagger.io"}).
			WithExec([]string{"git", "config", "--global", "user.name", "Dagger Agent"}).
			WithExec([]string{"gh", "auth", "setup-git"}).
			WithExec([]string{"gh", "repo", "clone", repoURL, "/src"}).
			WithWorkdir("/src"),
		BranchName: branchName,
	}
}

// Set the branch name to a new unique name
func (m *FeatureBranch) WithNewUniqueBranchName() *FeatureBranch {
	m.BranchName = m.BranchName + "-" + uuid.New().String()[:8]
	return m
}

// Set changeset of the feature branch
func (m *FeatureBranch) WithChanges(changes *dagger.Directory) *FeatureBranch {
	m.Changes = changes
	return m
}

func applyChanges(baseImage *dagger.Container, changes *dagger.Directory) *dagger.Container {
	return baseImage.
		WithMountedDirectory("/changes", changes).
		WithExec([]string{"rsync", "-a", "/changes/", "/src"})
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

	return applyChanges(m.Ctr, m.Changes).
		WithExec(diffArgs).
		Stdout(ctx)
}

// Commit the changes
func (m *FeatureBranch) Commit(ctx context.Context, message string) (*FeatureBranch, error) {
	if m.Changes == nil {
		return nil, errors.New("no changes to commit")
	}

	m.Ctr = applyChanges(m.Ctr, m.Changes).
		WithExec([]string{"git", "checkout", "-b", m.BranchName}).
		WithExec([]string{"git", "add", "."}).
		WithExec([]string{"git", "commit", "-m", message})

	_, err := m.Ctr.Sync(ctx)
	return m, err
}

// Push the changes to the remote branch
func (m *FeatureBranch) Push(ctx context.Context) (*FeatureBranch, error) {
	_, err := m.Ctr.WithExec([]string{"git", "push", "origin", m.BranchName}).Sync(ctx)
	return m, err
}

// Opens a Pull Request on GitHub
func (m *FeatureBranch) CreatePullRequest(ctx context.Context, title string, body string, draft bool) (string, error) {
	prArgs := []string{"gh", "pr", "create", "--head", m.BranchName}
	if title == "" || body == "" {
		prArgs = append(prArgs, "--fill")
	} else {
		prArgs = append(prArgs, "--title", title, "--body", body)
	}
	if draft {
		prArgs = append(prArgs, "--draft")
	}

	out, err := m.Ctr.WithExec(prArgs).Stdout(ctx)
	if err != nil {
		return "", err
	}

	// Grab the last line of the output, which is the PR URL
	lines := strings.Split(strings.TrimSpace(out), "\n")
	prURL := strings.TrimSpace(lines[len(lines)-1])
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

	env := dag.Env().
		WithStringInput("diff", diff, "This the diff of the PullRequest to review").
		WithStringInput("additionalContext", additionalContext, "Optional additional context for the PR review (can be empty)")

	// Augment the PR with the LLM
	llm := dag.LLM().
		WithEnv(env).
		WithPrompt(`Generate a detailed description of the changes in the PR.
Include the following information:
- The changes made to the code
- The rationale for the changes
- Any potential risks or considerations
- Any other relevant details

PR Diff:
$diff

Additional Context:
$additionalContext
`)
	generatedDescription, err := llm.LastReply(ctx)
	if err != nil {
		return "", err
	}
	// Update the PR with the LLM's description
	_, err = m.Ctr.
		WithMountedDirectory("/input", dag.Directory().WithNewFile("body-file.txt", generatedDescription)).
		WithExec([]string{"gh", "pr", "edit", "--body-file", "/input/body-file.txt"}).
		WithExec([]string{"gh", "pr", "ready"}).
		Sync(ctx)

	return prURL, err
}

// Checkout a Pull Request code
// Query is any argument supported by the gh cli (gh pr checkout [<number> | <url> | <branch>])
func (m *FeatureBranch) CheckoutPullRequest(ctx context.Context, query string) (*FeatureBranch, error) {
	m.Ctr = m.Ctr.
		WithExec([]string{"gh", "pr", "checkout", query})

	return m, nil
}

// Get the body of a Pull Request, returns title, body in a slice
func (m *FeatureBranch) GetPullRequestBodyTitle(ctx context.Context) ([]string, error) {
	data, err := m.Ctr.
		WithExec([]string{"gh", "pr", "view", "--json", "body", "--json", "title"}).
		Stdout(ctx)

	if err != nil {
		return nil, err
	}

	var body struct {
		Body  string `json:"body"`
		Title string `json:"title"`
	}

	if err := json.Unmarshal([]byte(data), &body); err != nil {
		return nil, err
	}

	return []string{body.Title, body.Body}, nil
}

// Get the diff of a Pull Request
func (m *FeatureBranch) GetPullRequestDiff(ctx context.Context, query string) (string, error) {
	diff, err := m.Ctr.WithExec([]string{"gh", "pr", "diff"}).Stdout(ctx)
	if err != nil {
		return "", err
	}

	return diff, nil
}

// Add a comment on the PullRequest, set editLast to true to edit the last comment
func (m *FeatureBranch) AddPullRequestComment(ctx context.Context, comment string, editLast bool) (string, error) {
	ghArgs := []string{"gh", "pr", "comment", "--body", comment}
	if editLast {
		ghArgs = append(ghArgs, "--edit-last")
	}
	out, err := m.Ctr.WithExec(ghArgs).Stdout(ctx)
	return out, err
}
