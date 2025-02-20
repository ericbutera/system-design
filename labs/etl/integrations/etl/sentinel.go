package etl

import (
	"log/slog"
)

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
				Integration: e.integration,
				VendorID:    asset.Asset,
				IPAddress:   asset.IPAddress,
				Hostname:    asset.Asset,
			})
		}
		return nil
	})
}
