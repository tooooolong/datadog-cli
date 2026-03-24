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

func NewClient(site string) (context.Context, *datadog.APIClient, error) {
	apiKey := os.Getenv("DD_API_KEY")
	appKey := os.Getenv("DD_APP_KEY")
	if apiKey == "" {
		return nil, nil, fmt.Errorf("DD_API_KEY environment variable is not set")
	}
	if appKey == "" {
		return nil, nil, fmt.Errorf("DD_APP_KEY environment variable is not set")
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
