package ports

import "context"

type VCS interface {
	CloneAndCheckout(ctx context.Context, repoUrl, ref, destDir string) error
}
