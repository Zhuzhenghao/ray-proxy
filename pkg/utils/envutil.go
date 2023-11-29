package utils

import (
	"os"
	"strings"
)

// AllowDuplicateKubeSystemID is whether to allow to add repeated k8s cluster.
var AllowDuplicateKubeSystemID = strings.ToLower(GetEnvWithDefault("ALLOW_DUPLICATE_KUBESYSTEMID", "")) == "true"

// DisableEgressProxy is whether to disable access member cluster by egress.
var DisableEgressProxy = strings.ToLower(GetEnvWithDefault("DISABLE_EGRESS_PROXY", "")) == "true"

func GetEnvWithDefault(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
