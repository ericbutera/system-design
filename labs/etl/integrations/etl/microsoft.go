package etl

type MicrosoftDefenderEtl struct{ BaseEtl }

type MicrosoftDefenderAsset struct {
	DeviceID        string   `json:"device_id" example:"cs-205"`
	Hostname        string   `json:"hostname" example:"test-vm"`
	IP              string   `json:"ip" example:"10.0.0.01"`
	OS              string   `json:"os" example:"Debian 10"`
	RiskLevel       string   `json:"risk_level" example:"High"`
	LastSeen        string   `json:"last_seen" example:"2021-01-01T00:00:00Z"`
	ThreatsDetected []string `json:"threats_detected" example:"Trojan.Win32.Generic,Exploit.CVE-2023-4567"`
}

func (e *MicrosoftDefenderEtl) Transform(params TransformParams) (TransformResult, error) {
	return Transformer(params, func(data []MicrosoftDefenderAsset, assets *[]Asset) error {
		for _, asset := range data {
			*assets = append(*assets, Asset{
				Integration: e.integration,
				VendorID:    asset.DeviceID,
				IPAddress:   asset.IP,
				Hostname:    asset.Hostname,
				OS:          asset.OS,
				// TODO: os, risk_level, last_seen, threats_detected
			})
		}
		return nil
	})
}
