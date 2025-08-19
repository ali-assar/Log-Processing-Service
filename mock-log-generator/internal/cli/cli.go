package cli

import (
	"log"

	"github.com/spf13/cobra"
)

type Config struct {
	Url        string
	IntervalMS int
}

func NewRoot() (*cobra.Command, *Config) {
	cfg := &Config{}

	cmd := &cobra.Command{
		Use:   "mock-log-generator",
		Short: "Start the mock WebSocket log generator",
	}

	// Flags
	cmd.Flags().StringVarP(&cfg.Url, "url", "u", ":8080", "HTTP listen address (e.g., :8080)")
	cmd.Flags().IntVarP(&cfg.IntervalMS, "interval-ms", "i", 0, "default interval in ms (0 = random per connection)")
	err := cmd.MarkFlagRequired("url")
	if err != nil {
		log.Fatal(err)
	}
	return cmd, cfg
}

func Parse() (*Config, error) {
	cmd, cfg := NewRoot()
	if err := cmd.Execute(); err != nil {
		return nil, err
	}
	return cfg, nil
}
