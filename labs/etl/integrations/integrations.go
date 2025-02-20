package integrations

type Integrations string

const CrowdstrikeFalcon Integrations = "crowdstrike_falcon"
const AzureSentinel Integrations = "azure_sentinel"
const MicrosoftDefender Integrations = "microsoft_defender"
const QualysVM Integrations = "qualys_vm"
const Rapid7InsightVM Integrations = "rapid7_insightvm"
const TenableNessus Integrations = "tenable_nessus"

func GetIntegrations() []Integrations {
	return []Integrations{
		AzureSentinel,
		CrowdstrikeFalcon,
		MicrosoftDefender,
		QualysVM,
		Rapid7InsightVM,
		TenableNessus,
	}
}
