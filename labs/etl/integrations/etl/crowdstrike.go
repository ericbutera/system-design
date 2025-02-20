package etl

type CrowdstrikeFalconEtl struct{ BaseEtl }

type CrowdstrikeAsset struct {
	DeviceID string `json:"device_id" example:"cs-205"`
	Hostname string `json:"hostname" example:"test-vm"`
	IP       string `json:"ip" example:"10.0.0.01"`
	OS       string `json:"os" example:"Debian 10"`
	Status   string `json:"status" example:"Offline"`
}

func (e *CrowdstrikeFalconEtl) Transform(params TransformParams) (TransformResult, error) {
	return Transformer(params, func(data []CrowdstrikeAsset, assets *[]Asset) error {
		for _, asset := range data {
			*assets = append(*assets, Asset{
				Integration: e.integration,
				VendorID:    asset.DeviceID,
				IPAddress:   asset.IP,
				Hostname:    asset.Hostname,
				OS:          asset.OS,
			})
		}
		return nil
	})
}
