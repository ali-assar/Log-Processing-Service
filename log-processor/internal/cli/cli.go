package cli

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	Urls     []string
	HTTPAddr string
}

func NewRoot() (*cobra.Command, *string, *string) {
	var urls, httpAddr string

	cmd := &cobra.Command{
		Use:   "log-processor",
		Short: "Start the log processor",
	}

	// Example: --urls "ws://localhost:8080/ws/logs,ws://localhost:9090/ws/logs"
	cmd.Flags().StringVarP(&urls, "urls", "u", "", "Comma-separated list of WebSocket URLs (e.g., ws://localhost:8080/ws/logs)")
	cmd.Flags().StringVarP(&httpAddr, "http-addr", "a", ":9090", "HTTP listen address for the stats API (e.g., :9090)")

	err := cmd.MarkFlagRequired("urls")
	if err != nil {
		log.Fatal(err)
	}
	return cmd, &urls, &httpAddr
}

func UrlToSlice(urls string) []string {
	if urls == "" {
		return nil
	}
	// Split and trim spaces
	parts := strings.Split(urls, ",")
	out := make([]string, 0, len(parts))
	for _, v := range parts {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func Parse() (*Config, error) {
	cmd, urls, addr := NewRoot()
	if err := cmd.Execute(); err != nil {
		return nil, err
	}
	return &Config{
		Urls:     UrlToSlice(*urls),
		HTTPAddr: *addr,
	}, nil
}
