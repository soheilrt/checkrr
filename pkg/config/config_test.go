package config

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	yamlData := `
clients:
  - name: client1
    api_key: test_key
    host: http://example.com
    conditions:
      waiting_threshold: 10s
      download_timeout_threshold: 30s
      average_speed_threshold: 1.5
    options:
      keep_in_client: true
      blocklist: false
      skip_redownload: true
sleep_time: 5s
log_level: debug
`

	r := bytes.NewReader([]byte(yamlData))
	config, err := LoadConfig(r)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Clients) != 1 {
		t.Errorf("Expected 1 client, got %d", len(config.Clients))
	}

	client := config.Clients[0]
	if client.Name != "client1" {
		t.Errorf("Expected client name 'client1', got '%s'", client.Name)
	}

	if client.APIKey != "test_key" {
		t.Errorf("Expected APIKey 'test_key', got '%s'", client.APIKey)
	}

	if client.Conditions.WaitingThreshold != 10*time.Second {
		t.Errorf("Expected WaitingThreshold 10s, got %v", client.Conditions.WaitingThreshold)
	}

	if client.Options.KeepInClient != true {
		t.Errorf("Expected KeepInClient true, got %v", client.Options.KeepInClient)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	yamlData := `
clients:
  - name: client1
    api_key: ""
  - name: client2
`

	r := bytes.NewReader([]byte(yamlData))
	os.Setenv("API_KEY_CLIENT1", "override_key_1")
	defer os.Unsetenv("API_KEY_CLIENT1")
	os.Setenv("API_KEY_CLIENT2", "override_key_2")
	defer os.Unsetenv("API_KEY_CLIENT2")

	config, err := LoadConfig(r)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.Clients[0].APIKey != "override_key_1" {
		t.Errorf("Expected APIKey 'override_key_1', got '%s'", config.Clients[0].APIKey)
	}
	if config.Clients[1].APIKey != "override_key_2" {
		t.Errorf("Expected APIKey 'override_key_2', got '%s'", config.Clients[1].APIKey)
	}
}
