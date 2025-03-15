// A generated module for DockerfileOptimizer functions
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
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"dagger/dockerfile-optimizer/internal/dagger"

	"github.com/dustin/go-humanize"
)

type DockerfileOptimizer struct{}

// Build the image from the Dockerfile, returns the number of layers and the size of the image
func imageInfo(ctx context.Context, dir *dagger.Directory, path string) ([]int, error) {
	dirname, filename := filepath.Split(path)
	ctr := dag.Container().
		Build(dir.Directory(dirname), dagger.ContainerBuildOpts{Dockerfile: filename})

	// Mount the OCI image and run tests
	out, err := dag.Container().From("wagoodman/dive:latest").
		WithMountedFile("/tmp/image.tar", ctr.AsTarball(
			// Layer compression seems to cause issues with dive in some cases
			dagger.ContainerAsTarballOpts{ForcedCompression: dagger.ImageLayerCompressionUncompressed},
		)).
		WithMountedDirectory("/workspace", dir.Directory(dirname)).
		WithExec([]string{"dive", "--json", "/tmp/img-info.json", "--ci", "docker-archive:///tmp/image.tar"}).
		File("/tmp/img-info.json").Contents(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to run dive: %w", err)
	}

	var imgInfo struct {
		Layer []struct{} `json:"layer"`
		Image struct {
			SizeBytes int64 `json:"sizeBytes"`
		} `json:"image"`
	}

	if err := json.Unmarshal([]byte(out), &imgInfo); err != nil {
		return nil, fmt.Errorf("failed to parse image info: %w", err)
	}

	numLayers := len(imgInfo.Layer)
	return []int{numLayers, int(imgInfo.Image.SizeBytes)}, nil
}

func askLLM(ws *dagger.Workspace, dockerfile, extraContext string) *dagger.LLM {
	llm := dag.Llm().
		WithWorkspace(ws).
		WithPromptVar("dockerfile", dockerfile).
		WithPromptVar("extra_context", extraContext).
		WithPrompt(`
You are a Platform Engineer with deep knowledge of Dockerfiles. You have access to a workspace.
Use the read, write and build tools to complete the following assignment:

Assignment: Optimize the Dockerfile for reducing its size, number of layers, and build time. And when possible, increasing the security level of the image by implementing best practices.

Follow these guidelines:
- Make all the optimizations you can think of at once, don't try to optimize it step by step.
- Make sure to never downgrade any image version found in the Dockerfile.
- If the Dockerfile is already optimized, just return an explanation that you couldn't optimize it.
- Make sure the new Dockerfile builds correctly and write it to the workspace, replacing the old one.
- Skip explanations of intermediate steps before the final answer.

At the end, return an explanation of the changes you made to the Dockerfile.
$extra_context
`)

	return llm
}

// Optimize a Dockerfile from a directory
func (m *DockerfileOptimizer) optimizeDockerfile(ctx context.Context, src *dagger.Directory) (*dagger.Directory, []string, string, error) {
	// Create a new workspace, using third-party module
	ws := dag.Workspace(src)
	originalWorkdir := ws.Workdir()

	// Find the Dockerfile
	// FIXME: handle multiple Dockerfiles
	dockerfiles, err := ws.Workdir().Glob(ctx, "*Dockerfile*")
	if err != nil {
		return nil, nil, "", fmt.Errorf("cannot read the directory: %w", err)
	}

	if len(dockerfiles) == 0 {
		return nil, nil, "", fmt.Errorf("no Dockerfile found")
	}

	dockerfile := dockerfiles[0]

	// Get the image info
	originalImgInfo, err := imageInfo(ctx, ws.Workdir(), dockerfile)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get image info: %w", err)
	}

	extraContext := ""
	answer := ""
	var lastState *dagger.Workspace
	var lastImgInfo []int
	// Try 5 times to optimize the Dockerfile
	for range make([]int, 5) {
		// Ask the LLM to optimize the Dockerfile
		llm := askLLM(ws, dockerfile, extraContext)
		answer, err = llm.LastReply(ctx)
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to ask LLM: %w", err)
		}

		lastState = llm.Workspace()

		// Compare the optimized Dockerfile with the original one
		lastImgInfo, err = imageInfo(ctx, lastState.Workdir(), dockerfile)
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to get image info: %w", err)
		}

		// We consider the optimization satisfactory if the size of the image is smaller
		if lastImgInfo[1] < originalImgInfo[1] {
			break
		}

		// Otherwise we give extra context to the LLM and try again
		extraContext = "\n\nYou previously attempted to optimize the Dockerfile, but the changes were not satisfactory. Here are the details:\n\n"
		extraContext += fmt.Sprintf("- The number of layers is %d in the original image, and %d layers in the optimized version.\n", originalImgInfo[0], lastImgInfo[0])
		extraContext += fmt.Sprintf("- The original image size is %d bytes, and the optimized image size is %d bytes.\n\n", originalImgInfo[1], lastImgInfo[1])
		extraContext += "Please make the necessary changes to the Dockerfile to improve the image size and number of layers.\n"
		// FIXME: add the modified Dockerfile to the extra context?
	}

	// Check if the workspace has been modified
	diff, err := originalWorkdir.Diff(lastState.Workdir()).Entries(ctx)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get workspace diff: %w", err)
	}

	if len(diff) == 0 {
		return nil, nil, "", fmt.Errorf("failed to optimize the Dockerfile")
	}

	answer += "\n\nImage info:\n"
	answer += fmt.Sprintf("- The original image has %d layers and is %s in size.\n", originalImgInfo[0], humanize.Bytes(uint64(originalImgInfo[1])))
	answer += fmt.Sprintf("- The optimized image has %d layers and is %s in size.\n", lastImgInfo[0], humanize.Bytes(uint64(lastImgInfo[1])))

	return lastState.Workdir(), diff, answer, nil
}

// Create a new PullRequest with the changes in the workspace, the given title and body, returns the PR URL
func createPR(ctx context.Context, githubToken *dagger.Secret, repoURL string, src *dagger.Directory, llmAnswer string) (string, error) {
	// Create a new feature branch
	featureBranch := dag.
		FeatureBranch(githubToken, repoURL, "dockerfile-improvements").
		WithChanges(src).
		Commit("Optimize Dockerfile").
		Push()

	// Make sure changes have been made to the workspace
	diff, err := featureBranch.Diff(ctx, true)
	if err != nil {
		return "", fmt.Errorf("failed to get branch diff: %w", err)
	}

	if diff == "" {
		return "", fmt.Errorf("got empty diff on feature branch (llm did not make any changes)")
	}

	return featureBranch.CreatePullRequestWithLlm(ctx, llmAnswer)
}

// Optimize a Dockerfile from a remote Github repository, and open a PR with the changes
func (m *DockerfileOptimizer) OptimizeDockerfileFromGithub(ctx context.Context, githubToken *dagger.Secret, repoURL string) (string, error) {
	if !strings.HasSuffix(repoURL, ".git") {
		repoURL = repoURL + ".git"
	}

	output, _, answer, err := m.optimizeDockerfile(ctx, dag.Git(repoURL).Head().Tree())
	if err != nil {
		return "", fmt.Errorf("failed to optimize the Dockerfile: %w", err)
	}

	return createPR(ctx, githubToken, repoURL, output, answer)
}

// Optimize a Dockerfile from a directory, only returns the initial directory with the optimized Dockerfile
func (m *DockerfileOptimizer) OptimizeDockerfileFromDirectory(ctx context.Context, src *dagger.Directory) (*dagger.Directory, error) {
	output, _, _, err := m.optimizeDockerfile(ctx, src)
	if err != nil {
		return nil, fmt.Errorf("failed to optimize the Dockerfile: %w", err)
	}

	return output, nil
}
