package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/segmentio/kafka-go"
)

var generators = []func() any{
	func() any { return generateFake[CrowdstrikeAsset]() },
	func() any { return generateFake[MicrosoftDefenderAsset]() },
	func() any { return generateFake[QualysVMAsset]() },
	func() any { return generateFake[Rapid7InsightVMAsset]() },
	func() any { return generateFake[AzureAsset]() },
	func() any { return generateFake[TenableNessusAsset]() },
}

func main() {
	ctx := context.Background()

	broker := os.Getenv("BROKER")
	topic := os.Getenv("TOPIC")
	if broker == "" || topic == "" {
		log.Fatal("BROKER and TOPIC environment variables must be set")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
	}

	for range 100 {
		for _, generator := range generators {
			write(ctx, writer, generator())
		}
		time.Sleep(1 * time.Second)
	}
}

func write(ctx context.Context, writer *kafka.Writer, data any) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return err
	}

	slog.Info("Sending message", "json", string(encoded))
	err = writer.WriteMessages(ctx, kafka.Message{Value: encoded})
	if err != nil {
		log.Printf("Error producing message: %v", err)
		return err
	}
	return nil
}

func generateFake[T any]() T {
	var asset T
	faker.FakeData(&asset)
	return asset
}

type CrowdstrikeAsset struct {
	DeviceID string `json:"device_id" example:"cs-205" faker:"uuid_hyphenated"`
	Hostname string `json:"hostname" example:"test-vm" faker:"domain_name"`
	IP       string `json:"ip" example:"10.0.0.01" faker:"ipv4"`
	OS       string `json:"os" example:"Debian 10" faker:"word"`
	Status   string `json:"status" example:"Offline" faker:"word"`
}

type MicrosoftDefenderAsset struct {
	DeviceID        string   `json:"device_id" example:"cs-205" faker:"uuid_hyphenated"`
	Hostname        string   `json:"hostname" example:"test-vm" faker:"domain_name"`
	IP              string   `json:"ip" example:"10.0.0.01" faker:"ipv4"`
	OS              string   `json:"os" example:"Debian 10" faker:"word"`
	RiskLevel       string   `json:"risk_level" example:"High" faker:"word"`
	LastSeen        string   `json:"last_seen" example:"2021-01-01T00:00:00Z" faker:"timestamp"`
	ThreatsDetected []string `json:"threats_detected" example:"Trojan.Win32.Generic,Exploit.CVE-2023-4567" faker:"[]word"`
}

type QualysVMAsset struct {
	AssetID string `json:"asset_id" example:"qualys-205" faker:"uuid_hyphenated"`
	FQDN    string `json:"fqdn" example:"mail.example.com" faker:"domain_name"`
	IP      string `json:"ip" example:"10.10.10.01" faker:"ipv4"`
	OS      string `json:"os" example:"Debian 10" faker:"word"`
}

type Rapid7InsightVMAsset struct {
	AssetID   string  `json:"asset_id" example:"r7-101" faker:"uuid_hyphenated"`
	Hostname  string  `json:"hostname" example:"host-1" faker:"domain_name"`
	IPAddress string  `json:"ip_address" example:"10.0.0.1" faker:"ipv4"`
	OS        string  `json:"os" example:"Windows" faker:"word"`
	RiskScore float64 `json:"risk_score" example:"400.5" faker:"amount"`
	LastScan  string  `json:"last_scan" example:"2021-01-01T00:00:00Z" faker:"timestamp"`
}

type AzureAsset struct {
	IncidentID   string       `json:"incident_id" example:"azsent-501" faker:"uuid_hyphenated"`
	Asset        string       `json:"asset" example:"vm-web-01" faker:"word"`
	IPAddress    string       `json:"ip_address" example:"10.0.0.1" faker:"ipv4"`
	ThreatLevel  string       `json:"threat_level" example:"High" faker:"word"`
	LastDetected string       `json:"last_detected" example:"2025-02-19T10:30:00Z" faker:"timestamp"`
	Alerts       []AzureAlert `json:"alerts" faker:"-"`
}

type TenableNessusAsset struct {
	ID              string                       `json:"id" example:"ten-003" faker:"uuid_hyphenated"`
	Hostname        string                       `json:"hostname" example:"appserver.internal" faker:"domain_name"`
	IPAddress       string                       `json:"ip" example:"10.0.0.1" faker:"ipv4"`
	OS              string                       `json:"operating_system" example:"RHEL 8" faker:"word"`
	MacAddress      string                       `json:"mac_address" example:"00:0a:95:9d:68:16" faker:"mac_address"`
	LastSeen        string                       `json:"last_seen" example:"2025-02-19T10:30:00Z" faker:"timestamp"`
	Vulnerabilities []TenableNessusVulnerability `json:"vulnerabilities" faker:"-"`
}

type TenableNessusVulnerability struct {
	ID                  string  `json:"id" example:"vul-001" faker:"uuid_hyphenated"`
	Severity            string  `json:"severity" example:"High" faker:"word"`
	ExploitabilityScore float64 `json:"exploitability_score" example:"8.6" faker:"amount"`
}

type AzureAlert struct {
	AlertID   string `json:"alert_id" example:"alert-001" faker:"uuid_hyphenated"`
	AlertName string `json:"alert_name" example:"Suspicious Activity" faker:"word"`
}
