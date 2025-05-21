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

// AdminTemplate enthält die HTML-Vorlage für die Admin-Seite
const AdminTemplate = `
<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Netzwerk-Analyzer Remote Agent</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            margin: 0;
            padding: 20px;
            color: #333;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        h1 {
            color: #2c3e50;
            margin-top: 0;
        }
        .status-card {
            background-color: #ebf5fb;
            border-left: 4px solid #3498db;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 4px;
        }
        .status-item {
            margin-bottom: 10px;
        }
        .status-label {
            font-weight: bold;
            margin-right: 10px;
        }
        .config-form {
            margin-top: 20px;
        }
        .form-group {
            margin-bottom: 15px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        input[type="text"], select {
            width: 100%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
        }
        button {
            background-color: #3498db;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover {
            background-color: #2980b9;
        }
        .success-message {
            color: #27ae60;
            margin-top: 10px;
            display: none;
        }
        .error-message {
            color: #e74c3c;
            margin-top: 10px;
            display: none;
        }
        .interface-card {
            background-color: #f9f9f9;
            border: 1px solid #ddd;
            padding: 10px;
            margin-bottom: 10px;
            border-radius: 4px;
        }
        .interface-name {
            font-weight: bold;
        }
        .interface-details {
            margin-top: 5px;
            font-size: 14px;
            color: #666;
        }
        .bridge-interface {
            background-color: #e8f8f5;
            border-color: #2ecc71;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Netzwerk-Analyzer Remote Agent</h1>
        
        <div class="status-card">
            <div class="status-item">
                <span class="status-label">Status:</span>
                <span id="agent-status">{{.Status}}</span>
            </div>
            <div class="status-item">
                <span class="status-label">Name:</span>
                <span>{{.Name}}</span>
            </div>
            <div class="status-item">
                <span class="status-label">Erfasste Pakete:</span>
                <span id="packets-captured">{{.PacketsCaptured}}</span>
            </div>
            <div class="status-item">
                <span class="status-label">Aktive Schnittstelle:</span>
                <span id="active-interface">{{.Interface}}</span>
            </div>
            <div class="status-item">
                <span class="status-label">Server-Verbindung:</span>
                <span id="server-connection">{{if .Connected}}Verbunden{{else}}Nicht verbunden{{end}}</span>
            </div>
        </div>
        
        <h2>Konfiguration</h2>
        <form id="config-form" class="config-form">
            <div class="form-group">
                <label for="server-url">Hauptserver-URL:</label>
                <input type="text" id="server-url" name="server_url" value="{{.ServerURL}}" placeholder="http://server-ip:9090">
            </div>
            
            <div class="form-group">
                <label for="agent-name">Agent-Name:</label>
                <input type="text" id="agent-name" name="name" value="{{.Name}}" placeholder="up-board-agent">
            </div>
            
            <div class="form-group">
                <label for="interface">Erfassungsschnittstelle:</label>
                <select id="interface" name="interface">
                    {{range .Interfaces}}
                    <option value="{{.Name}}" {{if eq $.Interface .Name}}selected{{end}}>{{.Name}} {{if .IsBridge}}(Bridge){{end}} - {{.IPs}}</option>
                    {{end}}
                </select>
            </div>
            
            <div class="form-group">
                <label for="api-key">API-Schlüssel:</label>
                <input type="text" id="api-key" name="api_key" value="{{.APIKey}}" placeholder="Optional: Authentifizierungsschlüssel">
            </div>
            
            <button type="submit">Konfiguration speichern</button>
            <div id="success-message" class="success-message">Konfiguration erfolgreich gespeichert!</div>
            <div id="error-message" class="error-message"></div>
        </form>
        
        <h2>Netzwerkschnittstellen</h2>
        <div id="interfaces-list">
            {{range .Interfaces}}
            <div class="interface-card {{if .IsBridge}}bridge-interface{{end}}">
                <div class="interface-name">{{.Name}} {{if .IsBridge}}(Bridge){{end}}</div>
                <div class="interface-details">
                    <div><strong>MAC:</strong> {{.MAC}}</div>
                    <div><strong>IP-Adressen:</strong> {{.IPs}}</div>
                    {{if .IsBridge}}
                    <div><strong>Bridge-Ports:</strong> {{.BridgePorts}}</div>
                    {{end}}
                </div>
            </div>
            {{end}}
        </div>
        
        <h2>Aktionen</h2>
        <div>
            <button id="restart-button">Agent neustarten</button>
            <button id="register-button">Bei Server registrieren</button>
        </div>
    </div>
    
    <script>
        // Formular absenden
        document.getElementById('config-form').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const formData = {
                server_url: document.getElementById('server-url').value,
                name: document.getElementById('agent-name').value,
                interface: document.getElementById('interface').value,
                api_key: document.getElementById('api-key').value
            };
            
            fetch('/admin/config', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            })
            .then(response => response.json())
            .then(data => {
                const successMsg = document.getElementById('success-message');
                const errorMsg = document.getElementById('error-message');
                
                if (data.success) {
                    successMsg.style.display = 'block';
                    errorMsg.style.display = 'none';
                    setTimeout(() => {
                        successMsg.style.display = 'none';
                    }, 3000);
                } else {
                    errorMsg.textContent = data.error || 'Fehler beim Speichern der Konfiguration';
                    errorMsg.style.display = 'block';
                    successMsg.style.display = 'none';
                }
            })
            .catch(err => {
                const errorMsg = document.getElementById('error-message');
                errorMsg.textContent = 'Fehler bei der Kommunikation mit dem Server: ' + err.message;
                errorMsg.style.display = 'block';
            });
        });
        
        // Neustart-Button
        document.getElementById('restart-button').addEventListener('click', function() {
            if (confirm('Möchten Sie den Agent wirklich neustarten?')) {
                fetch('/admin/restart', {
                    method: 'POST'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        alert('Agent wird neu gestartet...');
                    } else {
                        alert('Fehler beim Neustarten: ' + (data.error || 'Unbekannter Fehler'));
                    }
                })
                .catch(err => {
                    alert('Fehler bei der Kommunikation mit dem Server: ' + err.message);
                });
            }
        });
        
        // Registrierungs-Button
        document.getElementById('register-button').addEventListener('click', function() {
            fetch('/admin/register', {
                method: 'POST'
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Registrierung erfolgreich!');
                    document.getElementById('server-connection').textContent = 'Verbunden';
                } else {
                    alert('Fehler bei der Registrierung: ' + (data.error || 'Unbekannter Fehler'));
                }
            })
            .catch(err => {
                alert('Fehler bei der Kommunikation mit dem Server: ' + err.message);
            });
        });
        
        // Status-Updates
        function updateStatus() {
            fetch('/admin/status')
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('agent-status').textContent = data.data.status;
                    document.getElementById('packets-captured').textContent = data.data.packets_captured;
                    document.getElementById('active-interface').textContent = data.data.interface;
                    document.getElementById('server-connection').textContent = 
                        data.data.connected ? 'Verbunden' : 'Nicht verbunden';
                }
            })
            .catch(err => console.error('Fehler beim Abrufen des Status:', err));
        }
        
        // Status alle 5 Sekunden aktualisieren
        setInterval(updateStatus, 5000);
    </script>
</body>
</html>
`

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
	// Admin-Hauptseite
	router.HandleFunc("/admin", a.adminHandler).Methods("GET")

	// API-Endpunkte für Admin-Aktionen
	router.HandleFunc("/admin/config", a.configHandler).Methods("POST")
	router.HandleFunc("/admin/status", a.adminStatusHandler).Methods("GET")
	router.HandleFunc("/admin/restart", a.restartHandler).Methods("POST")
	router.HandleFunc("/admin/register", a.registerHandler).Methods("POST")
}

// adminHandler zeigt die Admin-Weboberfläche an
func (a *CaptureAgent) adminHandler(w http.ResponseWriter, r *http.Request) {
	// Template erstellen
	tmpl, err := template.New("admin").Parse(AdminTemplate)
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

	// Konfiguration speichern
	if err := a.saveConfig(); err != nil {
		log.Printf("Fehler beim Speichern der Konfiguration: %v", err)
		respondWithErrorJSON(w, fmt.Sprintf("Konfiguration konnte nicht gespeichert werden: %v", err))
		return
	}

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
	}

	// Weitere Pfade hinzufügen, die wir versuchen werden, nach Priorität sortiert
	execDir, _ := os.Executable()
	execDir = filepath.Dir(execDir)

	configPaths = append(configPaths,
		filepath.Join(execDir, "configs", "agent.json"),                  // /opt/ki-network-analyzer/configs/agent.json
		"/etc/ki-network-analyzer/agent.json",                            // Standard-Systemkonfiguration
		filepath.Join(execDir, "agent.json"),                             // Direkt im Executable-Verzeichnis
		filepath.Join(os.TempDir(), "ki-network-analyzer", "agent.json"), // Temp-Verzeichnis als letzten Ausweg
	)

	// Variable für den letzten Fehler
	var lastErr error

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

		// Wenn die Datei neu erstellt wurde, Berechtigungen setzen
		if !fileExists {
			if err := os.Chmod(configPath, 0664); err != nil { // rw-rw-r--
				log.Printf("Warnung: Konnte Berechtigungen nicht setzen: %v", err)
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
