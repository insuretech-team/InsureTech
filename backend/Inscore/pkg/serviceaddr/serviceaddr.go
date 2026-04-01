package serviceaddr

import (
	"fmt"
	"os"
	"strings"

	"github.com/newage-saint/insuretech/ops/env"
)

var isRunningInDocker = env.IsRunningInDocker

type Ports struct {
	Grpc int `yaml:"grpc"`
	Http int `yaml:"http"`
}

type Service struct {
	Name  string `yaml:"name"`
	Ports Ports  `yaml:"ports"`
}

type ServicesConfig struct {
	Services map[string]Service `yaml:"services"`
}

func ResolveGRPCAddr(explicit, overrideHost, serviceKey string, grpcPort int) string {
	if v := strings.TrimSpace(explicit); v != "" {
		return v
	}
	if strings.TrimSpace(serviceKey) == "" || grpcPort <= 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", DefaultHost(overrideHost, serviceKey), grpcPort)
}

func ResolveFromServicesMap(explicit string, services map[string]Service, overrideHost, serviceKey string) string {
	if v := strings.TrimSpace(explicit); v != "" {
		return v
	}
	svc, ok := services[serviceKey]
	if !ok || svc.Ports.Grpc <= 0 {
		return ""
	}
	return ResolveGRPCAddr("", overrideHost, serviceKey, svc.Ports.Grpc)
}

func DefaultHost(overrideHost, serviceKey string) string {
	if host := strings.TrimSpace(overrideHost); host != "" {
		return host
	}
	if host := strings.TrimSpace(os.Getenv("SERVICE_DISCOVERY_HOST")); host != "" {
		return host
	}

	switch strings.ToLower(strings.TrimSpace(os.Getenv("ENVIRONMENT"))) {
	case "production":
		return serviceKey
	case "development":
		return "localhost"
	}

	if isRunningInDocker() {
		return serviceKey
	}
	return "localhost"
}
