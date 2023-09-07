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

	// Build Container
	// - compile
	// - (assemble image contents)
	// - build image
	// - push image

	logger := stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))
	stdr.SetVerbosity(internal.Debug)

	ctx := internal.SetupSignalHandler()
	ctx = logr.NewContext(ctx, logger)

	if err := testjob.Run(ctx, os.Stdout); err != nil {
		panic(err)
	}
}
