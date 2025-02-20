package etl

// Azure SDK:
// "github.com/Azure/azure-sdk-for-go/sdk/azcore"
// "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
// "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights"

type AzureSentinelEtl struct{ BaseEtl }

/*
Example implementation of data extraction from Azure Sentinel:
- heartbeat page data for resuming on error
- concurrency on page requests

func (e *AzureSentinelEtl) Load(params LoadParams) (LoadResult, error) {
	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights
    client, _ := armsecurityinsights.NewIncidentsClient(subscription, cred, nil)
    pager := client.NewListPager(group, workspace, nil)
    for pager.More() {
        page, _ := pager.NextPage(ctx)
        for _, incident := range page.Value {
			// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights#Incident
            fmt.Printf("Incident ID: %s\n", *incident.ID)
        }
    }
}
*/

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
