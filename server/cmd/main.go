package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/sayedamirkarim/ki-network-analyzer/internal/api"
	"github.com/sayedamirkarim/ki-network-analyzer/internal/config"
	"github.com/sayedamirkarim/ki-network-analyzer/internal/packet"
	"github.com/sayedamirkarim/ki-network-analyzer/pkg/models"
)

var (
	configFile     = flag.String("config", "", "Pfad zur Konfigurationsdatei")
	pcapFile       = flag.String("pcap", "", "Pfad zur PCAP-Datei für die Analyse")
	listenAddr     = flag.String("listen", "", "Adresse und Port zum Lauschen (überschreibt Konfiguration)")
	debug          = flag.Bool("debug", false, "Debug-Modus aktivieren")
	liveCapture    = flag.Bool("live", false, "Aktiviere Live-Capture")
	interface_name = flag.String("interface", "", "Netzwerkschnittstelle für Live-Capture")
)

// Globale Variablen für aktive WebSocket-Verbindungen
var (
	activeWebSockets = make(map[*websocket.Conn]bool)
	wsLock           = make(chan bool, 1) // Ein einfacher Mutex für Zugriff auf activeWebSockets
)

func main() {
	flag.Parse()

	// Konfiguration laden
	var cfg *config.Config
	var err error

	if *configFile != "" {
		cfg, err = config.LoadConfig(*configFile)
		if err != nil {
			log.Fatalf("Fehler beim Laden der Konfiguration: %v", err)
		}
	} else {
		// Standardkonfiguration verwenden
		cfg = config.DefaultConfig()
	}

	// Befehlszeilenargumente überschreiben Konfiguration
	if *listenAddr != "" {
		cfg.Server.Host = *listenAddr
	}

	// Live-Capture-Konfiguration
	if *liveCapture {
		cfg.Capture.EnableLive = true
	}

	if *interface_name != "" {
		cfg.Capture.Interface = *interface_name
	}

	// Ausgabeverzeichnisse erstellen
	createDirs(cfg)

	// Signalbehandlung für sauberes Herunterfahren
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(cancel)

	// PCAP-Capturer erstellen
	capturer := packet.NewPcapCapturer(cfg)
	defer capturer.Close()

	// API-Router initialisieren
	router := mux.NewRouter()

	// API-Handler registrieren
	registerAPIHandlers(router, capturer, cfg)

	// Statische Dateien bereitstellen
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.Server.StaticDir)))

	// Server starten
	listenAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	go func() {
		log.Printf("Server gestartet auf %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Serverfehler: %v", err)
		}
	}()

	// PCAP-Datei verarbeiten, falls angegeben
	if *pcapFile != "" {
		log.Printf("Analysiere PCAP-Datei: %s", *pcapFile)

		err := capturer.OpenPcapFile(*pcapFile)
		if err != nil {
			log.Fatalf("Fehler beim Öffnen der PCAP-Datei: %v", err)
		}

		packetChan, errChan := capturer.StartCapture(ctx)

		// Pakete verarbeiten
		go processPackets(packetChan, errChan)
	} else if cfg.Capture.EnableLive && cfg.Capture.Interface != "" {
		// Live-Capture starten
		log.Printf("Starte Live-Capture auf Schnittstelle: %s", cfg.Capture.Interface)

		err := capturer.OpenLiveCapture(cfg.Capture.Interface)
		if err != nil {
			log.Fatalf("Fehler beim Öffnen der Netzwerkschnittstelle %s: %v",
				cfg.Capture.Interface, err)
		}

		packetChan, errChan := capturer.StartCapture(ctx)

		// Pakete live verarbeiten und an WebSockets streamen
		go processLivePackets(packetChan, errChan)
	}

	// Auf Kontext-Abbruch warten
	<-ctx.Done()

	// Server herunterfahren
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Fehler beim Herunterfahren des Servers: %v", err)
	}

	log.Println("Server erfolgreich beendet")
}

// addWebSocket fügt eine neue WebSocket-Verbindung zur aktiven Liste hinzu
func addWebSocket(conn *websocket.Conn) {
	wsLock <- true
	activeWebSockets[conn] = true
	<-wsLock
}

// removeWebSocket entfernt eine WebSocket-Verbindung aus der aktiven Liste
func removeWebSocket(conn *websocket.Conn) {
	wsLock <- true
	delete(activeWebSockets, conn)
	<-wsLock
}

// broadcastPacketInfo sendet ein PacketInfo an alle aktiven WebSockets
func broadcastPacketInfo(packet *models.PacketInfo) {
	wsLock <- true
	defer func() { <-wsLock }()

	// Nur Gateway-Pakete übertragen, um Bandbreite zu sparen
	if !packet.IsGatewayTraffic {
		return
	}

	// Erstelle eine vereinfachte Zusammenfassung für die Übertragung
	summary := models.PacketSummary{
		Timestamp:        packet.Timestamp,
		SourceIP:         packet.SourceIP.String(),
		DestinationIP:    packet.DestinationIP.String(),
		Protocol:         packet.Protocol,
		Length:           packet.Length,
		IsGatewayTraffic: packet.IsGatewayTraffic,
		Summary:          createPacketSummary(packet),
	}

	// An alle aktiven WebSockets senden
	for conn := range activeWebSockets {
		err := conn.WriteJSON(map[string]interface{}{
			"type": "packet",
			"data": summary,
		})

		if err != nil {
			log.Printf("Fehler beim Senden an WebSocket: %v", err)
			conn.Close()
			delete(activeWebSockets, conn)
		}
	}
}

// createPacketSummary erstellt eine menschenlesbare Zusammenfassung des Pakets
func createPacketSummary(packet *models.PacketInfo) string {
	summary := fmt.Sprintf("%s: %s → %s", packet.Protocol, packet.SourceIP, packet.DestinationIP)

	// Details je nach Protokoll hinzufügen
	switch packet.Protocol {
	case "DNS":
		if packet.DNSInfo != nil {
			if packet.DNSInfo.IsQuery {
				for _, q := range packet.DNSInfo.Queries {
					summary += fmt.Sprintf(", Query: %s", q.Name)
				}
			} else if packet.DNSInfo.IsAnswer {
				summary += ", DNS-Antwort"
			}
		}
	case "DHCP":
		if packet.DHCPInfo != nil {
			summary += fmt.Sprintf(", Typ: %s", packet.DHCPInfo.MessageType)
		}
	case "ARP":
		if packet.ARPInfo != nil {
			summary += fmt.Sprintf(", %s", packet.ARPInfo.Operation)
		}
	}

	return summary
}

// createDirs erstellt die erforderlichen Verzeichnisse
func createDirs(cfg *config.Config) {
	dirs := []string{
		filepath.Dir(cfg.Storage.Path),
		cfg.Capture.PCAPDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warnung: Verzeichnis %s konnte nicht erstellt werden: %v", dir, err)
		}
	}
}

// setupSignalHandler richtet die Signalbehandlung ein
func setupSignalHandler(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Signal empfangen: %v", sig)
		cancel()
	}()
}

// registerAPIHandlers registriert die API-Handler
func registerAPIHandlers(router *mux.Router, capturer *packet.PcapCapturer, cfg *config.Config) {
	// API-Unterrouter für /api-Pfade
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Websocket-Upgrade-Handler
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // In Produktion einschränken
		},
	}

	// API-Endpunkte
	apiRouter.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")

	// PCAP-Upload- und Analyse-Endpunkt
	apiRouter.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		api.AnalyzePcapHandler(w, r, capturer)
	}).Methods("POST")

	// Websocket-Endpunkt für Live-Updates
	apiRouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Fehler beim Upgrade zu WebSocket: %v", err)
			return
		}

		// WebSocket zur aktiven Liste hinzufügen
		addWebSocket(conn)

		// Verbindung verwalten
		go handleWebSocketConnection(conn)
	})

	// Spezifische Gateway-Analyse-Endpunkte
	apiRouter.HandleFunc("/gateways", api.GetGatewaysHandler).Methods("GET")
	apiRouter.HandleFunc("/traffic/gateway", api.GetGatewayTrafficHandler).Methods("GET")
	apiRouter.HandleFunc("/events/gateway", api.GetGatewayEventsHandler).Methods("GET")

	// Verfügbare Netzwerkschnittstellen auflisten
	apiRouter.HandleFunc("/interfaces", api.GetInterfacesHandler).Methods("GET")

	// Live-Capture starten/stoppen
	apiRouter.HandleFunc("/live/start", func(w http.ResponseWriter, r *http.Request) {
		api.StartLiveCaptureHandler(w, r, capturer)
	}).Methods("POST")

	apiRouter.HandleFunc("/live/stop", func(w http.ResponseWriter, r *http.Request) {
		api.StopLiveCaptureHandler(w, r)
	}).Methods("POST")

	// Remote-Agent-Management-Endpunkte
	apiRouter.HandleFunc("/agents", api.ListAgentsHandler).Methods("GET")
	apiRouter.HandleFunc("/agents/register", api.RegisterAgentHandler).Methods("POST")
	apiRouter.HandleFunc("/agents/unregister", api.UnregisterAgentHandler).Methods("POST")
	apiRouter.HandleFunc("/agents/heartbeat", api.HeartbeatHandler).Methods("POST")
	apiRouter.HandleFunc("/agents/capture/start", api.StartAgentCaptureHandler).Methods("POST")
	apiRouter.HandleFunc("/agents/capture/stop", api.StopAgentCaptureHandler).Methods("POST")
	apiRouter.HandleFunc("/agents/set-interface", api.SetInterfaceHandler).Methods("POST")

	// Status-Prüfung für Agents starten
	go api.CheckAgentsStatus()
}

// handleWebSocketConnection verwaltet eine WebSocket-Verbindung
func handleWebSocketConnection(conn *websocket.Conn) {
	defer conn.Close()
	defer removeWebSocket(conn)

	// WebSocket-Handler-Loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket-Fehler: %v", err)
			}
			break
		}

		// Einfache Echo-Antwort
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Fehler beim Senden der Nachricht: %v", err)
			break
		}
	}
}

// processPackets verarbeitet Pakete aus dem Kanal
func processPackets(packetChan <-chan *models.PacketInfo, errChan <-chan error) {
	var packetCount int
	var gatewayPackets int

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case p, ok := <-packetChan:
			if !ok {
				log.Printf("Paketverarbeitung abgeschlossen: %d Pakete (davon %d Gateway-Pakete)",
					packetCount, gatewayPackets)
				return
			}

			packetCount++
			if p.IsGatewayTraffic {
				gatewayPackets++
			}

			// Hier könnten wir Pakete weiter verarbeiten oder speichern

		case err, ok := <-errChan:
			if !ok {
				continue
			}
			log.Printf("Fehler bei der Paketverarbeitung: %v", err)

		case <-ticker.C:
			if packetCount > 0 {
				log.Printf("Verarbeitet: %d Pakete (davon %d Gateway-Pakete, %.1f%%)",
					packetCount, gatewayPackets, float64(gatewayPackets)/float64(packetCount)*100)
			}
		}
	}
}

// processLivePackets verarbeitet Pakete in Echtzeit und streamt sie an WebSockets
func processLivePackets(packetChan <-chan *models.PacketInfo, errChan <-chan error) {
	var packetCount int
	var gatewayPackets int

	// Statistik-Ticker
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Pakete in Echtzeit verarbeiten
	for {
		select {
		case p, ok := <-packetChan:
			if !ok {
				log.Printf("Live-Capture beendet: %d Pakete (davon %d Gateway-Pakete)",
					packetCount, gatewayPackets)
				return
			}

			packetCount++
			if p.IsGatewayTraffic {
				gatewayPackets++

				// Paket an alle WebSockets senden
				broadcastPacketInfo(p)

				// Paketkennzahlen sammeln für Dashboard-Updates
				// Hier könnte z.B. eine Aktualisierung von Statistiken erfolgen
			}

		case err, ok := <-errChan:
			if !ok {
				continue
			}
			log.Printf("Fehler bei der Live-Capture: %v", err)

		case <-ticker.C:
			// Periodische Stats-Updates
			if packetCount > 0 {
				gwPercentage := float64(gatewayPackets) / float64(packetCount) * 100
				log.Printf("Live-Capture läuft: %d Pakete (davon %d Gateway-Pakete, %.1f%%)",
					packetCount, gatewayPackets, gwPercentage)

				// Stats an alle WebSockets senden
				statsUpdate := map[string]interface{}{
					"type": "stats",
					"data": map[string]interface{}{
						"total_packets":      packetCount,
						"gateway_packets":    gatewayPackets,
						"gateway_percentage": gwPercentage,
						"timestamp":          time.Now(),
					},
				}

				wsLock <- true
				for conn := range activeWebSockets {
					if err := conn.WriteJSON(statsUpdate); err != nil {
						log.Printf("Fehler beim Senden der Statistik: %v", err)
						conn.Close()
						delete(activeWebSockets, conn)
					}
				}
				<-wsLock
			}
		}
	}
}
