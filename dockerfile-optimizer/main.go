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

func askLLM(ws *dagger.Workspace, dockerfile, extraContext string) *dagger.Llm {
	llm := dag.Llm().
		WithWorkspace(ws).
		WithPromptVar("dockerfile", dockerfile).
		WithPromptVar("extra_context", extraContext).
		WithPrompt(`
You are a Platform Engineer with deep knowledge of Dockerfiles. You have access to a workspace.
Use the read, write and build tools to complete the following assignment.

- Build the Dockerfile in the provided workspace at the path: "$dockerfile"
- Optimize the Dockerfile for reducing its size, number of layers, and build time. And if possible, increasing the security level of the image by implementing best practices.
- Ensure to not downgrade any version found in the Dockerfile.
- If the Dockerfile is already optimized, just return an explanation that you couldn't optimize it.
- Make sure the Dockerfile builds correctly before going to the next step.
- If you have changes to make, write the optimized Dockerfile to the workspace, at the same path "$dockerfile".

At the end, return an explanation of the changes you made to the Dockerfile.
$extra_context
`)

	return llm
}

// Optimize a Dockerfile
func (m *DockerfileOptimizer) OptimizeDockerfile(ctx context.Context, githubToken *dagger.Secret, repoURL string) (string, error) {
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// Create a new workspace, using third-party module
	ws := dag.Workspace(githubToken, repoURL)
	originalWorkdir := ws.Workdir()

	// Find the Dockerfile
	// FIXME: handle multiple Dockerfiles
	dockerfiles, err := ws.Workdir().Glob(ctx, "*Dockerfile*")
	if err != nil {
		return "", fmt.Errorf("cannot read the directory: %w", err)
	}

	if len(dockerfiles) == 0 {
		return "", fmt.Errorf("no Dockerfile found")
	}

	dockerfile := dockerfiles[0]

	// Get the image info
	originalImgInfo, err := imageInfo(ctx, ws.Workdir(), dockerfile)
	if err != nil {
		return "", fmt.Errorf("failed to get image info: %w", err)
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
			return "", fmt.Errorf("failed to ask LLM: %w", err)
		}

		lastState = llm.Workspace()

		// Compare the optimized Dockerfile with the original one
		lastImgInfo, err = imageInfo(ctx, lastState.Workdir(), dockerfile)
		if err != nil {
			return "", fmt.Errorf("failed to get image info: %w", err)
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
		return "", fmt.Errorf("failed to get workspace diff: %w", err)
	}

	// DEBUG
	dag.Container().From("cgr.dev/chainguard/wolfi-base:latest").
		WithMountedDirectory("/orig_workspace", originalWorkdir).
		WithMountedDirectory("/workspace", lastState.Workdir()).
		Terminal().Sync(ctx)

	if len(diff) == 0 {
		return answer, fmt.Errorf("failed to optimize the Dockerfile")
	}

	answer += "\n\nImage info:\n"
	answer += fmt.Sprintf("- The original image has %d layers and is %d bytes in size.\n", originalImgInfo[0], originalImgInfo[1])
	answer += fmt.Sprintf("- The optimized image has %d layers and is %d bytes in size.\n", lastImgInfo[0], lastImgInfo[1])

	return answer, nil
}

// // Create a new PullRequest with the changes in the workspace, the given title and body, returns the PR URL
// func (w *Workspace) CreatePR(ctx context.Context, title, body string) (string, error) {
// 	// generate a random branch name
// 	branchName := "dockerfile-improvements-" + uuid.New().String()[:8]
// 	// The changeset needs to contain only the Dockerfile otherwise the diff will fail (FIXME?)
// 	changeset := dag.Directory().WithFile("Dockerfile", w.Container.File("Dockerfile"))
// 	// Create a new feature branch
// 	featureBranch := dag.FeatureBranch(w.GitHubToken, w.RepoURL, branchName).
// 		WithChanges(changeset)

// 	// Make sure changes have been made to the workspace
// 	diff, err := featureBranch.Diff(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	if diff == "" {
// 		return "", fmt.Errorf("got empty diff on feature branch (llm did not make any changes)")
// 	}

// 	return featureBranch.PullRequest(ctx, title, body)
// }
