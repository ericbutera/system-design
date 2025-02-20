// SaaS platform SDK
package sdk

import "log/slog"

type Asset struct {
	VendorID  string
	IPAddress string
	Hostname  string
	OS        string
}

type Platform interface {
	SaveAsset(asset *Asset) error
}

type GrpcAPI struct{}

func NewGrpc() Platform {
	return &GrpcAPI{}
}

func (p *GrpcAPI) SaveAsset(asset *Asset) error {
	slog.Debug("SaaS SDK: Saving asset", "asset", asset)
	return nil
}
