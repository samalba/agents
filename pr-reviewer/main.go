// A generated module for PrReviewer functions
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
	"dagger/pr-reviewer/internal/dagger"
)

type PrReviewer struct{}

// Review a PR, by default it will review the PR description and the diff.
// Query is any argument supported by the gh cli (gh pr view [<number> | <url> | <branch>]).
// Additional instructions can be provided to the LLM to guide the review.
func (m *PrReviewer) ReviewPr(ctx context.Context, githubToken *dagger.Secret, repoURL string, query, additionalInstructions string) ([]string, error) {
	prCheckout := dag.FeatureBranch(githubToken, repoURL, "unused-branch-name").CheckoutPullRequest(query)

	prBody, err := prCheckout.GetPullRequestBody(ctx)
	if err != nil {
		return nil, err
	}

	prDiff, err := prCheckout.GetPullRequestDiff(ctx, query)
	if err != nil {
		return nil, err
	}

	// TODO: ask LLM

	return nil, nil
}
