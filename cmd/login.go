package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Configure Datadog API credentials",
	Long: `Interactively set your Datadog API key and Application key.

Keys are saved to ~/.local/config/datadog-cli/config.json.
Environment variables DD_API_KEY and DD_APP_KEY take priority over saved config.`,
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	existing, _ := ddclient.LoadConfig()
	isTTY := term.IsTerminal(int(os.Stdin.Fd()))
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your Datadog credentials.")
	fmt.Printf("Keys will be saved to: %s\n\n", ddclient.ConfigPath())

	apiKey, err := readSecret(reader, isTTY, "DD API Key", existingHint(existing, "api"))
	if err != nil {
		return err
	}
	if apiKey == "" {
		return fmt.Errorf("API Key cannot be empty")
	}

	appKey, err := readSecret(reader, isTTY, "DD App Key", existingHint(existing, "app"))
	if err != nil {
		return err
	}
	if appKey == "" {
		return fmt.Errorf("App Key cannot be empty")
	}

	fmt.Println()
	fmt.Printf("  API Key: %s\n", maskKey(apiKey))
	fmt.Printf("  App Key: %s\n", maskKey(appKey))
	fmt.Println()

	if isTTY {
		fmt.Print("Press Enter to save, or Ctrl+C to cancel...")
		reader.ReadBytes('\n')
	}

	cfg := &ddclient.Config{
		APIKey: apiKey,
		AppKey: appKey,
	}
	if err := ddclient.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Credentials saved to %s\n", ddclient.ConfigPath())
	return nil
}

func readSecret(reader *bufio.Reader, isTTY bool, prompt, hint string) (string, error) {
	if hint != "" {
		fmt.Printf("%s [current: %s]: ", prompt, hint)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	if isTTY {
		raw, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		value := strings.TrimSpace(string(raw))
		if value != "" {
			fmt.Printf("%s\n", maskKey(value))
		} else {
			fmt.Println()
		}
		return value, nil
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func existingHint(cfg *ddclient.Config, which string) string {
	if cfg == nil {
		return ""
	}
	switch which {
	case "api":
		if cfg.APIKey != "" {
			return maskKey(cfg.APIKey)
		}
	case "app":
		if cfg.AppKey != "" {
			return maskKey(cfg.AppKey)
		}
	}
	return ""
}
