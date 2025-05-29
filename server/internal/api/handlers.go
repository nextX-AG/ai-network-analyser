package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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

// InterfaceInfo enthält Informationen über eine Netzwerkschnittstelle
type InterfaceInfo struct {
	Name        string   `json:"name"`
	Index       int      `json:"index"`
	MacAddress  string   `json:"mac_address,omitempty"`
	IPAddresses []string `json:"ip_addresses,omitempty"`
	IsUp        bool     `json:"is_up"`
	IsLoopback  bool     `json:"is_loopback"`
}

// LiveCaptureRequest enthält die Anfrageparameter für die Live-Capture
type LiveCaptureRequest struct {
	Interface string `json:"interface"`
	Filter    string `json:"filter,omitempty"`
}

var (
	startTime = time.Now()

	// Für die Verwaltung der aktiven Live-Capture
	activeCaptureContext context.Context
	activeCaptureCancel  context.CancelFunc
	activeCaptureStatus  string = "idle" // "idle", "running", "error"
	captureStatusMutex   sync.Mutex
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

// GetInterfacesHandler gibt eine Liste aller verfügbaren Netzwerkschnittstellen zurück
func GetInterfacesHandler(w http.ResponseWriter, r *http.Request) {
	// Alle Netzwerkschnittstellen abfragen
	interfaces, err := net.Interfaces()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Fehler beim Abrufen der Netzwerkschnittstellen: %v", err))
		return
	}

	// Ergebnis aufbereiten
	var interfaceInfos []InterfaceInfo

	for _, iface := range interfaces {
		// Nur aktive und nicht-virtuelle Schnittstellen berücksichtigen
		if iface.Flags&net.FlagUp == 0 {
			continue // Deaktivierte Schnittstellen überspringen
		}

		// IP-Adressen der Schnittstelle abrufen
		var ipAddresses []string
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				ipAddresses = append(ipAddresses, addr.String())
			}
		}

		// Schnittstellen-Info erstellen
		interfaceInfo := InterfaceInfo{
			Name:        iface.Name,
			Index:       iface.Index,
			MacAddress:  iface.HardwareAddr.String(),
			IPAddresses: ipAddresses,
			IsUp:        iface.Flags&net.FlagUp != 0,
			IsLoopback:  iface.Flags&net.FlagLoopback != 0,
		}

		interfaceInfos = append(interfaceInfos, interfaceInfo)
	}

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Data:    interfaceInfos,
	}

	w.Header().Set("Content-Type", "application/json")
	// Stellen wir sicher, dass keine Fehler bei der JSON-Kodierung auftreten
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Fehler bei der JSON-Kodierung: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}
}

// StartLiveCaptureHandler startet die Live-Capture auf einer Netzwerkschnittstelle
func StartLiveCaptureHandler(w http.ResponseWriter, r *http.Request, capturer *packet.PcapCapturer) {
	// Prüfen, ob bereits eine Capture läuft
	captureStatusMutex.Lock()
	if activeCaptureStatus == "running" {
		captureStatusMutex.Unlock()
		respondWithError(w, http.StatusConflict, "Eine Live-Capture läuft bereits")
		return
	}
	captureStatusMutex.Unlock()

	// Request-Body einlesen
	var request LiveCaptureRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Ungültiges Anfrageformat")
		return
	}

	// Netzwerkschnittstelle prüfen
	if request.Interface == "" {
		respondWithError(w, http.StatusBadRequest, "Keine Netzwerkschnittstelle angegeben")
		return
	}

	// Capture starten
	err = capturer.OpenLiveCapture(request.Interface)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Fehler beim Öffnen der Netzwerkschnittstelle: %v", err))
		return
	}

	// Neuen Kontext für die Capture erstellen
	activeCaptureContext, activeCaptureCancel = context.WithCancel(context.Background())

	// Capture starten
	packetChan, errChan := capturer.StartCapture(activeCaptureContext)

	// Status aktualisieren
	captureStatusMutex.Lock()
	activeCaptureStatus = "running"
	captureStatusMutex.Unlock()

	// Verarbeitung in Goroutine starten
	go func() {
		var packetCount int
		var gatewayCount int

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case p, ok := <-packetChan:
				if !ok {
					log.Println("Paketkanal geschlossen")
					captureStatusMutex.Lock()
					activeCaptureStatus = "idle"
					captureStatusMutex.Unlock()
					return
				}

				packetCount++
				if p.IsGatewayTraffic {
					gatewayCount++
				}

				// Hier könnte die Verarbeitung erfolgen und Daten an WebSockets gesendet werden
				// Diese Logik ist bereits in processLivePackets im Hauptprogramm implementiert

			case err, ok := <-errChan:
				if !ok {
					continue
				}
				log.Printf("Fehler bei der Live-Capture: %v", err)

			case <-ticker.C:
				log.Printf("Live-Capture Status: %d Pakete erfasst, davon %d Gateway-Pakete",
					packetCount, gatewayCount)

			case <-activeCaptureContext.Done():
				log.Println("Live-Capture wurde gestoppt")
				captureStatusMutex.Lock()
				activeCaptureStatus = "idle"
				captureStatusMutex.Unlock()
				return
			}
		}
	}()

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("Live-Capture auf Schnittstelle %s gestartet", request.Interface),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StopLiveCaptureHandler stoppt die aktive Live-Capture
func StopLiveCaptureHandler(w http.ResponseWriter, r *http.Request) {
	// Prüfen, ob eine Capture läuft
	captureStatusMutex.Lock()
	if activeCaptureStatus != "running" {
		captureStatusMutex.Unlock()
		respondWithError(w, http.StatusBadRequest, "Keine aktive Live-Capture vorhanden")
		return
	}
	captureStatusMutex.Unlock()

	// Capture stoppen
	if activeCaptureCancel != nil {
		activeCaptureCancel()
	}

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Message: "Live-Capture wurde gestoppt",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
