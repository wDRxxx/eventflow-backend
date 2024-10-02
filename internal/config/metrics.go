package config

import (
	"net"
	"os"
)

type MetricsConfig struct {
	appName        string
	prometheusHost string
	prometheusPort string
}

func (c *MetricsConfig) AppName() string {
	return c.appName
}

func (c *MetricsConfig) PrometheusAddress() string {
	return net.JoinHostPort(c.prometheusHost, c.prometheusPort)
}

func NewMetricsConfig() *MetricsConfig {
	appName := os.Getenv("METRICS_APPNAME")
	if appName == "" {
		panic("METRICS_APPNAME environment variable is empty")
	}
	prometheusHost := os.Getenv("METRICS_PROMETHEUS_HOST")
	if prometheusHost == "" {
		panic("METRICS_PROMETHEUS_HOST environment variable is empty")
	}
	prometheusPort := os.Getenv("METRICS_PROMETHEUS_PORT")
	if prometheusPort == "" {
		panic("METRICS_PROMETHEUS_PORT environment variable is empty")
	}

	return &MetricsConfig{
		appName:        appName,
		prometheusHost: prometheusHost,
		prometheusPort: prometheusPort,
	}
}
