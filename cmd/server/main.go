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
	configFile = flag.String("config", "", "Pfad zur Konfigurationsdatei")
	pcapFile   = flag.String("pcap", "", "Pfad zur PCAP-Datei für die Analyse")
	listenAddr = flag.String("listen", "", "Adresse und Port zum Lauschen (überschreibt Konfiguration)")
	debug      = flag.Bool("debug", false, "Debug-Modus aktivieren")
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
		api.WebSocketHandler(w, r, upgrader)
	})

	// Spezifische Gateway-Analyse-Endpunkte
	apiRouter.HandleFunc("/gateways", api.GetGatewaysHandler).Methods("GET")
	apiRouter.HandleFunc("/traffic/gateway", api.GetGatewayTrafficHandler).Methods("GET")
	apiRouter.HandleFunc("/events/gateway", api.GetGatewayEventsHandler).Methods("GET")
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
