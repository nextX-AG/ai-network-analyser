package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config repräsentiert die Hauptkonfiguration der Anwendung
type Config struct {
	Server  ServerConfig  `json:"server"`
	Capture CaptureConfig `json:"capture"`
	Storage StorageConfig `json:"storage"`
	AI      AIConfig      `json:"ai"`
	Speech  SpeechConfig  `json:"speech"`
	Gateway GatewayConfig `json:"gateway"`
	Agent   *AgentConfig  `json:"agent,omitempty"`
}

// ServerConfig enthält die Konfiguration für den HTTP-Server
type ServerConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	EnableWebSocket bool   `json:"enable_websocket"`
	EnableCORS      bool   `json:"enable_cors"`
	StaticDir       string `json:"static_dir"`
}

// CaptureConfig enthält die Konfiguration für die Paketerfassung
type CaptureConfig struct {
	PCAPDir     string `json:"pcap_dir"`
	Interface   string `json:"interface"`
	PromiscMode bool   `json:"promisc_mode"`
	SnapLen     int    `json:"snap_len"`
	Filter      string `json:"filter"`
	BufferSize  int    `json:"buffer_size"`
	EnableLive  bool   `json:"enable_live"`
}

// StorageConfig enthält die Konfiguration für die Datenspeicherung
type StorageConfig struct {
	Type       string `json:"type"` // sqlite, memory
	Path       string `json:"path"` // Pfad zur SQLite-Datei
	AutoVacuum bool   `json:"auto_vacuum"`
	MaxPackets int    `json:"max_packets"` // Max. Anzahl zu speichernder Pakete
}

// AIConfig enthält die Konfiguration für KI-Integration
type AIConfig struct {
	Enabled     bool    `json:"enabled"`
	Provider    string  `json:"provider"` // openai, local
	APIKey      string  `json:"api_key,omitempty"`
	Endpoint    string  `json:"endpoint,omitempty"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// SpeechConfig enthält die Konfiguration für Speech2Text
type SpeechConfig struct {
	Enabled   bool   `json:"enabled"`
	Engine    string `json:"engine"` // whisper_local, whisper_api
	ModelPath string `json:"model_path,omitempty"`
	Language  string `json:"language"`
	APIKey    string `json:"api_key,omitempty"`
}

// GatewayConfig enthält gateway-spezifische Analysekonfigurationen
type GatewayConfig struct {
	DetectGateways       bool     `json:"detect_gateways"`
	KnownGateways        []string `json:"known_gateways"`
	TrackNAT             bool     `json:"track_nat"`
	TrackDNS             bool     `json:"track_dns"`
	TrackDHCP            bool     `json:"track_dhcp"`
	TrackARP             bool     `json:"track_arp"`
	DetectPortForwarding bool     `json:"detect_port_forwarding"`
	DetectDMZ            bool     `json:"detect_dmz"`
	DetectUPnP           bool     `json:"detect_upnp"`
	EnableAlerts         bool     `json:"enable_alerts"`
}

// AgentConfig enthält die Konfiguration für den Remote-Agent
type AgentConfig struct {
	// Auf welcher Adresse/Port der Agent lauscht
	Listen string `json:"listen"`

	// URL des Hauptservers für die Registrierung
	ServerURL string `json:"server_url"`

	// Zu verwendende Netzwerkschnittstelle für Packet-Capture
	Interface string `json:"interface"`

	// Name des Agents für die Identifikation
	Name string `json:"name"`

	// API-Schlüssel für die Authentifizierung mit dem Hauptserver
	APIKey string `json:"api_key,omitempty"`
}

// LoadConfig lädt die Konfiguration aus einer Datei
func LoadConfig(configPath string) (*Config, error) {
	// Standardkonfiguration
	config := DefaultConfig()

	// Konfigurationsdatei lesen, falls vorhanden
	if configPath != "" {
		file, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("Fehler beim Lesen der Konfigurationsdatei: %w", err)
		}

		if err := json.Unmarshal(file, config); err != nil {
			return nil, fmt.Errorf("Fehler beim Parsen der Konfigurationsdatei: %w", err)
		}
	}

	return config, nil
}

// DefaultConfig erstellt eine Standardkonfiguration
func DefaultConfig() *Config {
	execDir, _ := os.Executable()
	baseDir := filepath.Dir(execDir)

	return &Config{
		Server: ServerConfig{
			Host:            "127.0.0.1",
			Port:            8080,
			EnableWebSocket: true,
			EnableCORS:      true,
			StaticDir:       filepath.Join(baseDir, "web"),
		},
		Capture: CaptureConfig{
			PCAPDir:     filepath.Join(baseDir, "pcaps"),
			Interface:   "",
			PromiscMode: true,
			SnapLen:     65535,
			Filter:      "",
			BufferSize:  2 * 1024 * 1024, // 2MB
			EnableLive:  false,
		},
		Storage: StorageConfig{
			Type:       "sqlite",
			Path:       filepath.Join(baseDir, "data", "packets.db"),
			AutoVacuum: true,
			MaxPackets: 1000000,
		},
		AI: AIConfig{
			Enabled:     false,
			Provider:    "openai",
			Model:       "gpt-4",
			MaxTokens:   1000,
			Temperature: 0.1,
		},
		Speech: SpeechConfig{
			Enabled:   false,
			Engine:    "whisper_local",
			ModelPath: filepath.Join(baseDir, "models", "whisper.bin"),
			Language:  "auto",
		},
		Gateway: GatewayConfig{
			DetectGateways:       true,
			KnownGateways:        []string{},
			TrackNAT:             true,
			TrackDNS:             true,
			TrackDHCP:            true,
			TrackARP:             true,
			DetectPortForwarding: true,
			DetectDMZ:            true,
			DetectUPnP:           true,
			EnableAlerts:         true,
		},
	}
}

// SaveConfig speichert die Konfiguration in eine Datei
func SaveConfig(config *Config, configPath string) error {
	// Zielverzeichnis erstellen, falls es nicht existiert
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Fehler beim Erstellen des Konfigurationsverzeichnisses: %w", err)
	}

	// Konfiguration als JSON speichern
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("Fehler beim Umwandeln der Konfiguration in JSON: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("Fehler beim Speichern der Konfigurationsdatei: %w", err)
	}

	return nil
}
