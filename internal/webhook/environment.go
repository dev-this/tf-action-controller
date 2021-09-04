package webhook

import "os"

var (
	GhOwner = getEnv("GH_OWNER", "")
	tfPath  = getEnv("TF_BIN", "/usr/bin/terraform")
)

func getEnv(key, fallback string) *string {
	if value, ok := os.LookupEnv(key); ok {
		return &value
	}

	if fallback != "" {
		return &fallback
	}

	return nil
}
