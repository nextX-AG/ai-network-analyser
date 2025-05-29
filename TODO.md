# TODO Liste

## Aktuelle Probleme in Bearbeitung [🔄]

### Remote-Agents UI Probleme
- [✓] Filter-Funktion wird nicht im Remote-Agents Tab angezeigt
  - [✓] Test-Text aus AgentCard.jsx entfernt
  - [✓] Filter-Button und Panel sind korrekt implementiert
  - [✓] Filter-Komponente wird beim Klick auf den Button angezeigt
  - [✓] Filter-State-Management pro Agent ist implementiert

### Packet Capture Bugs
- [🔄] Packet-Zählung funktioniert nicht korrekt
  - [ ] WebSocket-Verbindung für Packet-Updates überprüfen
    - [ ] WebSocket-Handler im Backend analysieren
    - [ ] WebSocket-Event-Format überprüfen
    - [ ] WebSocket-Verbindung im Frontend debuggen
  - [ ] Packet Counter State in RemoteAgentsContainer debuggen
    - [ ] State-Update-Logik überprüfen
    - [ ] Event-Handler für Packet-Updates testen
  - [ ] Backend Packet-Zähler in PcapCapturer analysieren
    - [ ] Zähler-Logik überprüfen
    - [ ] Event-Emission testen
  - [ ] WebSocket Event Handler für Packet-Updates korrigieren
    - [ ] Event-Format standardisieren
    - [ ] Error-Handling implementieren
  - [ ] Packet Counter Reset-Logik bei Capture Start/Stop überarbeiten
    - [ ] Reset-Event implementieren
    - [ ] Frontend-State-Reset sicherstellen

### Nächste Schritte
1. Filter-Integration:
   - [ ] NetworkFilterPanel.jsx aus networkCapture in remoteAgents/components kopieren
   - [ ] Komponente für Agent-spezifische Verwendung anpassen
   - [ ] Filter-State in remoteAgents/hooks implementieren
   - [ ] Integration in AgentCard.jsx

2. Packet Counter Fix:
   - [ ] WebSocket-Verbindung in Browser Dev Tools analysieren
   - [ ] Backend Logs für Packet-Events aktivieren
   - [ ] Packet Counter State Management überprüfen
   - [ ] WebSocket Reconnect-Logik testen

## Integration der Filterfunktionalität in Remote-Agents UI

### Vorbereitende Analyse
- [🔄] Identifizierung der relevanten Komponenten im Remote-Agents-Tab
- [🔄] Analyse der bestehenden Agent-Karten-Struktur
- [🔄] Festlegung optimaler Positionierung der Filter-UI innerhalb der Agent-Karte

## API-Routen

### Phase 7: API-Überprüfung und Debugging
- [x] API-Routen im Code überprüfen
  - [x] Server-Routen-Registrierung analysieren
  - [x] API-Handler-Implementierung prüfen
  - [x] Middleware-Konfiguration überprüfen
- [x] Server-Logs auf Fehler untersuchen
  - [x] Start-up Logs analysieren
  - [x] Runtime-Fehler identifizieren
  - [x] CORS und Routing-Probleme prüfen
- [x] API-Dokumentation überprüfen
  - [x] Endpunkt-Definitionen validieren
  - [x] Route-Präfixe verifizieren
  - [x] API-Versioning prüfen
- [x] API-Tests durchführen
  - [x] Basis-Endpunkte testen (/api/, /api/status)
  - [x] Client-spezifische Endpunkte testen
  - [x] WebSocket-Verbindungen testen

### Implementierte Routen
- [x] GET /api/health - Systemstatus und Version
- [x] GET /api/interfaces - Verfügbare Netzwerkschnittstellen
- [x] POST /api/analyze - PCAP-Datei analysieren
- [x] POST /api/live/start - Live-Capture starten
- [x] POST /api/live/stop - Live-Capture stoppen
- [x] GET /api/gateways - Gateway-Informationen
- [x] GET /api/traffic/gateway - Gateway-Traffic-Statistiken
- [x] GET /api/events/gateway - Gateway-Ereignisse
- [x] WebSocket /api/ws - Echtzeit-Updates

### Remote-Agent-Routen
- [x] GET /api/agents - Liste aller Agents
- [x] POST /api/agents/register - Agent registrieren
- [x] POST /api/agents/unregister - Agent abmelden
- [x] POST /api/agents/heartbeat - Agent-Status aktualisieren
- [x] POST /api/agents/capture/start - Capture auf Agent starten
- [x] POST /api/agents/capture/stop - Capture auf Agent stoppen
- [x] POST /api/agents/set-interface - Interface auf Agent setzen

## Frontend-Tasks
- [x] API-Service-Integration testen
- [x] Fehlerbehandlung implementieren
- [x] Loading-States hinzufügen
- [x] Retry-Mechanismus für fehlgeschlagene Anfragen
- [x] WebSocket-Reconnect-Logik implementieren

## Backend-Tasks
- [ ] API-Key-Authentifizierung implementieren
- [ ] Rate-Limiting einführen
- [ ] Logging verbessern
- [ ] Fehlerbehandlung vereinheitlichen
- [ ] Performance-Monitoring hinzufügen

## Sicherheit
- [ ] HTTPS-Konfiguration überprüfen
- [ ] CORS-Einstellungen überprüfen
- [ ] Input-Validierung verstärken
- [ ] API-Dokumentation aktualisieren 