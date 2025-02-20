package etl

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ericbutera/system-design/labs/etl/integrations"
	i "github.com/ericbutera/system-design/labs/etl/integrations"
	sdk "github.com/ericbutera/system-design/labs/etl/saas"
)

type Asset struct {
	Integration integrations.Integrations
	VendorID    string // Represents the unique identifier of the asset in the source system. This should be the correlation ID between "theirs" and "ours"
	IPAddress   string
	Hostname    string
	OS          string // Operating System
}

type ExtractParams struct {
}

type ExtractResult struct {
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

func New(integration i.Integrations, saasClient sdk.Platform) (Etl, error) {
	base := BaseEtl{
		integration: integration,
		saasClient:  saasClient,
	}

	switch integration {
	case i.AzureSentinel:
		return &AzureSentinelEtl{BaseEtl: base}, nil
	case i.CrowdstrikeFalcon:
		return &CrowdstrikeFalconEtl{BaseEtl: base}, nil
	case i.MicrosoftDefender:
		return &MicrosoftDefenderEtl{BaseEtl: base}, nil
	case i.QualysVM:
		return &QualysVMEtl{BaseEtl: base}, nil
	case i.Rapid7InsightVM:
		return &Rapid7InsightVMEtl{BaseEtl: base}, nil
	case i.TenableNessus:
		return &TenableNessusEtl{BaseEtl: base}, nil
	}

	return nil, ErrInvalidIntegration
}

type BaseEtl struct {
	integration integrations.Integrations
	saasClient  sdk.Platform
}

func (b *BaseEtl) Extract(params ExtractParams) (ExtractResult, error) {
	slog.Debug("Extracting data", "integration", b.integration, "params", params)
	// note: this would normally fetch & write data (extract)
	// it should support rate limiting, retries, backoff, auth, pagination
	// also, extract should be durable so it doesn't access the remote source again

	path, err := FetchData(b.integration)
	if err != nil {
		return ExtractResult{}, err
	}

	return ExtractResult{
		BlobStorage: path,
	}, nil
}

func (b *BaseEtl) Transform(params TransformParams) (TransformResult, error) {
	panic("transform not implemented for " + b.integration)
}

func (b *BaseEtl) Load(params LoadParams) (LoadResult, error) {
	// note: this would normally write data to a warehouse
	for _, asset := range params.Assets {
		_ = b.saasClient.SaveAsset(&sdk.Asset{
			VendorID:  asset.VendorID,
			IPAddress: asset.IPAddress,
			Hostname:  asset.Hostname,
			OS:        asset.OS,
		})
	}
	return LoadResult{}, nil
}

// handles converting raw integration data back into a typed form
// use fn callback to pass in a transform mapping function
func Transformer[T any](params TransformParams, fn func(data T, assets *[]Asset) error) (TransformResult, error) {
	slog.Debug("Transformer", "params", params)

	// note: it can be a good idea to convert raw data into a more structured and easier to parse format like Avro
	blobs, err := GetBlobs(params.BlobStorage)
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

// (simulate) fetch data from source
func FetchData(integration integrations.Integrations) (string, error) {
	slog.Debug("Fetching data", "integration", integration)
	// TODO: write source data to cloud storage: /integration/{integration}/{date}/{chunk}.json
	// pass data reference to data between stages
	return GetPath(integration), nil
}

const Pattern = `testdata/%s.json`

// object storage path for integration data
func GetPath(integration integrations.Integrations) string {
	return fmt.Sprintf(Pattern, integration)
}

type Blobs [][]byte

// (simulate) read from object storage
func GetBlobs(path string) (Blobs, error) {
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

func SaveAsset(asset Asset) {
	slog.Debug("Saving asset", "asset", asset)
}
