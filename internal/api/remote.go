package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// RemoteAgent enthält Informationen zu einem Remote-Capture-Agent
type RemoteAgent struct {
	Name             string                   `json:"name"`
	URL              string                   `json:"url"`
	Status           string                   `json:"status"` // "online", "offline", "capturing"
	LastSeen         time.Time                `json:"last_seen"`
	Interfaces       []string                 `json:"interfaces"`
	InterfaceDetails []map[string]interface{} `json:"interface_details"`
	ActiveInterface  string                   `json:"active_interface"`
	Version          string                   `json:"version"`
	OS               string                   `json:"os"`
	Hostname         string                   `json:"hostname"`
}

// AgentRegistration enthält die Informationen für die Agentenregistrierung
type AgentRegistration struct {
	Name             string                   `json:"name"`
	URL              string                   `json:"url"`
	Interfaces       []string                 `json:"interfaces"`
	InterfaceDetails []map[string]interface{} `json:"interface_details"`
	Version          string                   `json:"version"`
	OS               string                   `json:"os"`
	Hostname         string                   `json:"hostname"`
}

var (
	// Verwaltung der registrierten Agents
	remoteAgents      = make(map[string]*RemoteAgent)
	remoteAgentsMutex sync.RWMutex
)

// RegisterAgentHandler verarbeitet die Registrierung eines Remote-Agents
func RegisterAgentHandler(w http.ResponseWriter, r *http.Request) {
	// Authentifizierung prüfen (falls vorhanden)
	// apiKey := r.Header.Get("X-API-Key")
	// TODO: Implementiere richtige Validierung des API-Keys

	// Request-Body parsen
	var reg AgentRegistration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ungültiges Anfrageformat")
		return
	}

	// Agent-Informationen validieren
	if reg.Name == "" || reg.URL == "" {
		respondWithError(w, http.StatusBadRequest, "Name und URL sind erforderlich")
		return
	}

	// Neuen Agent erstellen oder bestehenden aktualisieren
	agent := &RemoteAgent{
		Name:             reg.Name,
		URL:              reg.URL,
		Status:           "online",
		LastSeen:         time.Now(),
		Interfaces:       reg.Interfaces,
		InterfaceDetails: reg.InterfaceDetails,
		Version:          reg.Version,
		OS:               reg.OS,
		Hostname:         reg.Hostname,
	}

	// In der Map speichern
	remoteAgentsMutex.Lock()
	remoteAgents[reg.Name] = agent
	remoteAgentsMutex.Unlock()

	log.Printf("Agent '%s' registriert: %s", reg.Name, reg.URL)

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("Agent '%s' erfolgreich registriert", reg.Name),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UnregisterAgentHandler behandelt die Abmeldung eines Agents
func UnregisterAgentHandler(w http.ResponseWriter, r *http.Request) {
	// API-Key prüfen
	// apiKey := r.Header.Get("X-API-Key")
	// TODO: Implementiere richtige Validierung des API-Keys

	// Request-Body parsen
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ungültiges Anfrageformat")
		return
	}

	// Agent-Namen validieren
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Agent-Name ist erforderlich")
		return
	}

	// Agent aus der Map entfernen
	remoteAgentsMutex.Lock()
	delete(remoteAgents, req.Name)
	remoteAgentsMutex.Unlock()

	log.Printf("Agent '%s' abgemeldet", req.Name)

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("Agent '%s' erfolgreich abgemeldet", req.Name),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HeartbeatHandler verarbeitet Heartbeat-Anfragen von Agents
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	// API-Key prüfen
	// apiKey := r.Header.Get("X-API-Key")
	// TODO: Implementiere richtige Validierung des API-Keys

	// Request-Body parsen
	var req struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ungültiges Anfrageformat")
		return
	}

	// Agent-Namen validieren
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Agent-Name ist erforderlich")
		return
	}

	// Agent in der Map aktualisieren
	remoteAgentsMutex.Lock()
	agent, exists := remoteAgents[req.Name]
	if exists {
		agent.LastSeen = time.Now()
		if req.Status != "" {
			agent.Status = req.Status
		}
	}
	remoteAgentsMutex.Unlock()

	if !exists {
		respondWithError(w, http.StatusNotFound, "Agent nicht registriert")
		return
	}

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListAgentsHandler gibt eine Liste aller registrierten Agents zurück
func ListAgentsHandler(w http.ResponseWriter, r *http.Request) {
	// Alle Agents aus der Map abrufen
	remoteAgentsMutex.RLock()
	agents := make([]*RemoteAgent, 0, len(remoteAgents))
	for _, agent := range remoteAgents {
		agents = append(agents, agent)
	}
	remoteAgentsMutex.RUnlock()

	// Erfolgreiche Antwort senden
	response := APIResponse{
		Success: true,
		Data:    agents,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StartAgentCaptureHandler startet eine Capture auf einem Remote-Agent
func StartAgentCaptureHandler(w http.ResponseWriter, r *http.Request) {
	// Request-Body parsen
	var req struct {
		Name      string `json:"name"`
		Interface string `json:"interface"`
		Filter    string `json:"filter,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ungültiges Anfrageformat")
		return
	}

	// Agent-Namen validieren
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Agent-Name ist erforderlich")
		return
	}

	// Agent in der Map finden
	remoteAgentsMutex.RLock()
	agent, exists := remoteAgents[req.Name]
	remoteAgentsMutex.RUnlock()

	if !exists {
		respondWithError(w, http.StatusNotFound, "Agent nicht gefunden")
		return
	}

	// Capture-Anfrage an den Agent senden
	captureReq := map[string]string{
		"interface": req.Interface,
		"filter":    req.Filter,
	}
	jsonData, err := json.Marshal(captureReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fehler bei der JSON-Kodierung")
		return
	}

	// URL zusammensetzen
	url := fmt.Sprintf("%s/capture/start", agent.URL)

	// HTTP-Request senden mit den serialisierten Daten
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Fehler bei der Kommunikation mit dem Agent: %v", err))
		return
	}
	defer resp.Body.Close()

	// Antwort des Agents parsen
	var agentResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fehler beim Parsen der Agent-Antwort")
		return
	}

	// Status des Agents aktualisieren
	if agentResp.Success {
		remoteAgentsMutex.Lock()
		agent.Status = "capturing"
		remoteAgentsMutex.Unlock()
	}

	// Antwort des Agents weiterleiten
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentResp)
}

// StopAgentCaptureHandler stoppt eine Capture auf einem Remote-Agent
func StopAgentCaptureHandler(w http.ResponseWriter, r *http.Request) {
	// Request-Body parsen
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ungültiges Anfrageformat")
		return
	}

	// Agent-Namen validieren
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Agent-Name ist erforderlich")
		return
	}

	// Agent in der Map finden
	remoteAgentsMutex.RLock()
	agent, exists := remoteAgents[req.Name]
	remoteAgentsMutex.RUnlock()

	if !exists {
		respondWithError(w, http.StatusNotFound, "Agent nicht gefunden")
		return
	}

	// URL zusammensetzen
	url := fmt.Sprintf("%s/capture/stop", agent.URL)

	// HTTP-Request senden
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Post(url, "application/json", http.NoBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Fehler bei der Kommunikation mit dem Agent: %v", err))
		return
	}
	defer resp.Body.Close()

	// Antwort des Agents parsen
	var agentResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fehler beim Parsen der Agent-Antwort")
		return
	}

	// Status des Agents aktualisieren
	if agentResp.Success {
		remoteAgentsMutex.Lock()
		agent.Status = "online"
		remoteAgentsMutex.Unlock()
	}

	// Antwort des Agents weiterleiten
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentResp)
}

// CheckAgentsStatus prüft regelmäßig den Status der Agents und markiert inaktive als offline
func CheckAgentsStatus() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		remoteAgentsMutex.Lock()
		for name, agent := range remoteAgents {
			// Wenn ein Agent seit mehr als 2 Minuten keinen Heartbeat gesendet hat,
			// markieren wir ihn als offline
			if time.Since(agent.LastSeen) > 2*time.Minute {
				agent.Status = "offline"
				log.Printf("Agent '%s' ist offline (kein Heartbeat seit %v)", name, time.Since(agent.LastSeen))
			}
		}
		remoteAgentsMutex.Unlock()
	}
}

// SetInterfaceHandler verarbeitet Anfragen zum Setzen der aktiven Schnittstellte auf einem Agent
func SetInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	// Request-Body parsen
	var req struct {
		Name      string `json:"name"`      // Name des Agents
		Interface string `json:"interface"` // Name der zu aktivierenden Schnittstelle
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Ungültiges Anfrageformat")
		return
	}

	// Agent-Namen und Schnittstelle validieren
	if req.Name == "" || req.Interface == "" {
		respondWithError(w, http.StatusBadRequest, "Agent-Name und Schnittstelle sind erforderlich")
		return
	}

	// Agent in der Map finden
	remoteAgentsMutex.RLock()
	agent, exists := remoteAgents[req.Name]
	remoteAgentsMutex.RUnlock()

	if !exists {
		respondWithError(w, http.StatusNotFound, "Agent nicht gefunden")
		return
	}

	// Überprüfen, ob die angegebene Schnittstelle auf dem Agent existiert
	interfaceExists := false
	for _, ifName := range agent.Interfaces {
		if ifName == req.Interface {
			interfaceExists = true
			break
		}
	}

	if !interfaceExists {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Schnittstelle '%s' existiert nicht auf Agent '%s'", req.Interface, req.Name))
		return
	}

	// Anfrage an den Agent senden, um die Schnittstelle zu aktivieren
	setInterfaceReq := map[string]string{
		"interface": req.Interface,
	}
	jsonData, err := json.Marshal(setInterfaceReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fehler bei der JSON-Kodierung")
		return
	}

	// URL zusammensetzen
	url := fmt.Sprintf("%s/capture/set-interface", agent.URL)

	// HTTP-Request senden
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Fehler bei der Kommunikation mit dem Agent: %v", err))
		return
	}
	defer resp.Body.Close()

	// Antwort des Agents parsen
	var agentResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Fehler beim Parsen der Agent-Antwort")
		return
	}

	// Status des Agents aktualisieren
	if agentResp.Success {
		remoteAgentsMutex.Lock()
		agent.ActiveInterface = req.Interface
		remoteAgentsMutex.Unlock()
	}

	// Antwort des Agents weiterleiten
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentResp)
}
