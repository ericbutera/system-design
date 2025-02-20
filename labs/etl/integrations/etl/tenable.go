package etl

import "log/slog"

type TenableNessusEtl struct{ BaseEtl }

type TenableNessusAsset struct {
	ID              string                       `json:"id" example:"ten-003"`
	Hostname        string                       `json:"hostname" example:"appserver.internal"`
	IPAddress       string                       `json:"ip" example:"10.0.0.1"`
	OS              string                       `json:"operating_system" example:"RHEL 8"`
	MacAddress      string                       `json:"mac_address" example:"00:0a:95:9d:68:16"`
	LastSeen        string                       `json:"last_seen" example:"2025-02-19T10:30:00Z"`
	Vulnerabilities []TenableNessusVulnerability `json:"vulnerabilities"`
}

type TenableNessusVulnerability struct {
	ID                  string  `json:"id" example:"vul-001"`
	Severity            string  `json:"severity" example:"High"`
	ExploitabilityScore float64 `json:"exploitability_score" example:"8.6"`
}

func (e *TenableNessusEtl) Transform(params TransformParams) (TransformResult, error) {
	return Transformer(params, func(data []TenableNessusAsset, assets *[]Asset) error {
		slog.Info("Transforming data", "data", data)
		for _, asset := range data {
			*assets = append(*assets, Asset{
				Integration: e.integration,
				VendorID:    asset.ID,
				IPAddress:   asset.IPAddress,
				Hostname:    asset.Hostname,
				OS:          asset.OS,
			})
		}
		return nil
	})
}
