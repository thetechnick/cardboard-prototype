package steps

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os/exec"
	"syscall"
	"time"

	"cardboard.package-operator.run/internal"
)

// RunStep runs an executable at the given path.
type run struct {
	internal.Named
	Path string
	Args []string

	Inputs  []Input
	Outputs []Output
}

const (
	WorkDirInputName = "_WorkDir"
	ImageInputName   = "_Image"
)

type Input struct {
	internal.Named
	Path string // optional
}

type Output struct {
	internal.Named
	Path string // optional
}

func Run(name, path string, args []string) *run {
	return &run{
		Named: internal.Named(name),
		Path:  path,
		Args:  args,
	}
}

func (r *run) Run(ctx context.Context, out io.Writer) error {
	err := r.run(ctx, out)

	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		if pathErr.Err == syscall.ETXTBSY {
			// retry if file is busy, which might happen when it has just been written.
			time.Sleep(100 * time.Millisecond)
			return r.run(ctx, out)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *run) run(ctx context.Context, out io.Writer) error {
	cmd := exec.CommandContext(ctx, r.Path, r.Args...)
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}
