package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
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
	ActiveFilter    string    `json:"active_filter,omitempty"`
}

// AgentInfo enthält die Registrierungsinformationen für den Server
type AgentInfo struct {
	Name             string                   `json:"name"`
	URL              string                   `json:"url"`
	Interfaces       []string                 `json:"interfaces"`
	InterfaceDetails []map[string]interface{} `json:"interface_details"`
	Version          string                   `json:"version"`
	OS               string                   `json:"os"`
	Hostname         string                   `json:"hostname"`
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

	// Sicherstellen, dass Interface im Status gesetzt ist
	a.statusMutex.Lock()
	a.status.Interface = a.config.Agent.Interface
	a.statusMutex.Unlock()

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
	var interfaceDetails []map[string]interface{}

	// Hostname für die Registrierung abrufen
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = a.config.Agent.Name
	}

	// Die tatsächliche IP-Adresse ermitteln, unter der der Agent erreichbar ist
	host, port, err := parseListenAddress(a.config.Agent.Listen)
	if err != nil {
		return fmt.Errorf("failed to parse listen address: %v", err)
	}

	// Wenn der Host 0.0.0.0 ist, müssen wir die tatsächliche IP-Adresse ermitteln
	actualIP := host
	if host == "0.0.0.0" || host == "::" || host == "" {
		// Die Server-URL parsen, um die Netzwerk-Route zu bestimmen
		serverURL, err := url.Parse(a.config.Agent.ServerURL)
		if err != nil {
			log.Printf("Warnung: Konnte Server-URL nicht parsen: %v", err)
		} else {
			// Die tatsächliche IP-Adresse ermitteln, die zum Server routet
			actualIP, err = getOutboundIP(serverURL.Hostname())
			if err != nil {
				log.Printf("Warnung: Konnte ausgehende IP nicht ermitteln: %v", err)
				// Fallback auf eine nicht-lokale IP
				actualIP = getFirstNonLocalhostIP()
			}
		}
	}

	// Die vollständige URL des Agents erstellen
	agentURL := fmt.Sprintf("http://%s:%s", actualIP, port)

	// Detaillierte Schnittstelleninformationen sammeln
	for _, iface := range ifaces {
		// Lokale Loopback-Schnittstellen und inaktive Schnittstellen überspringen
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// IP-Adressen der Schnittstelle abrufen
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("Error getting addresses for interface %s: %v", iface.Name, err)
			continue
		}

		// IP-Adressen sammeln
		var ipStrings []string
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil || ipNet.IP.To16() != nil {
					ipStrings = append(ipStrings, ipNet.IP.String())
				}
			}
		}

		// Bridge-Status prüfen (Linux-spezifisch)
		isBridge := false
		bridgePorts := ""

		// Bridge-Schnittstellen erkennen (Linux-spezifisch)
		if _, err := os.Stat(fmt.Sprintf("/sys/class/net/%s/bridge", iface.Name)); err == nil {
			isBridge = true

			// Bridge-Ports auslesen
			files, err := os.ReadDir(fmt.Sprintf("/sys/class/net/%s/brif", iface.Name))
			if err == nil {
				var ports []string
				for _, file := range files {
					ports = append(ports, file.Name())
				}
				bridgePorts = strings.Join(ports, ", ")
			}
		}

		// Schnittstelle zur Liste hinzufügen
		interfaceNames = append(interfaceNames, iface.Name)

		// Detaillierte Informationen hinzufügen
		ifaceDetails := map[string]interface{}{
			"name":         iface.Name,
			"mac":          iface.HardwareAddr.String(),
			"ips":          ipStrings,
			"is_bridge":    isBridge,
			"bridge_ports": bridgePorts,
			"flags":        iface.Flags.String(),
			"mtu":          iface.MTU,
		}
		interfaceDetails = append(interfaceDetails, ifaceDetails)
	}

	// AgentInfo erstellen
	info := AgentInfo{
		Name:             a.config.Agent.Name,
		URL:              agentURL,
		Interfaces:       interfaceNames,
		InterfaceDetails: interfaceDetails,
		Version:          "0.1.0", // TODO: aus Versionsdatei lesen
		OS:               runtime.GOOS,
		Hostname:         hostname,
	}

	// JSON-Kodierung
	jsonData, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal agent info: %v", err)
	}

	// Überprüfe die Server-URL
	if a.config.Agent.ServerURL == "" {
		return fmt.Errorf("server URL is not configured")
	}

	// Registrierungs-URL zusammensetzen
	registerURL := fmt.Sprintf("%s/api/agents/register", a.config.Agent.ServerURL)
	log.Printf("Sending registration request to: %s", registerURL)

	// HTTP-Request senden
	req, err := http.NewRequest("POST", registerURL, bytes.NewBuffer(jsonData))
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
		// Status auf error setzen
		a.statusMutex.Lock()
		a.status.Status = "error"
		a.status.Error = fmt.Sprintf("Verbindung zum Server fehlgeschlagen: %v", err)
		a.statusMutex.Unlock()
		return fmt.Errorf("failed to send registration request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Status auf error setzen
		a.statusMutex.Lock()
		a.status.Status = "error"
		a.status.Error = fmt.Sprintf("Server antwortete mit Status: %d", resp.StatusCode)
		a.statusMutex.Unlock()
		return fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}

	// Status auf idle setzen bei erfolgreicher Registrierung
	a.statusMutex.Lock()
	a.status.Status = "idle"
	a.status.Error = ""
	a.statusMutex.Unlock()

	log.Printf("Agent registered successfully with server %s", a.config.Agent.ServerURL)
	return nil
}

// Hilfsfunktionen für die IP-Adressermittlung

// parseListenAddress zerlegt eine Adresse im Format "host:port" oder ":port"
func parseListenAddress(addr string) (host, port string, err error) {
	// Standardwert für Port setzen, falls nicht angegeben
	if !strings.Contains(addr, ":") {
		return addr, "8090", nil
	}

	host, port, err = net.SplitHostPort(addr)
	if err != nil {
		return "", "", err
	}

	// Wenn kein Host angegeben ist (z.B. ":8090"), leeren String zurückgeben
	if host == "" {
		host = "0.0.0.0"
	}

	return host, port, nil
}

// getOutboundIP ermittelt die lokale IP-Adresse, die für eine Verbindung nach außen verwendet wird
func getOutboundIP(destination string) (string, error) {
	// Wenn das Ziel keine gültige Adresse ist, verwenden wir einen öffentlichen DNS-Server
	if net.ParseIP(destination) == nil && !strings.Contains(destination, ".") {
		destination = "8.8.8.8:80" // Google DNS
	} else if !strings.Contains(destination, ":") {
		// Port hinzufügen, wenn keiner angegeben ist
		destination = destination + ":80"
	}

	// UDP-Verbindung herstellen (wird nicht wirklich aufgebaut)
	conn, err := net.Dial("udp", destination)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Die lokale Adresse der Verbindung ermitteln
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// getFirstNonLocalhostIP ermittelt die erste nicht-lokale IP-Adresse
func getFirstNonLocalhostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

// Unregister meldet den Agent vom Hauptserver ab
func (a *CaptureAgent) Unregister() error {
	// TODO: Implementieren Sie die Abmeldung
	return nil
}

// RegisterRoutes registriert die HTTP-Endpunkte für den Agent
func (a *CaptureAgent) RegisterRoutes(router *mux.Router) {
	// CORS-Middleware für alle Routen hinzufügen
	router.Use(a.corsMiddleware)

	router.HandleFunc("/health", a.healthHandler).Methods("GET")
	router.HandleFunc("/status", a.statusHandler).Methods("GET")
	router.HandleFunc("/capture/start", a.startCaptureHandler).Methods("POST")
	router.HandleFunc("/capture/stop", a.stopCaptureHandler).Methods("POST")
	router.HandleFunc("/capture/set-interface", a.setInterfaceHandler).Methods("POST")
	router.HandleFunc("/ws", a.websocketHandler)

	// Weitere Routen hier registrieren...
}

// corsMiddleware ist ein Middleware für CORS-Unterstützung (Cross-Origin Resource Sharing)
func (a *CaptureAgent) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS-Header für alle Anfragen setzen
		w.Header().Set("Access-Control-Allow-Origin", "*") // Im Produktivbetrieb einschränken
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		// Preflight-Anfragen mit OPTIONS direkt beantworten
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Hauptanfrage verarbeiten
		next.ServeHTTP(w, r)
	})
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

	// BPF-Filter setzen, wenn angegeben
	if request.Filter != "" {
		// Filter in Konfiguration setzen
		a.config.Capture.Filter = request.Filter
		log.Printf("BPF-Filter für Capture gesetzt: %s", request.Filter)
	} else {
		// Filter zurücksetzen, wenn keiner angegeben wurde
		a.config.Capture.Filter = ""
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
	a.status.ActiveFilter = request.Filter
	a.statusMutex.Unlock()

	// Paketverarbeitung in Goroutine starten
	go a.processPackets(packetChan, errChan)

	// Erfolgreiche Antwort senden
	responseMessage := fmt.Sprintf("Capture started on interface %s", captureInterface)
	if request.Filter != "" {
		responseMessage += fmt.Sprintf(" with filter: %s", request.Filter)
	}

	response := APIResponse{
		Success: true,
		Message: responseMessage,
		Data: map[string]interface{}{
			"interface": captureInterface,
			"filter":    request.Filter,
		},
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
	a.status.ActiveFilter = ""
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
		currentStatus := a.status.Status
		packetsCaptured := a.status.PacketsCaptured
		a.statusMutex.Unlock()

		// Heartbeat an den Hauptserver senden, wenn eine Server-URL konfiguriert ist
		if a.config.Agent.ServerURL != "" {
			// Aktiver Filter für Heartbeat
			a.statusMutex.RLock()
			activeFilter := a.status.ActiveFilter
			a.statusMutex.RUnlock()

			// Heartbeat-Daten vorbereiten
			heartbeatData := map[string]interface{}{
				"name":             a.config.Agent.Name,
				"status":           currentStatus,
				"packets_captured": packetsCaptured,
				"interface":        a.config.Agent.Interface,
				"active_filter":    activeFilter,
			}

			jsonData, err := json.Marshal(heartbeatData)
			if err != nil {
				log.Printf("Fehler beim Erstellen des Heartbeats: %v", err)
				continue
			}

			// Heartbeat-URL zusammensetzen
			heartbeatURL := fmt.Sprintf("%s/api/agents/heartbeat", a.config.Agent.ServerURL)

			// HTTP-Request senden
			req, err := http.NewRequest("POST", heartbeatURL, bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("Fehler beim Erstellen des Heartbeat-Requests: %v", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")
			if a.config.Agent.APIKey != "" {
				req.Header.Set("X-API-Key", a.config.Agent.APIKey)
			}

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Fehler beim Senden des Heartbeats: %v", err)
				continue
			}

			// Antwort verwerfen
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("Heartbeat wurde vom Server nicht akzeptiert. Status: %d", resp.StatusCode)
			} else {
				log.Printf("Heartbeat: Agent %s ist aktiv und mit dem Server verbunden (Pakete: %d)", a.config.Agent.Name, packetsCaptured)
			}
		} else {
			// Kein Server konfiguriert, lokale Protokollierung
			log.Printf("Heartbeat: Agent %s is alive (keine Server-Verbindung, Pakete: %d)", a.config.Agent.Name, packetsCaptured)
		}
	}
}

// processPackets verarbeitet eingehende Pakete
func (a *CaptureAgent) processPackets(packetChan <-chan *models.PacketInfo, errChan <-chan error) {
	// Debug-Ausgabe beim Start
	log.Println("DEBUG: Paketverarbeitung gestartet")

	// Zähler für das Debugging
	debugCounter := 0
	lastLogTime := time.Now()

	for {
		select {
		case packet, ok := <-packetChan:
			if !ok {
				log.Println("Packet channel closed")
				return
			}

			// Debug-Ausgabe alle 10 Pakete oder alle 5 Sekunden
			debugCounter++
			if debugCounter%10 == 0 || time.Since(lastLogTime) > 5*time.Second {
				log.Printf("DEBUG: Paket %d empfangen: %s -> %s (Protokoll: %s, Länge: %d)",
					debugCounter, packet.SourceIP, packet.DestinationIP, packet.Protocol, packet.Length)
				lastLogTime = time.Now()
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

// setInterfaceHandler setzt die aktive Schnittstelle für die Datenerfassung
func (a *CaptureAgent) setInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	// Anfrage-Body parsen
	var req struct {
		Interface string `json:"interface"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Fehler beim Parsen der Anfrage: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Prüfen, ob die gewünschte Schnittstelle existiert
	validInterface := false
	ifaces, err := net.Interfaces()
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Fehler beim Auflisten der Netzwerkschnittstellen: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	for _, iface := range ifaces {
		if iface.Name == req.Interface {
			validInterface = true
			break
		}
	}

	if !validInterface {
		response := APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Netzwerkschnittstelle '%s' nicht gefunden", req.Interface),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Die Schnittstelle in der Konfiguration aktualisieren
	a.config.Agent.Interface = req.Interface

	// Auch den Status aktualisieren
	a.statusMutex.Lock()
	a.status.Interface = req.Interface
	a.statusMutex.Unlock()

	// Die Konfiguration speichern
	if err := a.saveConfig(); err != nil {
		log.Printf("Warnung: Fehler beim Speichern der Konfiguration nach Interface-Änderung: %v", err)
	}

	// Auch dem Capturer mitteilen, dass sich die Schnittstelle geändert hat
	a.capturer.UpdateInterface(req.Interface)

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("Netzwerkschnittstelle erfolgreich auf '%s' gesetzt", req.Interface),
		Data: map[string]string{
			"interface": req.Interface,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
