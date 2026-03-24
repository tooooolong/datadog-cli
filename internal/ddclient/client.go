package ddclient

import (
	"context"
	"fmt"
	"os"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
)

// FormatAPIError extracts a readable message from Datadog API errors.
func FormatAPIError(err error) string {
	if genErr, ok := err.(datadog.GenericOpenAPIError); ok {
		return fmt.Sprintf("%s — %s", err.Error(), string(genErr.Body()))
	}
	return err.Error()
}

// ResolveKeys returns API key and App key from env vars (priority) or config file (fallback).
func ResolveKeys() (apiKey, appKey string, source string) {
	apiKey = os.Getenv("DD_API_KEY")
	appKey = os.Getenv("DD_APP_KEY")
	if apiKey != "" && appKey != "" {
		return apiKey, appKey, "env"
	}

	cfg, err := LoadConfig()
	if err == nil {
		if apiKey == "" {
			apiKey = cfg.APIKey
		}
		if appKey == "" {
			appKey = cfg.AppKey
		}
		if apiKey != "" && appKey != "" {
			return apiKey, appKey, "config"
		}
	}

	return apiKey, appKey, ""
}

func NewClient(site string) (context.Context, *datadog.APIClient, error) {
	apiKey, appKey, _ := ResolveKeys()
	if apiKey == "" {
		return nil, nil, fmt.Errorf("DD_API_KEY not set. Run 'datadog login' or export DD_API_KEY")
	}
	if appKey == "" {
		return nil, nil, fmt.Errorf("DD_APP_KEY not set. Run 'datadog login' or export DD_APP_KEY")
	}

	ctx := context.WithValue(context.Background(), datadog.ContextAPIKeys, map[string]datadog.APIKey{
		"apiKeyAuth": {Key: apiKey},
		"appKeyAuth": {Key: appKey},
	})

	cfg := datadog.NewConfiguration()
	if site != "" {
		ctx = context.WithValue(ctx, datadog.ContextServerVariables, map[string]string{
			"site": site,
		})
	}

	return ctx, datadog.NewAPIClient(cfg), nil
}
