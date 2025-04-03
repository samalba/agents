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
	"fmt"
)

type PrReviewer struct{}

// Review a PR, by default it will review the PR description and the diff.
// Query is any argument supported by the gh cli (gh pr view [<number> | <url> | <branch>]).
// Additional instructions can be provided to the LLM to guide the review.
// Returns the URL of the PR comment created by the LLM
func (m *PrReviewer) ReviewPr(ctx context.Context, githubToken *dagger.Secret, query, repoURL,
	// +optional
	additionalInstructions string,
) (string, error) {
	prCheckout := dag.FeatureBranch(githubToken, repoURL, "unused-branch-name").
		CheckoutPullRequest(query)

	prInfo, err := prCheckout.GetPullRequestBodyTitle(ctx)
	if err != nil {
		return "", err
	}

	prTitle := prInfo[0]
	prBody := prInfo[1]

	prDiff, err := prCheckout.GetPullRequestDiff(ctx, query)
	if err != nil {
		return "", err
	}

	if additionalInstructions == "" {
		additionalInstructions = fmt.Sprintf("\n\nAdditional Instructions: %s\n", additionalInstructions)
	}

	env := dag.Env().
		WithStringInput("prTitle", prTitle, "The title of the PullRequest to review").
		WithStringInput("prBody", prBody, "The body of the PullRequest to review").
		WithStringInput("prDiff", prDiff, "The diff of the PullRequest to review").
		WithStringInput("additionalInstructions", additionalInstructions, "Optional additional context for the PR review (can be empty)")

	llm := dag.LLM().
		WithEnv(env).
		WithPrompt(`Review the following Pull Request:

PR Title:
$prTitle

PR Body:
$prBody

PR Diff:
$prDiff
$additionalInstructions

Generate a review of the Pull Request. Include the following information:
- The changes made to the code
- The rationale for the changes
- Any potential risks or considerations
- Any other relevant details

In the review, make a recommendation for merging the PR or requesting changes,
but do not repeat the PR title or body, or summarizing the changes, focus on the
merge recommendation and assessment of the changes.

At the very end of the message, mentions if you recommends merging the PR or requesting changes, in bold, with a corresponding emoji.

Only output the review, nothing else.`)

	review, err := llm.LastReply(ctx)
	if err != nil {
		return "", err
	}

	// Add the review as a comment
	url, err := prCheckout.AddPullRequestComment(ctx, review, false)
	if err != nil {
		return "", err
	}

	return url, nil
}
