package main

import (
	"log/slog"
	"os"

	"github.com/ericbutera/system-design/labs/etl/integrations"
	"github.com/ericbutera/system-design/labs/etl/integrations/etl"
)

// Tenets of ETL:
// data is immutable
// transforms yield new data
// process is idempotent
// prefer passing data location over actual data

// TODO: intermediary step between extract to convert to rows instead of one big array

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	errs := 0
	for _, integration := range integrations.GetIntegrations() {
		if err := run(integration); err != nil { // worker pool
			slog.Error("Error", "err", err)
			errs++
		}
	}

	if errs > 0 {
		os.Exit(1)
	}
}

func run(integration integrations.Integrations) error {
	slog.Info("Running ETL", "integration", integration)
	pipeline, err := etl.New(integration)
	if err != nil {
		return err
	}

	extractResult, err := pipeline.Extract(etl.ExtractParams{})
	if err != nil {
		return err
	}

	transformResult, err := pipeline.Transform(etl.TransformParams{BlobStorage: extractResult.BlobStorage})
	if err != nil {
		return err
	}

	loadResult, err := pipeline.Load(etl.LoadParams{Assets: transformResult.Assets})
	if err != nil {
		return err
	}

	slog.Info("Results", "results", loadResult)

	return nil
}
