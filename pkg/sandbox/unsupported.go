//go:build !darwin

package sandbox

import (
	"context"
	"errors"
)

func Serve(ctx context.Context, opts Options) error {
	return errors.New("the sandbox data plane is only supported on darwin guests")
}
