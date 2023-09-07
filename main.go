package main

import (
	"log"
	"os"

	"cardboard.package-operator.run/internal"
	"cardboard.package-operator.run/internal/job"
	"cardboard.package-operator.run/internal/resources"
	"cardboard.package-operator.run/internal/steps"
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

func main() {
	// Test Job
	testjob := job.Job{
		Resources: []resources.Interface{
			resources.Folder("src", "."),
		},
		Steps: []steps.Interface{
			steps.Run("unit-test", "./hack/scripts/test.sh", nil),
		},
	}
	var _ = testjob

	src := resources.Folder("src", ".")

	// BuildJob
	buildJob := job.Job{
		Resources: []resources.Interface{
			src,
		},
		Steps: []steps.Interface{
			steps.Run("go-build", "go", []string{"build", "-v", "-o", "bin/cardboard", "./cmd/cardboard"}),
			// TODO: reusable thing?
			steps.Run("image-build", "podman", []string{"build", "-f", "config/images/cardboard.Containerfile", "-o", "bin/cardboard.tar.gz", "."}),
			// step.Put? src.PutStep
		},
	}
	// Build Container
	// - compile
	// - (assemble image contents)
	// - build image
	// - push image

	logger := stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))
	stdr.SetVerbosity(internal.Debug)

	ctx := internal.SetupSignalHandler()
	ctx = logr.NewContext(ctx, logger)

	if err := buildJob.Run(ctx, os.Stdout); err != nil {
		panic(err)
	}

	// if err := testjob.Run(ctx, os.Stdout); err != nil {
	// 	panic(err)
	// }
}
