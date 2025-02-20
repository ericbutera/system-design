package etl

type QualysVMEtl struct{ BaseEtl }

type QualysVMAsset struct {
	AssetID string `json:"asset_id" example:"qualys-205"`
	FQDN    string `json:"fqdn" example:"mail.example.com"`
	IP      string `json:"ip" example:"10.10.10.01"`
	OS      string `json:"os" example:"Debian 10"`
}

func (e *QualysVMEtl) Transform(params TransformParams) (TransformResult, error) {
	return Transformer(params, func(data []QualysVMAsset, assets *[]Asset) error {
		for _, asset := range data {
			*assets = append(*assets, Asset{
				Integration: e.integration,
				VendorID:    asset.AssetID,
				Hostname:    asset.FQDN,
				IPAddress:   asset.IP,
				OS:          asset.OS,
			})
		}
		return nil
	})
}
