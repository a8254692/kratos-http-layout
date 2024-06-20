package env

import "gitlab.top.slotssprite.com/my/api-layout/party/util"

const (
	ModeLocal      = "local"
	ModeDevelop    = "dev"
	ModeTest       = "test"
	ModeProduction = "prod"
)

// IsDevelopment ...
func IsDevelopment() bool {
	return util.IsContainsString(GetMode(), []string{ModeLocal, ModeDevelop, ModeTest})
}
