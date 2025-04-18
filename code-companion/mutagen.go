package main

import (
	"context"
	"dagger/host-sync/internal/dagger"
	"errors"
)

// Creates a Mutagen agent that can be used to sync with a host
//
// $ dagger call mutagen-agent --authorized_keys ~/.ssh/id_ed25519.pub up --ports 1222:22
// $ mutagen sync create --name=MyCode ./my-code root@localhost:1222:~/dagger
// $ mutagen sync monitor MyCode
func (m *CodeCompanion) MutagenAgent(
	ctx context.Context,
	// +optional
	authorizedKeys *dagger.File,
	// +optional
	publicKey string,
) (*dagger.Service, error) {

	if publicKey != "" {
		authorizedKeys = dag.
			Directory().
			WithNewFile("authorized_keys", publicKey).
			File("authorized_keys")
	}
	if authorizedKeys == nil {
		return nil, errors.New("authorizedKeys or publicKey must be provided")
	}

	return dag.Container().
		From("hermsi/alpine-sshd").
		WithEnvVariable("ROOT_KEYPAIR_LOGIN_ENABLED", "true").
		WithFile("/root/.ssh/authorized_keys", authorizedKeys, dagger.ContainerWithFileOpts{
			Permissions: 0600,
		}).
		// Ignore the error with chown (unknown user root.root)
		WithExec([]string{"sh", "-c", "sed -i 's/set -e//' /entrypoint.sh"}).
		WithMountedCache("/root/dagger", dag.CacheVolume("MyCode"), dagger.ContainerWithMountedCacheOpts{
			Sharing: dagger.CacheSharingModeShared,
		}).
		WithExposedPort(22).
		AsService(dagger.ContainerAsServiceOpts{
			UseEntrypoint: true,
		}), nil
}
