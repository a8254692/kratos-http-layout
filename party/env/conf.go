package env

import "github.com/spf13/viper"

const (
	PathMode           = "mode"
	PathServiceName    = "name"
	PathServiceVersion = "version"

	PathMonitorHttpHost = "server.monitorHttp.host"
	PathMonitorHttpPort = "server.monitorHttp.port"
)

// GetMode ...
func GetMode() string {
	return viper.GetString(PathMode)
}

// GetServiceName ...
func GetServiceName() string {
	return viper.GetString(PathServiceName)
}

// GetServiceVersion ...
func GetServiceVersion() string {
	return viper.GetString(PathServiceVersion)
}
