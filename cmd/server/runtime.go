package main

import (
	"log"
	"os"
	"strconv"
)

type RuntimeParams struct {
	servicePort string

	githubWebhookSecret string

	githubAppID             int64
	githubAppInstallationID int64
	githubAppPrivateKey     string
}

// checkEnvKeys will do the world's most basic validation and check for env var existence.
func checkEnvKeys() {
	missingEnvKeys := []string{}

	for _, requiredEnvKey := range requiredEnvKeys {
		if _, ok := os.LookupEnv(requiredEnvKey); !ok {
			missingEnvKeys = append(missingEnvKeys, requiredEnvKey)
		}
	}

	if len(missingEnvKeys) > 0 {
		log.Fatalf("missing defined env vars: %s", missingEnvKeys)
	}
}

func parsePortNumber() string {
	if envPort, ok := os.LookupEnv("PORT"); ok && len(envPort) > 0 {
		return envPort
	}

	return DefaultPort
}

func ParseRuntimeParameters() RuntimeParams {
	port := parsePortNumber()

	appID, _ := strconv.Atoi(os.Getenv("APP_ID"))
	installID, _ := strconv.Atoi(os.Getenv("INSTALLATION_ID"))

	return RuntimeParams{
		servicePort: port,

		githubWebhookSecret: os.Getenv("GH_SECRET"),

		githubAppID:             int64(appID),
		githubAppInstallationID: int64(installID),
		githubAppPrivateKey:     os.Getenv("PRIVATE_KEY"),
	}
}
