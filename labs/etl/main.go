package main

import (
	"log/slog"
	"os"

	"github.com/ericbutera/system-design/labs/etl/integrations"
	"github.com/ericbutera/system-design/labs/etl/integrations/etl"
	sdk "github.com/ericbutera/system-design/labs/etl/saas"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	client := sdk.NewGrpc()

	errs := 0
	for _, integration := range integrations.GetIntegrations() {
		if err := run(integration, client); err != nil { // worker pool
			slog.Error("Error", "err", err)
			errs++
		}
	}

	if errs > 0 {
		slog.Info("Errors occurred", "count", errs)
		os.Exit(1)
	}
}

func run(integration integrations.Integrations, saasClient sdk.Platform) error {
	slog.Info("Running ETL", "integration", integration)
	pipeline, err := etl.New(integration, saasClient)
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
