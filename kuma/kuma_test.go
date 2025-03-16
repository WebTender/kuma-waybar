package kuma

import (
	"testing"
)

func TestParseKumaMetrics(t *testing.T) {
	t.Run("De-duplication case", func(t *testing.T) {
		metrics, monitors := parseBody(`
# Comments allowed
monitor_cert_days_remaining{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 66
monitor_cert_is_valid{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 1
monitor_response_time{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 2.07
monitor_response_time{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 2.07
monitor_status{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 1
monitor_status{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 0
`)
		if len(metrics) != 6 {
			// Maybe they should by the key prefix
			t.Errorf("metrics are not de-duplicated %v", len(metrics))
		}

		if len(monitors) != 1 {
			t.Errorf("Error: %v", monitors)
		}

		if monitors[0].Name != "My Website" {
			t.Errorf("Error: %v", monitors[0].Name)
		}

		if monitors[0].Type != "http" {
			t.Errorf("Error: %v", monitors[0].Type)
		}

		if monitors[0].Url != "https://mywebsite.com" {
			t.Errorf("Error: %v", monitors[0].Url)
		}

		if monitors[0].Status != 0 {
			t.Errorf("Error: %v should be the last status", monitors[0].Status)
		}
	})

	t.Run("Missing Status", func(t *testing.T) {
		metrics, monitors := parseBody(`
# Comments allowed
monitor_cert_days_remaining{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 66
monitor_cert_is_valid{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 1
monitor_response_time{monitor_name="My Website",monitor_type="http",monitor_url="https://mywebsite.com"} 2.07
`)

		if len(metrics) != 3 {
			t.Errorf("metrics are not de-duplicated %v", len(metrics))
		}

		if len(monitors) != 1 {
			t.Fatalf("Error: %v", monitors)
		}

		if monitors[0].Status != Paused {
			t.Errorf("Error: Assumes status should default to paused %v", monitors[0].Status)
		}
	})
}
