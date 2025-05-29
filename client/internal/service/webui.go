package agent

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sayedamirkarim/ki-network-analyzer/internal/config"
)

// InterfaceInfo enthält Informationen über eine Netzwerkschnittstelle für die UI
type InterfaceInfo struct {
	Name        string
	MAC         string
	IPs         string
	IsBridge    bool
	BridgePorts string
}

// AdminPageData enthält die Daten für die Admin-Seite
type AdminPageData struct {
	Name            string
	Status          string
	PacketsCaptured int
	Interface       string
	ServerURL       string
	APIKey          string
	Connected       bool
	Interfaces      []InterfaceInfo
}

// RegisterAdminHandlers registriert die HTTP-Handler für die Admin-Weboberfläche
func (a *CaptureAgent) RegisterAdminHandlers(router *mux.Router) {
	// Admin-Subrouter mit CORS-Middleware erstellen
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(a.corsMiddleware)

	// Statische Dateien bereitstellen
	fileServer := http.FileServer(http.Dir("ui/static"))
	adminRouter.PathPrefix("/static/").Handler(http.StripPrefix("/admin/static/", fileServer))

	// Admin-Hauptseite
	adminRouter.HandleFunc("", a.adminHandler).Methods("GET")

	// API-Endpunkte für Admin-Aktionen
	adminRouter.HandleFunc("/config", a.configHandler).Methods("POST")
	adminRouter.HandleFunc("/status", a.adminStatusHandler).Methods("GET")
	adminRouter.HandleFunc("/restart", a.restartHandler).Methods("POST")
	adminRouter.HandleFunc("/register", a.registerHandler).Methods("POST")
}

// adminHandler zeigt die Admin-Weboberfläche an
func (a *CaptureAgent) adminHandler(w http.ResponseWriter, r *http.Request) {
	// Template laden
	tmpl, err := template.ParseFiles("ui/templates/admin.html")
	if err != nil {
		log.Printf("Error parsing admin template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Netzwerkschnittstellen für die Anzeige aufbereiten
	var interfaces []InterfaceInfo

	// Netzwerkschnittstellen erkennen
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting network interfaces: %v", err)
	} else {
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

			// IP-Adressen als Zeichenkette formatieren
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

			interfaces = append(interfaces, InterfaceInfo{
				Name:        iface.Name,
				MAC:         iface.HardwareAddr.String(),
				IPs:         strings.Join(ipStrings, ", "),
				IsBridge:    isBridge,
				BridgePorts: bridgePorts,
			})
		}
	}

	// Daten für die Template-Verarbeitung vorbereiten
	a.statusMutex.RLock()
	data := AdminPageData{
		Name:            a.config.Agent.Name,
		Status:          a.status.Status,
		PacketsCaptured: a.status.PacketsCaptured,
		Interface:       a.status.Interface,
		ServerURL:       a.config.Agent.ServerURL,
		APIKey:          a.config.Agent.APIKey,
		Connected:       a.status.Status != "error",
		Interfaces:      interfaces,
	}
	a.statusMutex.RUnlock()

	// Template rendern
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing admin template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// configHandler verarbeitet Konfigurationsänderungen über die Admin-Oberfläche
func (a *CaptureAgent) configHandler(w http.ResponseWriter, r *http.Request) {
	// Anfrage-Body parsen
	var req struct {
		ServerURL string `json:"server_url"`
		Name      string `json:"name"`
		Interface string `json:"interface"`
		APIKey    string `json:"api_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithErrorJSON(w, "Ungültiges Anfrageformat")
		return
	}

	// Konfiguration aktualisieren
	a.config.Agent.ServerURL = req.ServerURL
	a.config.Agent.Name = req.Name
	a.config.Agent.Interface = req.Interface
	a.config.Agent.APIKey = req.APIKey

	// Status aktualisieren - Wichtig für sofortige UI-Updates
	a.statusMutex.Lock()
	a.status.Name = req.Name
	a.status.Interface = req.Interface
	a.statusMutex.Unlock()

	// Konfiguration speichern
	if err := a.saveConfig(); err != nil {
		log.Printf("Fehler beim Speichern der Konfiguration: %v", err)
		respondWithErrorJSON(w, fmt.Sprintf("Konfiguration konnte nicht gespeichert werden: %v", err))
		return
	}

	// Auch den Status im Capturer aktualisieren, damit die Änderungen beim nächsten Neustart erhalten bleiben
	a.capturer.UpdateInterface(req.Interface)

	// Erfolgreiche Antwort senden
	respondWithSuccessJSON(w, "Konfiguration erfolgreich gespeichert", nil)
}

// saveConfig speichert die Konfiguration in die Konfigurationsdatei
func (a *CaptureAgent) saveConfig() error {
	// Liste der Pfade, wo wir versuchen könnten, die Konfiguration zu speichern
	configPaths := []string{}

	// Nach dem Argument "--config" suchen
	configArgPath := ""
	for i, arg := range os.Args {
		if arg == "--config" && i+1 < len(os.Args) {
			configArgPath = os.Args[i+1]
			break
		} else if strings.HasPrefix(arg, "--config=") {
			configArgPath = strings.TrimPrefix(arg, "--config=")
			break
		}
	}

	// Wenn ein Konfigurationspfad als Argument übergeben wurde, versuchen wir zuerst dort
	if configArgPath != "" {
		configPaths = append(configPaths, configArgPath)
		log.Printf("Konfigurationspfad aus Befehlszeile: %s", configArgPath)
	}

	// Weitere Pfade hinzufügen, die wir versuchen werden, nach Priorität sortiert
	execDir, _ := os.Executable()
	execDir = filepath.Dir(execDir)

	additionalPaths := []string{
		filepath.Join(execDir, "configs", "agent.json"),
		"/etc/ki-network-analyzer/agent.json",
		filepath.Join(execDir, "agent.json"),
		filepath.Join(os.TempDir(), "ki-network-analyzer", "agent.json"),
	}

	for _, path := range additionalPaths {
		// Prüfe, ob die Datei bereits existiert
		_, err := os.Stat(path)
		exists := !os.IsNotExist(err)

		if exists {
			log.Printf("Vorhandene Konfigurationsdatei gefunden: %s", path)
		}

		configPaths = append(configPaths, path)
	}

	// Variable für den letzten Fehler
	var lastErr error

	// Variable für erfolgreichen Pfad
	var successPath string

	// Versuche alle Pfade nacheinander
	for _, configPath := range configPaths {
		log.Printf("Versuche Konfiguration zu speichern in: %s", configPath)

		// Konfigurationsverzeichnis erstellen, falls nicht vorhanden
		configDir := filepath.Dir(configPath)
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				log.Printf("Konnte Konfigurationsverzeichnis nicht erstellen: %v", err)
				lastErr = err
				continue
			}
		}

		// Prüfen, ob die Datei existiert und beschreibbar ist
		var fileExists bool
		if _, err := os.Stat(configPath); err == nil {
			fileExists = true
			// Versuche, die Datei zum Schreiben zu öffnen
			testFile, err := os.OpenFile(configPath, os.O_WRONLY, 0)
			if err == nil {
				testFile.Close()
			} else {
				log.Printf("Datei existiert, ist aber nicht beschreibbar: %v", err)
				lastErr = err
				continue
			}
		}

		// Konfiguration speichern
		if err := config.SaveConfig(a.config, configPath); err != nil {
			log.Printf("Fehler beim Speichern der Konfiguration in %s: %v", configPath, err)
			lastErr = err
			continue
		}

		// Bei Erfolg eine Meldung ausgeben
		log.Printf("Konfiguration erfolgreich in %s gespeichert", configPath)
		successPath = configPath

		// Wenn die Datei neu erstellt wurde, Berechtigungen setzen
		if !fileExists {
			if err := os.Chmod(configPath, 0664); err != nil { // rw-rw-r--
				log.Printf("Warnung: Konnte Berechtigungen nicht setzen: %v", err)
			}
		}

		// Speichere den erfolgreichen Pfad für zukünftige Verwendung
		if successPath != "" {
			// Speichere den Pfad in einer Datei im Ausführungsverzeichnis
			configInfoPath := filepath.Join(execDir, "last_config_path")
			if err := os.WriteFile(configInfoPath, []byte(successPath), 0664); err != nil {
				log.Printf("Warnung: Konnte letzten Konfigurationspfad nicht speichern: %v", err)
			} else {
				log.Printf("Letzter erfolgreicher Konfigurationspfad gespeichert in: %s", configInfoPath)
			}
		}

		return nil
	}

	// Wenn wir hier ankommen, ist das Speichern fehlgeschlagen
	return fmt.Errorf("konnte Konfiguration in keinem der Pfade speichern: %v", lastErr)
}

// statusHandler gibt den aktuellen Status des Agents zurück
func (a *CaptureAgent) adminStatusHandler(w http.ResponseWriter, r *http.Request) {
	a.statusMutex.RLock()
	statusData := map[string]interface{}{
		"status":           a.status.Status,
		"packets_captured": a.status.PacketsCaptured,
		"interface":        a.status.Interface,
		"connected":        a.status.Status != "error",
	}
	a.statusMutex.RUnlock()

	respondWithSuccessJSON(w, "", statusData)
}

// restartHandler startet den Agent neu
func (a *CaptureAgent) restartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Agent restart requested")

	// Aktuelle Erfassung beenden, falls eine läuft
	if a.cancelFunc != nil {
		log.Println("Stopping current capture before restart")
		a.cancelFunc()
	}

	// Sicherstellen, dass die aktuelle Konfiguration gespeichert wird, bevor wir neustarten
	if err := a.saveConfig(); err != nil {
		log.Printf("Warnung: Konfiguration konnte vor Neustart nicht gespeichert werden: %v", err)
	} else {
		log.Println("Konfiguration vor Neustart gespeichert")
	}

	// Informieren Sie alle WebSocket-Clients über den Neustart
	a.clientsMutex.Lock()
	restartMessage := map[string]interface{}{
		"type":    "system",
		"message": "Agent wird neu gestartet...",
	}
	messageJSON, _ := json.Marshal(restartMessage)

	for client := range a.clients {
		client.WriteMessage(websocket.TextMessage, messageJSON)
		client.Close()
	}
	// Clients-Map leeren
	a.clients = make(map[*websocket.Conn]bool)
	a.clientsMutex.Unlock()

	// Erfolgreiche Antwort senden, bevor der Neustart beginnt
	respondWithSuccessJSON(w, "Neustart eingeleitet", nil)

	// Neustart als Goroutine ausführen, damit die HTTP-Antwort zuerst gesendet wird
	go func() {
		log.Println("Performing agent restart...")
		time.Sleep(1 * time.Second) // Kurze Verzögerung, um sicherzustellen, dass die Antwort gesendet wurde

		// In einer Produktionsumgebung würden wir systemd oder einen ähnlichen Dienst-Manager verwenden
		// Für Entwicklungszwecke starten wir den Agenten neu, indem wir den Prozess neu starten

		// Ermitteln Sie den Pfad zur ausführbaren Datei und den Argumenten
		executable, err := os.Executable()
		if err != nil {
			log.Printf("Error getting executable path: %v", err)
			return
		}

		// Neustart per Exec
		log.Printf("Restarting agent process: %s %v", executable, os.Args[1:])
		if err := syscall.Exec(executable, append([]string{executable}, os.Args[1:]...), os.Environ()); err != nil {
			log.Printf("Failed to restart agent: %v", err)
		}
	}()
}

// registerHandler registriert den Agent manuell beim Hauptserver
func (a *CaptureAgent) registerHandler(w http.ResponseWriter, r *http.Request) {
	// Sicherstellen, dass die aktuelle Konfiguration verwendet wird
	serverURL := a.config.Agent.ServerURL
	log.Printf("Verwende Server-URL für manuelle Registrierung: %s", serverURL)

	if err := a.Register(); err != nil {
		log.Printf("Registrierung fehlgeschlagen: %v", err)
		respondWithErrorJSON(w, fmt.Sprintf("Registrierung fehlgeschlagen: %v", err))
		return
	}

	// Status auf "idle" setzen bei erfolgreicher Registrierung
	a.statusMutex.Lock()
	a.status.Status = "idle"
	a.status.Error = ""
	a.statusMutex.Unlock()

	respondWithSuccessJSON(w, "Erfolgreich beim Hauptserver registriert", nil)
}

// Hilfsfunktionen für JSON-Antworten
func respondWithSuccessJSON(w http.ResponseWriter, message string, data interface{}) {
	response := map[string]interface{}{
		"success": true,
	}

	if message != "" {
		response["message"] = message
	}

	if data != nil {
		response["data"] = data
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func respondWithErrorJSON(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"success": false,
		"error":   message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
