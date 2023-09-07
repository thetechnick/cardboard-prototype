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
	for _, step := range j.Steps {
		fmt.Fprintf(out, "--- STEP %s ---\n", step.Name())
		if err := step.Run(ctx, out); err != nil {
			fmt.Fprintf(out, "%v\n", err)
			fmt.Fprintln(out, "--- FAILED ---")
			return nil
		}
	}
	fmt.Fprintln(out, "--- SUCCESS ---")
	return nil
}

func (j *Job) Watch(ctx context.Context, out io.Writer) error {
	resources := resources.List(j.Resources)

	// IF continuous
	wch, err := resources.Watch(ctx)
	if err != nil {
		panic(err)
	}

	// run once at start
	if err := j.Run(ctx, out); err != nil {
		return err
	}

	// IF continuous
	for evnt := range wch {
		if evnt.Err != nil {
			panic(evnt.Err)
		}
		if err := j.Run(ctx, out); err != nil {
			return err
		}
	}
	return nil
}
