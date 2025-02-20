package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// Tenets of ETL:
// data is immutable
// transforms yield new data
// process is idempotent
// prefer passing data location over actual data

// TODO: intermediary step between extract to convert to rows instead of one big array

const Pattern = `testdata/%s.json`

type Integrations string

const IntegrationCrowdstrikeFalcon Integrations = "crowdstrike_falcon"
const IntegrationAzureSentinel Integrations = "azure_sentinel"
const IntegrationMicrosoftDefender Integrations = "microsoft_defender"
const IntegrationNessusStandalone Integrations = "nessus_standalone"
const IntegrationQualysVM Integrations = "qualys_vm"
const IntegrationRapid7InsightVM Integrations = "rapid7_insight_vm"
const IntegrationTenableNessus Integrations = "tenable_nessus"

type Asset struct {
	VendorID    string
	Integration Integrations // enum
	IPAddress   string
	Hostname    string
}

type ExtractParams struct {
	//Integration Integrations
}

type ExtractResult struct {
	//Integration Integrations
	BlobStorage string // cloud storage path
}

type TransformParams struct {
	BlobStorage string
}

type TransformResult struct {
	Assets []Asset
	Errors int
}

type LoadParams struct {
	Assets []Asset
}

type LoadResult struct {
	Success int
	Errors  int
}

type Etl interface {
	Extract(params ExtractParams) (ExtractResult, error)
	Transform(params TransformParams) (TransformResult, error)
	Load(params LoadParams) (LoadResult, error)
}

var (
	ErrInvalidIntegration = errors.New("invalid integration")
)

func getIntegrations() []Integrations {
	return []Integrations{
		IntegrationAzureSentinel,
		IntegrationCrowdstrikeFalcon,
		// 	IntegrationMicrosoftDefender,
		// 	IntegrationNessusStandalone,
		// 	IntegrationQualysVM,
		// 	IntegrationRapid7InsightVM,
		// 	IntegrationTenableNessus,
	}
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	errs := 0
	for _, integration := range getIntegrations() {
		if err := run(integration); err != nil { // worker pool
			slog.Error("Error", "err", err)
			errs++
		}
	}

	if errs > 0 {
		os.Exit(1)
	}
}

func run(integration Integrations) error {
	slog.Info("Running ETL", "integration", integration)
	etl, err := New(integration)
	if err != nil {
		return err
	}

	extractResult, err := etl.Extract(ExtractParams{})
	if err != nil {
		return err
	}

	transformResult, err := etl.Transform(TransformParams{BlobStorage: extractResult.BlobStorage})
	if err != nil {
		return err
	}

	loadResult, err := etl.Load(LoadParams{Assets: transformResult.Assets})
	if err != nil {
		return err
	}

	slog.Info("Results", "results", loadResult)

	return nil
}

type BaseEtl struct {
	integration Integrations
}

func (b *BaseEtl) Extract(params ExtractParams) (ExtractResult, error) {
	slog.Debug("Extracting data", "integration", b.integration, "params", params)
	// note: this would normally fetch & write data (extract)
	// it should support rate limiting, retries, backoff, auth, pagination
	// also, extract should be durable so it doesn't access the remote source again

	path, err := fetchData(b.integration)
	if err != nil {
		return ExtractResult{}, err
	}

	return ExtractResult{
		BlobStorage: path,
	}, nil
}

func (b *BaseEtl) Transform(params TransformParams) (TransformResult, error) {
	panic("not implemented")
}

func (b *BaseEtl) Load(params LoadParams) (LoadResult, error) {
	// note: this would normally write data to a warehouse
	for _, asset := range params.Assets {
		slog.Info("Loading asset", "asset", asset)
	}
	return LoadResult{}, nil
}

func New(integration Integrations) (Etl, error) {
	base := BaseEtl{integration: integration}

	switch integration {
	case IntegrationAzureSentinel:
		return &AzureSentinelEtl{BaseEtl: base}, nil
	case IntegrationCrowdstrikeFalcon:
		return &CrowdstrikeFalconEtl{BaseEtl: base}, nil
	case IntegrationMicrosoftDefender:
		return &MicrosoftDefenderEtl{BaseEtl: base}, nil
	case IntegrationNessusStandalone:
		return &NessusStandaloneEtl{BaseEtl: base}, nil
	case IntegrationQualysVM:
		return &QualysVMEtl{BaseEtl: base}, nil
	case IntegrationRapid7InsightVM:
		return &Rapid7InsightVMEtl{BaseEtl: base}, nil
	case IntegrationTenableNessus:
		return &TenableNessusEtl{BaseEtl: base}, nil
	}

	return nil, ErrInvalidIntegration
}

type AzureSentinelEtl struct{ BaseEtl }

type AzureAsset struct {
	IncidentID   string       `json:"incident_id" example:"azsent-501"`
	Asset        string       `json:"asset" example:"vm-web-01"`
	IPAddress    string       `json:"ip_address" example:"10.0.0.1"`
	ThreatLevel  string       `json:"threat_level" example:"High"`
	LastDetected string       `json:"last_detected" example:"2025-02-19T10:30:00Z"`
	Alerts       []AzureAlert `json:"alerts"`
}

type AzureAlert struct {
	AlertID     string `json:"alert_id" example:"AL-1001"`
	Description string `json:"description" example:"Brute force attack detected"`
	Severity    string `json:"severity" example:"High"`
}

func (e *AzureSentinelEtl) Transform(params TransformParams) (TransformResult, error) {
	return Transformer(params, func(data []AzureAsset, assets *[]Asset) error {
		slog.Info("Transforming data", "data", data)
		for _, asset := range data {
			*assets = append(*assets, Asset{
				VendorID:    asset.IncidentID,
				Integration: e.integration,
				IPAddress:   asset.IPAddress,
				Hostname:    asset.Asset,
			})
		}
		return nil
	})
}

type CrowdstrikeAsset struct {
	DeviceID string `json:"device_id" example:"cs-205"`
	Hostname string `json:"hostname" example:"test-vm"`
	IP       string `json:"ip" example:"10.0.0.01"`
	OS       string `json:"os" example:"Debian 10"`
	Status   string `json:"status" example:"Offline"`
}

type CrowdstrikeFalconEtl struct{ BaseEtl }

func (e *CrowdstrikeFalconEtl) Transform(params TransformParams) (TransformResult, error) {
	return Transformer(params, func(data []CrowdstrikeAsset, assets *[]Asset) error {
		for _, asset := range data {
			*assets = append(*assets, Asset{
				VendorID:    asset.DeviceID,
				Integration: e.integration,
				IPAddress:   asset.IP,
				Hostname:    asset.Hostname,
			})
		}
		return nil
	})
}

type MicrosoftDefenderEtl struct{ BaseEtl }

type NessusStandaloneEtl struct{ BaseEtl }

type QualysVMEtl struct{ BaseEtl }

type Rapid7InsightVMEtl struct{ BaseEtl }

type TenableNessusEtl struct{ BaseEtl }

func Transformer[T any](params TransformParams, fn func(data T, assets *[]Asset) error) (TransformResult, error) {
	slog.Debug("Transformer", "params", params)

	// note: it can be a good idea to convert raw data into a more structured and easier to parse format like Avro
	blobs, err := getBlobs(params.BlobStorage)
	if err != nil {
		return TransformResult{}, err
	}

	var (
		assets []Asset
		errs   int
	)

	for _, blob := range blobs {
		var integrationAssets T
		if err := json.Unmarshal(blob, &integrationAssets); err != nil {
			slog.Error("Error unmarshalling", "err", err)
			errs++
			continue
		}

		if err := fn(integrationAssets, &assets); err != nil {
			slog.Error("Error processing blob", "err", err)
			errs++
		}
	}

	return TransformResult{
		Assets: assets,
		Errors: errs,
	}, nil
}

// fetch data from source
func fetchData(integration Integrations) (string, error) {
	slog.Debug("Fetching data", "integration", integration)
	// TODO: write source to cloud storage: /integration/{integration}/{date}/{chunk}.json
	return getPath(integration), nil
}

// object storage path for integration data
func getPath(integration Integrations) string {
	return fmt.Sprintf(Pattern, integration)
}

type Blobs [][]byte

func getBlobs(path string) (Blobs, error) {
	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	blobs := make(Blobs, 0)
	for _, file := range files {
		slog.Debug("reading file", "file", file, "path", path)
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, data)
	}
	return blobs, nil
}
