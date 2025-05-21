package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sayedamirkarim/ki-network-analyzer/internal/packet"
	"github.com/sayedamirkarim/ki-network-analyzer/pkg/models"
	"github.com/sayedamirkarim/ki-network-analyzer/pkg/version"
)

// APIResponse ist eine generische API-Antwortstruktur
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HealthResponse ist die Antwort für den Health-Check-Endpunkt
type HealthResponse struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	BuildDate  string            `json:"build_date,omitempty"`
	Commit     string            `json:"commit,omitempty"`
	Uptime     string            `json:"uptime"`
	Components map[string]string `json:"components"`
}

var (
	startTime = time.Now()
)

// HealthCheckHandler liefert Informationen über den Serverstatus
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:    "healthy",
		Version:   version.Version,
		BuildDate: version.BuildDate,
		Commit:    version.CommitHash,
		Uptime:    time.Since(startTime).String(),
		Components: map[string]string{
			"server":  "healthy",
			"storage": "healthy",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// AnalyzePcapHandler verarbeitet den Upload und die Analyse einer PCAP-Datei
func AnalyzePcapHandler(w http.ResponseWriter, r *http.Request, capturer *packet.PcapCapturer) {
	// Maximale Dateigröße festlegen (100 MB)
	maxFileSize := int64(100 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	// Multipart-Anfrage analysieren
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Datei zu groß oder ungültiges Format")
		return
	}

	// Datei aus der Anfrage extrahieren
	file, fileHeader, err := r.FormFile("pcap")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Keine PCAP-Datei in der Anfrage gefunden")
		return
	}
	defer file.Close()

	// Temporäre Datei erstellen
	tempDir := os.TempDir()
	tempFileName := filepath.Join(tempDir, "upload-"+time.Now().Format("20060102-150405")+".pcap")
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fehler beim Erstellen der temporären Datei")
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFileName)

	// Datei speichern
	_, err = io.Copy(tempFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fehler beim Speichern der Datei")
		return
	}
	tempFile.Close() // Schließen, um die Datei für das Lesen zu öffnen

	// PCAP-Datei öffnen
	err = capturer.OpenPcapFile(tempFileName)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Fehler beim Öffnen der PCAP-Datei: %v", err))
		return
	}

	// Kontext erstellen
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Paketerfassung starten
	packetChan, errChan := capturer.StartCapture(ctx)

	// Pakete sammeln
	var packets []*models.PacketInfo
	var gatewayPackets []*models.PacketInfo
	var packetCount, gatewayCount int

	// Maximale Analyse-Zeit festlegen (30 Sekunden)
	timeout := time.After(30 * time.Second)

packetLoop:
	for {
		select {
		case p, ok := <-packetChan:
			if !ok {
				break packetLoop
			}
			packetCount++

			// Gateway-Pakete separat sammeln
			if p.IsGatewayTraffic {
				gatewayCount++
				gatewayPackets = append(gatewayPackets, p)
			}

			// Begrenzte Anzahl von Paketen für die Antwort
			if len(packets) < 1000 {
				packets = append(packets, p)
			}

		case err, ok := <-errChan:
			if !ok {
				continue
			}
			log.Printf("Fehler bei der Paketverarbeitung: %v", err)

		case <-timeout:
			log.Println("Timeout bei der Paketanalyse")
			break packetLoop
		}
	}

	// Antwort erstellen
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("PCAP-Datei '%s' erfolgreich analysiert", fileHeader.Filename),
		Data: map[string]interface{}{
			"total_packets":      packetCount,
			"gateway_packets":    gatewayCount,
			"gateway_percentage": float64(gatewayCount) / float64(packetCount) * 100,
			"sample_packets":     packets[:min(len(packets), 100)], // Kleine Stichprobe zurückgeben
		},
	}

	// Als JSON antworten
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// WebSocketHandler verwaltet Websocket-Verbindungen für Live-Updates
func WebSocketHandler(w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader) {
	// Upgrade der HTTP-Verbindung zu einer Websocket-Verbindung
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Fehler beim Upgrade zu WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// WebSocket-Handler-Loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Fehler beim Lesen der WebSocket-Nachricht: %v", err)
			break
		}

		// Echo-Antwort (für Demos)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("Fehler beim Senden der WebSocket-Nachricht: %v", err)
			break
		}
	}
}

// GetGatewaysHandler gibt erkannte Gateways zurück
func GetGatewaysHandler(w http.ResponseWriter, r *http.Request) {
	// Mock-Daten für den MVP
	gateways := []map[string]interface{}{
		{
			"ip":        "192.168.1.1",
			"mac":       "11:22:33:44:55:66",
			"is_active": true,
			"role":      "default_gateway",
			"services":  []string{"DHCP", "DNS"},
		},
	}

	response := APIResponse{
		Success: true,
		Data:    gateways,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGatewayTrafficHandler gibt Gateway-Traffic-Statistiken zurück
func GetGatewayTrafficHandler(w http.ResponseWriter, r *http.Request) {
	// Mock-Daten für den MVP
	trafficStats := map[string]interface{}{
		"total_packets":      1000,
		"gateway_packets":    650,
		"gateway_percentage": 65.0,
		"protocols": map[string]int{
			"DNS":   150,
			"HTTP":  200,
			"DHCP":  50,
			"ARP":   100,
			"Other": 150,
		},
	}

	response := APIResponse{
		Success: true,
		Data:    trafficStats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGatewayEventsHandler gibt Gateway-relevante Ereignisse zurück
func GetGatewayEventsHandler(w http.ResponseWriter, r *http.Request) {
	// Mock-Daten für den MVP
	events := []models.GatewayEvent{
		{
			Timestamp:   time.Now().Add(-5 * time.Minute),
			EventType:   "dhcp",
			Description: "DHCP-Lease-Erneuerung",
			Severity:    "info",
			GatewayIP:   "192.168.1.1",
			ClientIP:    "192.168.1.100",
		},
		{
			Timestamp:   time.Now().Add(-10 * time.Minute),
			EventType:   "dns",
			Description: "DNS-Auflösung für example.com",
			Severity:    "info",
			GatewayIP:   "192.168.1.1",
			ClientIP:    "192.168.1.100",
		},
		{
			Timestamp:   time.Now().Add(-15 * time.Minute),
			EventType:   "arp",
			Description: "ARP-Anfrage für Gateway-MAC",
			Severity:    "info",
			GatewayIP:   "192.168.1.1",
			ClientIP:    "192.168.1.100",
		},
	}

	response := APIResponse{
		Success: true,
		Data:    events,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// respondWithError sendet eine Fehlerantwort an den Client
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// min gibt das Minimum von zwei Zahlen zurück
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
