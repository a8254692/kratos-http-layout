package bootstrap

import (
	"fmt"
	"os"
)

// ServiceInfo ...
type ServiceInfo struct {
	Name     string
	Version  string
	Id       string
	Metadata map[string]string
}

// NewServiceInfo ...
func NewServiceInfo(name, version string) *ServiceInfo {
	id, _ := os.Hostname()
	return &ServiceInfo{
		Name:     name,
		Version:  version,
		Id:       id,
		Metadata: map[string]string{},
	}
}

// GetInstanceId ...
func (s *ServiceInfo) GetInstanceId() string {
	return fmt.Sprintf("%s.%s", s.Id, s.Name)
}

// SetMataData ...
func (s *ServiceInfo) SetMataData(k, v string) {
	s.Metadata[k] = v
}
