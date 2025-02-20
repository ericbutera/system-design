package etl

type Rapid7InsightVMEtl struct{ BaseEtl }

type Rapid7InsightVMAsset struct {
	AssetID   string  `json:"asset_id" example:"r7-101"`
	Hostname  string  `json:"hostname" example:"host-1"`
	IPAddress string  `json:"ip_address" example:"10.0.0.1"`
	OS        string  `json:"os" example:"Windows"`
	RiskScore float64 `json:"risk_score" example:"400.5"`
	LastScan  string  `json:"last_scan" example:"2021-01-01T00:00:00Z"`
}

func (e *Rapid7InsightVMEtl) Transform(params TransformParams) (TransformResult, error) {
	return Transformer(params, func(data []Rapid7InsightVMAsset, assets *[]Asset) error {
		for _, asset := range data {
			*assets = append(*assets, Asset{
				Integration: e.integration,
				VendorID:    asset.AssetID,
				IPAddress:   asset.IPAddress,
				Hostname:    asset.Hostname,
				OS:          asset.OS,
			})
		}
		return nil
	})
}
