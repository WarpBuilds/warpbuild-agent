//go:build !darwin

package sandbox

import (
	"context"
	"errors"
)

// The sandbox data plane is macOS-only: the transport is AF_VSOCK into a guest
// under Virtualization.framework, and the pause path speaks diskutil.
func Serve(ctx context.Context, opts Options) error {
	return errors.New("the sandbox data plane is only supported on darwin guests")
}
