package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/sayedamirkarim/ki-network-analyzer/internal/agent"
	"github.com/sayedamirkarim/ki-network-analyzer/internal/config"
)

var (
	configFile = flag.String("config", "", "Path to configuration file")
	listenAddr = flag.String("listen", "0.0.0.0:8090", "Address and port to listen on")
	serverAddr = flag.String("server", "http://localhost:9090", "Address of the main server")
	debug      = flag.Bool("debug", false, "Enable debug mode")
	interface_ = flag.String("interface", "", "Network interface to capture packets from")
	name       = flag.String("name", "", "Agent name (defaults to hostname)")
)

func main() {
	flag.Parse()

	// Check for root permissions
	if err := checkRootPermissions(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// Set up agent name
	agentName := *name
	if agentName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("Warning: Could not determine hostname: %v", err)
			agentName = "unknown-agent"
		} else {
			agentName = hostname
		}
	}

	// Determine config file path
	configFilePath := *configFile
	if configFilePath == "" {
		// Check for last known config path
		execDir, _ := os.Executable()
		execDir = filepath.Dir(execDir)
		lastConfigPath := filepath.Join(execDir, "last_config_path")

		if data, err := os.ReadFile(lastConfigPath); err == nil {
			savedPath := string(data)
			if _, err := os.Stat(savedPath); err == nil {
				configFilePath = savedPath
				log.Printf("Verwende gespeicherten Konfigurationspfad: %s", configFilePath)
			}
		}

		// Fallbacks prüfen, falls keine gespeicherte Konfiguration gefunden wurde
		if configFilePath == "" {
			potentialPaths := []string{
				filepath.Join(execDir, "configs", "agent.json"),
				"/etc/ki-network-analyzer/agent.json",
				filepath.Join(execDir, "agent.json"),
			}

			for _, path := range potentialPaths {
				if _, err := os.Stat(path); err == nil {
					configFilePath = path
					log.Printf("Verwende existierende Konfigurationsdatei: %s", configFilePath)
					break
				}
			}
		}
	}

	// Load configuration if specified
	var cfg *config.Config
	var err error
	if configFilePath != "" {
		log.Printf("Lade Konfiguration aus: %s", configFilePath)
		cfg, err = config.LoadConfig(configFilePath)
		if err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}
	} else {
		// Use default configuration
		log.Println("Verwende Standardkonfiguration (keine Konfigurationsdatei gefunden)")
		cfg = config.DefaultConfig()
	}

	// Override config with command line flags, but only if they are explicitly set
	if cfg.Agent == nil {
		cfg.Agent = &config.AgentConfig{
			Listen:    *listenAddr,
			ServerURL: *serverAddr,
			Interface: *interface_,
			Name:      agentName,
		}
	} else {
		// Werte nur überschreiben, wenn nicht leer
		if *listenAddr != "" {
			cfg.Agent.Listen = *listenAddr
		}
		if *serverAddr != "" && *serverAddr != "http://localhost:9090" {
			// Nur überschreiben, wenn explizit ein anderer Wert als der Standard angegeben wurde
			cfg.Agent.ServerURL = *serverAddr
		}
		if *interface_ != "" {
			cfg.Agent.Interface = *interface_
		}
		if *name != "" {
			cfg.Agent.Name = *name
		}
	}

	// Logge wichtige Konfigurationswerte für Debug-Zwecke
	log.Printf("Wichtige Konfigurationswerte:")
	log.Printf("- Server-URL: %s", cfg.Agent.ServerURL)
	log.Printf("- Interface: %s", cfg.Agent.Interface)
	log.Printf("- Agent-Name: %s", cfg.Agent.Name)

	// Create context for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handler
	setupSignalHandler(cancel)

	// Create and start the agent
	captureAgent := agent.NewCaptureAgent(cfg)
	if err := captureAgent.Init(); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// Register with the main server
	if err := captureAgent.Register(); err != nil {
		log.Printf("Warning: Failed to register with main server: %v", err)
	}

	// Set up the HTTP router
	router := mux.NewRouter()

	// Register API routes
	captureAgent.RegisterRoutes(router)

	// Register Admin UI routes
	captureAgent.RegisterAdminHandlers(router)

	// Set up and start the HTTP server
	server := &http.Server{
		Addr:    cfg.Agent.Listen,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting agent server on %s", cfg.Agent.Listen)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation (signal)
	<-ctx.Done()

	// Graceful shutdown
	log.Println("Shutting down agent...")

	// Unregister from the main server
	if err := captureAgent.Unregister(); err != nil {
		log.Printf("Warning: Failed to unregister from main server: %v", err)
	}

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Agent shutdown complete")
}

// setupSignalHandler sets up signal handling for graceful shutdown
func setupSignalHandler(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received signal: %v", sig)
		cancel()
	}()
}

// checkRootPermissions prüft, ob das Programm mit ausreichenden Rechten für die Paketerfassung läuft
func checkRootPermissions() error {
	if runtime.GOOS == "windows" {
		// Auf Windows müssen wir andere Überprüfungen durchführen
		// oder es wird normalerweise über die WinPcap/Npcap-Bibliothek geregelt
		return nil
	}

	// Auf Unix-ähnlichen Systemen prüfen wir, ob wir root sind
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("konnte aktuellen Benutzer nicht ermitteln: %v", err)
	}

	if currentUser.Uid != "0" {
		// Nicht als root, aber eventuell mit CAP_NET_RAW und CAP_NET_ADMIN Capabilities
		// Das kann man mit dem Befehl "getcap" überprüfen, aber das ist über Go schwieriger zu implementieren

		// Stattdessen eine Warnung ausgeben
		log.Println("WARNUNG: Der Agent läuft nicht als root-Benutzer.")
		log.Println("         Für die Paketerfassung werden root-Rechte oder bestimmte Capabilities benötigt.")
		log.Println("         Wenn die Paketzählung nicht funktioniert, starten Sie den Agent als root oder setzen Sie die notwendigen Capabilities:")
		log.Println("         sudo setcap 'cap_net_raw,cap_net_admin=eip' " + os.Args[0])

		// Versuchen, Schnittstellenzugriff zu testen
		iface := getSelectedInterface()
		if iface != "" {
			if err := testInterfaceAccess(iface); err != nil {
				return fmt.Errorf("unzureichende Rechte für Paketerfassung: %v", err)
			}
		}
	}

	return nil
}

// getSelectedInterface ermittelt die gewählte Schnittstelle aus den Kommandozeilen-Flags
func getSelectedInterface() string {
	if *interface_ != "" {
		return *interface_
	}

	// Versuchen, die Schnittstelle aus der Konfigurationsdatei zu lesen
	if *configFile != "" {
		cfg, err := config.LoadConfig(*configFile)
		if err == nil && cfg.Agent != nil && cfg.Agent.Interface != "" {
			return cfg.Agent.Interface
		}
	}

	return ""
}

// testInterfaceAccess versucht, die Schnittstelle für die Paketerfassung zu öffnen
func testInterfaceAccess(iface string) error {
	// Ein einfaches Kommando ausführen, um zu prüfen, ob wir auf die Schnittstelle zugreifen können
	// Dies ist ein vereinfachter Test, der in der Praxis nicht immer zuverlässig ist
	tempFile := os.TempDir() + "/pcap_test.log"
	cmd := exec.Command("tcpdump", "-i", iface, "-c", "1", "-w", tempFile)

	// Standardausgabe und -fehler erfassen
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Kommando mit Timeout ausführen
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("konnte tcpdump nicht starten: %v", err)
	}

	// Mit Timeout warten
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(2 * time.Second):
		// Timeout - tcpdump hängt oder wartet auf Pakete
		cmd.Process.Kill()
		return nil // Kein offensichtlicher Fehler
	case err := <-done:
		// Fertig
		if err != nil {
			// Fehler prüfen
			errorOutput := stderr.String()
			if strings.Contains(errorOutput, "permission denied") ||
				strings.Contains(errorOutput, "permissions") ||
				strings.Contains(errorOutput, "Operation not permitted") {
				return fmt.Errorf("keine Berechtigung für Paketerfassung auf Schnittstelle %s", iface)
			}
		}
	}

	// Temporäre Datei entfernen
	os.Remove(tempFile)
	return nil
}
