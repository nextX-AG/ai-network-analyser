package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/sayedamirkarim/ki-network-analyzer/internal/config"
	"github.com/sayedamirkarim/ki-network-analyzer/internal/packet"
	"github.com/sayedamirkarim/ki-network-analyzer/pkg/models"
)

// AgentStatus enthält die aktuellen Status-Informationen des Agents
type AgentStatus struct {
	Name            string    `json:"name"`
	Status          string    `json:"status"` // "idle", "capturing", "error"
	LastHeartbeat   time.Time `json:"last_heartbeat"`
	StartTime       time.Time `json:"start_time"`
	PacketsCaptured int       `json:"packets_captured"`
	Interface       string    `json:"interface"`
	Error           string    `json:"error,omitempty"`
}

// AgentInfo enthält die Registrierungsinformationen für den Server
type AgentInfo struct {
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	Interfaces []string `json:"interfaces"`
	Version    string   `json:"version"`
	OS         string   `json:"os"`
	Hostname   string   `json:"hostname"`
}

// APIResponse ist eine generische API-Antwortstruktur
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CaptureRequest enthält die Konfiguration für eine Capture-Anfrage
type CaptureRequest struct {
	Interface string `json:"interface"`
	Filter    string `json:"filter,omitempty"`
}

// CaptureAgent verwaltet die Packet-Capture und API-Kommunikation
type CaptureAgent struct {
	config       *config.Config
	status       AgentStatus
	statusMutex  sync.RWMutex
	capturer     *packet.PcapCapturer
	activeCtx    context.Context
	cancelFunc   context.CancelFunc
	clients      map[*websocket.Conn]bool
	clientsMutex sync.Mutex
}

// NewCaptureAgent erstellt eine neue Instanz des CaptureAgent
func NewCaptureAgent(config *config.Config) *CaptureAgent {
	return &CaptureAgent{
		config: config,
		status: AgentStatus{
			Name:          config.Agent.Name,
			Status:        "idle",
			StartTime:     time.Now(),
			LastHeartbeat: time.Now(),
			Interface:     config.Agent.Interface,
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// Init initialisiert den CaptureAgent
func (a *CaptureAgent) Init() error {
	a.capturer = packet.NewPcapCapturer(a.config)

	// Heartbeat-Routine starten
	go a.heartbeatRoutine()

	// Automatische Registrierung versuchen, wenn eine Server-URL konfiguriert ist
	if a.config.Agent.ServerURL != "" {
		log.Printf("Versuche automatische Registrierung beim Server: %s", a.config.Agent.ServerURL)
		go func() {
			// Kurz warten, um sicherzustellen, dass der Server gestartet ist
			time.Sleep(2 * time.Second)

			// Registrierung versuchen
			if err := a.Register(); err != nil {
				log.Printf("Automatische Registrierung fehlgeschlagen: %v", err)
				log.Println("Der Agent wird im Offline-Modus ausgeführt. Verwenden Sie die Web-UI zur manuellen Registrierung.")

				// Status aktualisieren
				a.statusMutex.Lock()
				a.status.Status = "error"
				a.status.Error = fmt.Sprintf("Registrierung fehlgeschlagen: %v", err)
				a.statusMutex.Unlock()
			} else {
				log.Println("Automatische Registrierung beim Server erfolgreich")

				// Status aktualisieren
				a.statusMutex.Lock()
				a.status.Status = "idle"
				a.status.Error = ""
				a.statusMutex.Unlock()
			}
		}()
	} else {
		log.Println("Keine Server-URL konfiguriert. Verwenden Sie die Web-UI zur manuellen Konfiguration und Registrierung.")
	}

	return nil
}

// Register registriert den Agent beim Hauptserver
func (a *CaptureAgent) Register() error {
	// Netzwerkschnittstellen abfragen
	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %v", err)
	}

	var interfaceNames []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			interfaceNames = append(interfaceNames, iface.Name)
		}
	}

	// Hostname für die Registrierung abrufen
	hostname, _ := a.config.Agent.Name, ""

	// AgentInfo erstellen
	info := AgentInfo{
		Name:       a.config.Agent.Name,
		URL:        fmt.Sprintf("http://%s", a.config.Agent.Listen),
		Interfaces: interfaceNames,
		Version:    "0.1.0", // TODO: aus Versionsdatei lesen
		OS:         "linux", // TODO: dynamisch ermitteln
		Hostname:   hostname,
	}

	// JSON-Kodierung
	jsonData, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal agent info: %v", err)
	}

	// Registrierungs-URL zusammensetzen
	url := fmt.Sprintf("%s/api/agents/register", a.config.Agent.ServerURL)

	// HTTP-Request senden
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if a.config.Agent.APIKey != "" {
		req.Header.Set("X-API-Key", a.config.Agent.APIKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send registration request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	log.Printf("Agent registered successfully with server %s", a.config.Agent.ServerURL)
	return nil
}

// Unregister meldet den Agent vom Hauptserver ab
func (a *CaptureAgent) Unregister() error {
	// TODO: Implementieren Sie die Abmeldung
	return nil
}

// RegisterRoutes registriert die HTTP-Endpunkte für den Agent
func (a *CaptureAgent) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", a.healthHandler).Methods("GET")
	router.HandleFunc("/status", a.statusHandler).Methods("GET")
	router.HandleFunc("/capture/start", a.startCaptureHandler).Methods("POST")
	router.HandleFunc("/capture/stop", a.stopCaptureHandler).Methods("POST")
	router.HandleFunc("/ws", a.websocketHandler)

	// Weitere Routen hier registrieren...
}

// healthHandler gibt den Gesundheitszustand des Agents zurück
func (a *CaptureAgent) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":     "healthy",
			"uptime":     time.Since(a.status.StartTime).String(),
			"agent_name": a.config.Agent.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// statusHandler gibt den aktuellen Status des Agents zurück
func (a *CaptureAgent) statusHandler(w http.ResponseWriter, r *http.Request) {
	a.statusMutex.RLock()
	status := a.status
	a.statusMutex.RUnlock()

	response := APIResponse{
		Success: true,
		Data:    status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// startCaptureHandler startet die Paketerfassung
func (a *CaptureAgent) startCaptureHandler(w http.ResponseWriter, r *http.Request) {
	a.statusMutex.Lock()
	if a.status.Status == "capturing" {
		a.statusMutex.Unlock()
		respondWithError(w, http.StatusConflict, "Capture already in progress")
		return
	}
	a.statusMutex.Unlock()

	// Anfrage parsen
	var request CaptureRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Interface überprüfen und ggf. aus Konfiguration verwenden
	captureInterface := request.Interface
	if captureInterface == "" {
		captureInterface = a.config.Agent.Interface
		if captureInterface == "" {
			respondWithError(w, http.StatusBadRequest, "No interface specified")
			return
		}
	}

	// Capture öffnen
	if err := a.capturer.OpenLiveCapture(captureInterface); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Failed to open interface %s: %v", captureInterface, err))
		return
	}

	// Context für die Capture erstellen
	ctx, cancel := context.WithCancel(context.Background())
	a.activeCtx = ctx
	a.cancelFunc = cancel

	// Capture starten
	packetChan, errChan := a.capturer.StartCapture(ctx)

	// Status aktualisieren
	a.statusMutex.Lock()
	a.status.Status = "capturing"
	a.status.Interface = captureInterface
	a.status.PacketsCaptured = 0
	a.status.Error = ""
	a.statusMutex.Unlock()

	// Paketverarbeitung in Goroutine starten
	go a.processPackets(packetChan, errChan)

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("Capture started on interface %s", captureInterface),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// stopCaptureHandler stoppt die Paketerfassung
func (a *CaptureAgent) stopCaptureHandler(w http.ResponseWriter, r *http.Request) {
	a.statusMutex.Lock()
	if a.status.Status != "capturing" {
		a.statusMutex.Unlock()
		respondWithError(w, http.StatusBadRequest, "No active capture to stop")
		return
	}
	a.statusMutex.Unlock()

	// Capture stoppen
	if a.cancelFunc != nil {
		a.cancelFunc()
	}

	// Status aktualisieren
	a.statusMutex.Lock()
	a.status.Status = "idle"
	a.statusMutex.Unlock()

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Message: "Capture stopped",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// websocketHandler verwaltet WebSocket-Verbindungen für Paket-Streaming
func (a *CaptureAgent) websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // In Produktion sollte dies eingeschränkt werden
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Client registrieren
	a.clientsMutex.Lock()
	a.clients[conn] = true
	a.clientsMutex.Unlock()

	// Handler für eingehende Nachrichten
	go func() {
		defer func() {
			conn.Close()
			a.clientsMutex.Lock()
			delete(a.clients, conn)
			a.clientsMutex.Unlock()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}
			// Hier könnten Befehle vom Client verarbeitet werden
		}
	}()
}

// heartbeatRoutine sendet regelmäßig Heartbeats an den Hauptserver
func (a *CaptureAgent) heartbeatRoutine() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		a.statusMutex.Lock()
		a.status.LastHeartbeat = time.Now()
		a.statusMutex.Unlock()

		// TODO: Heartbeat an den Hauptserver senden
		log.Printf("Heartbeat: Agent %s is alive", a.config.Agent.Name)
	}
}

// processPackets verarbeitet eingehende Pakete
func (a *CaptureAgent) processPackets(packetChan <-chan *models.PacketInfo, errChan <-chan error) {
	for {
		select {
		case packet, ok := <-packetChan:
			if !ok {
				log.Println("Packet channel closed")
				return
			}

			// Paket zählen
			a.statusMutex.Lock()
			a.status.PacketsCaptured++
			a.statusMutex.Unlock()

			// Paket an alle verbundenen Clients senden
			a.broadcastPacket(packet)

		case err, ok := <-errChan:
			if !ok {
				continue
			}
			log.Printf("Error during packet capture: %v", err)

			a.statusMutex.Lock()
			a.status.Error = err.Error()
			a.statusMutex.Unlock()

		case <-a.activeCtx.Done():
			log.Println("Capture context cancelled")
			return
		}
	}
}

// broadcastPacket sendet ein Paket an alle verbundenen WebSocket-Clients
func (a *CaptureAgent) broadcastPacket(packet *models.PacketInfo) {
	// Vereinfachte Paketdarstellung für die Übertragung
	packetData := map[string]interface{}{
		"type": "packet",
		"data": map[string]interface{}{
			"timestamp":  packet.Timestamp,
			"source_ip":  packet.SourceIP.String(),
			"dest_ip":    packet.DestinationIP.String(),
			"protocol":   packet.Protocol,
			"length":     packet.Length,
			"is_gateway": packet.IsGatewayTraffic,
			"summary":    fmt.Sprintf("%s: %s -> %s", packet.Protocol, packet.SourceIP, packet.DestinationIP),
		},
	}

	data, err := json.Marshal(packetData)
	if err != nil {
		log.Printf("Error marshaling packet data: %v", err)
		return
	}

	a.clientsMutex.Lock()
	defer a.clientsMutex.Unlock()

	for client := range a.clients {
		if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error sending packet to client: %v", err)
			client.Close()
			delete(a.clients, client)
		}
	}
}

// respondWithError sendet eine Fehlerantwort
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
