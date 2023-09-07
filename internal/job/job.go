package job

import (
	"context"
	"fmt"
	"io"

	"cardboard.package-operator.run/internal/resources"
	"cardboard.package-operator.run/internal/steps"
)

type Job struct {
	Resources []resources.Interface
	Steps     []steps.Interface
}

func (j *Job) Run(ctx context.Context, out io.Writer) error {
	resources := resources.List(j.Resources)
	steps := steps.List(j.Steps)

	// IF continuous
	wch, err := resources.Watch(ctx)
	if err != nil {
		panic(err)
	}

	// run once at start:
	if err := steps.Run(ctx, out); err != nil {
		fmt.Fprintln(out, "--- FAILED ---")
	} else {
		fmt.Fprintln(out, "--- SUCCESS ---")
	}

	// IF continuous
	for evnt := range wch {
		if evnt.Err != nil {
			panic(evnt.Err)
		}
		if err := steps.Run(ctx, out); err != nil {
			fmt.Fprintln(out, "--- FAILED ---")
		} else {
			fmt.Fprintln(out, "--- SUCCESS ---")
		}
	}

	return nil
}
