# AI Network Analyser - TODO Liste

## Priorität 1: Projektrestrukturierung

### Phase 1: Vorbereitung
- [x] Backup des aktuellen Projektstands erstellen
- [x] Git-Branch für Restrukturierung anlegen
- [x] Temporäres Arbeitsverzeichnis für Migration erstellen

### Phase 2: Neue Verzeichnisstruktur
- [x] Client-Verzeichnisstruktur erstellen (vormals Agent)
  - [x] cmd/
  - [x] configs/
  - [x] internal/
  - [x] ui/

- [x] Server-Verzeichnisstruktur erstellen
  - [x] cmd/
  - [x] configs/
  - [x] internal/
    - [x] ai/
    - [x] speech/
    - [x] timeline/
  - [x] ui/

- [x] Gemeinsame Verzeichnisse erstellen
  - [x] pkg/
    - [x] ai/
    - [x] common/
    - [x] models/
    - [x] protocol/
  - [x] docs/
  - [x] scripts/
  - [x] configs/

### Phase 3: Code-Migration
- [x] Client-Code migrieren (vormals Agent)
  - [x] Hauptanwendung
  - [x] Interne Pakete
  - [ ] UI-Komponenten

- [x] Server-Code migrieren
  - [x] Hauptanwendung
  - [x] API-Handler
  - [x] UI-Komponenten
  - [ ] KI-Module
  - [ ] Sprachverarbeitung
  - [ ] Timeline-Komponenten

- [x] Gemeinsamen Code migrieren
  - [x] Modelle (packet.go)
  - [x] Protokolle (dhcp_types.go)
  - [x] Utilities (common.go)
  - [ ] KI-Utilities

### Phase 4: KI-Integration
- [ ] LLM-Integration
  - [ ] Lokales Modell einbinden
  - [ ] OpenAI API anbinden
  - [ ] Modell-Verwaltung implementieren

- [ ] KI-Analysemodule
  - [ ] Paketanalyse-Modelle entwickeln
  - [ ] Anomalie-Erkennung implementieren
  - [ ] Musteranalyse integrieren

- [ ] KI-UI-Komponenten
  - [ ] Analyse-Dashboard erstellen
  - [ ] Modell-Verwaltung UI
  - [ ] Ergebnis-Visualisierung

### Phase 5: Sprachverarbeitung
- [ ] Whisper Integration
  - [ ] Whisper.cpp einbinden
  - [ ] Modell-Management implementieren
  - [ ] Performance-Optimierung

- [ ] Spracherkennungs-UI
  - [ ] Aufnahme-Interface
  - [ ] Transkriptions-Anzeige
  - [ ] Sprachsteuerung

- [ ] API-Integration
  - [ ] REST-Endpunkte
  - [ ] WebSocket-Streaming
  - [ ] Fehlerbehandlung

### Phase 6: Anpassungen
- [x] Import-Pfade aktualisieren
- [x] Build-Skripte anpassen
- [ ] Tests aktualisieren
- [ ] Dokumentation aktualisieren

### Phase 7: API-Überprüfung und Debugging
- [ ] API-Routen im Code überprüfen
  - [ ] Server-Routen-Registrierung analysieren
  - [ ] API-Handler-Implementierung prüfen
  - [ ] Middleware-Konfiguration überprüfen
- [ ] Server-Logs auf Fehler untersuchen
  - [ ] Start-up Logs analysieren
  - [ ] Runtime-Fehler identifizieren
  - [ ] CORS und Routing-Probleme prüfen
- [ ] API-Dokumentation überprüfen
  - [ ] Endpunkt-Definitionen validieren
  - [ ] Route-Präfixe verifizieren
  - [ ] API-Versioning prüfen
- [ ] API-Tests durchführen
  - [ ] Basis-Endpunkte testen (/api/, /api/status)
  - [ ] Client-spezifische Endpunkte testen
  - [ ] WebSocket-Verbindungen testen

### Phase 8: Validierung
- [x] Builds testen
- [ ] Unit-Tests durchführen
- [ ] Integration testen
- [x] UI-Tests durchführen
- [ ] KI-Modelle validieren
- [ ] Spracherkennung testen

### Phase 2: Dokumentation vervollständigen
- [x] README.md für Client-Komponente erstellen
- [x] README.md für Server-Komponente erstellen
- [ ] Inline-Dokumentation für Packages aktualisieren
- [ ] API-Dokumentation erstellen

## Priorität 2: Feature-Entwicklung

### KI-Features
- [ ] Implementierung der Paketanalyse-KI
- [ ] Integration der Anomalie-Erkennung
- [ ] Entwicklung der Musteranalyse
- [ ] Implementierung der KI-gestützten Filtervorschläge

### Sprachverarbeitung
- [ ] Implementierung der Sprachnotizen
- [ ] Integration der Sprachsteuerung
- [ ] Entwicklung der Sprachanalyse
- [ ] Implementierung der Sprachbasierten Suche

### Timeline-Features
- [ ] Implementierung der KI-gestützten Ereigniserkennung
- [ ] Integration der Sprachnotizen in Timeline
- [ ] Entwicklung der automatischen Ereignisklassifizierung

## Priorität 3: Optimierungen
- [ ] Performance-Optimierung der KI-Modelle
- [ ] Optimierung der Spracherkennung
- [ ] Verbesserung der Echtzeit-Analyse
- [ ] Optimierung der Datenspeicherung

## Abgeschlossene Aufgaben
- [x] Initiale Projektstruktur definiert
- [x] Architektur-Dokumentation erstellt
- [x] Restrukturierungsplan erstellt

# KI-Netzwerk-Analyzer TODO List

This document provides a high-level overview of the current tasks for the KI-Netzwerk-Analyzer project. For more detailed tasks and progress tracking, please see the [docs/TODO.md](docs/TODO.md).

## Current Focus Areas

### High Priority
- [x] Fix Remote Agent interface persistence issue - after restart selected interface is not preserved
- [x] Fix Remote Agent UI interface display - active interface is not shown in status
- [ ] Implement server-side network interface selection for agents
- [ ] Docker configuration for development environment
- [ ] SQLite integration for data persistence
- [ ] Implement optimizations for large PCAP files
- [ ] Complete the authentication and security concept for remote agents
- [ ] Implement Speech2Text module with Whisper.cpp integration

### Medium Priority
- [ ] Timeline visualization with Three.js
- [ ] AI annotation module with OpenAI GPT API integration
- [ ] Implement event and timeline module
- [ ] Extend test coverage for critical components
- [ ] Create Docker images for remote agents

### Low Priority
- [ ] Multi-agent capture synchronization
- [ ] Mobile view for frontend
- [ ] Prepare for AQEA compatibility
- [ ] CI/CD pipeline for automated tests and builds
- [ ] Cloud-based deployment option

## Completed Tasks
- [x] Project structure definition and initialization
- [x] Basic Go backend implementation
- [x] Integration of packet capture with gopacket
- [x] Basic React/Three.js frontend scaffold
- [x] Remote capture system for distributed capture
- [x] Web interface for agents
- [x] Automatic detection and registration of agents
- [x] Bridge optimization for MITM monitoring
- [x] Gateway detection and analysis implementation
- [x] REST API endpoints for gateway information
- [x] Systemd service templates for easy deployment
- [x] Fix Admin UI Route registration - made Admin UI accessible
- [x] Fix configuration file permissions - resolved read-only filesystem issue for configuration
- [x] Fix Remote Agent interface persistence and status display
- [x] Fix packet capturing permissions - ensured agent runs as root with proper capabilities

## Current Remote Agent Improvements
- [x] Added multiple configuration paths to handle read-only filesystems
- [x] Improved error handling in configuration saving
- [x] Updated installation script to use writable configuration paths
- [x] Added permission checks and fixes for configuration files
- [x] Fix interface persistence after agent restart
- [x] Fix active interface display in status UI
- [x] Added UpdateInterface method to PcapCapturer to ensure configuration is updated
- [x] Improved saveConfig function to try multiple paths if one fails
- [x] Ensured restart handler saves configuration before restarting
- [x] Added root permission check and explicit capability requirements in agent
- [x] Enhanced systemd service to ensure proper network capture permissions

## Server-Side Network Interface Selection Implementation
- [ ] Fix agent registration to use actual routable IP address instead of 0.0.0.0
- [ ] Enhance agent to send complete list of available network interfaces with details (IP, MAC, bridge status)
- [ ] Implement server-side API endpoint to select and activate interfaces on agents
- [ ] Update UI to display all available interfaces for each agent
- [ ] Add interface selection controls in Remote-Agents UI
- [ ] Implement WebSocket protocol for real-time capture status updates
- [ ] Add error handling for unreachable interfaces

## Current Agent Issues to Fix
- [x] Fix packet counter display in UI when packets are captured (Agent shows captured packets but UI doesn't)
- [x] Fix heartbeat mechanism to include captured packet count in status updates
- [x] Implement workaround for UI updating with real-time packet counts via polling
- [x] Ensure interface configuration is correctly persisted between agent restarts
- [x] Fix Server-URL configuration persistence and prioritization of saved values
- [x] Add detailed logging for agent configuration saving/loading process
- [x] Fix CORS issues with Agent API to allow cross-origin access from main server UI
- [ ] Implement proper error handling for WebSocket communication failures
- [ ] Add server-side packet counter validation against agent-reported values

## Next Action Items

1. Complete server-side network interface selection
2. Complete SQLite integration for data persistence
3. Implement the Speech2Text module
4. Begin development of the Three.js timeline visualization
5. Set up Docker configuration for development
6. Begin AI integration for packet analysis 

## Packet Filtering Implementation

### UI Components
- [x] Design and implement filter input section in server UI
- [x] Add text field for manual BPF filter syntax entry
- [x] Create UI for source/destination IP address filtering
- [x] Create UI for port filtering (source and destination)
- [x] Create UI for protocol filtering (TCP, UDP, ICMP, etc.)
- [x] Create UI for MAC address filtering
- [x] Implement filter combination mechanism (AND/OR operators)
- [x] Add filter presets for common use cases (HTTP/HTTPS, DNS, etc.)
- [x] Implement filter validation to prevent syntax errors
- [x] Create UI for saved filter management
- [ ] Integrate filter UI into Remote-Agents tab instead of separate page
- [ ] Update agent management UI to display active filters

### Server-Side Implementation
- [x] Extend API endpoints to accept filter parameters
- [x] Implement filter parameter validation on server
- [x] Create filter parser to convert UI filters to BPF syntax
- [x] Extend capture configuration to include filters
- [x] Implement filter state persistence in session

### Agent-Side Implementation
- [x] Extend agent capture API to accept BPF filter parameters
- [x] Apply BPF filters to PcapCapturer at capture start
- [x] Implement proper error handling for invalid filters
- [x] Add filter feedback mechanism to detect inefficient filters
- [x] Update agent status to include current active filter

### Testing and Documentation
- [ ] Create test cases for various filter combinations
- [ ] Document BPF syntax for advanced users
- [ ] Create example filters for common network analysis tasks
- [ ] Test filter performance on high-volume captures
- [ ] Document filter best practices in user guide 
- [ ] Document filter integration in agent management UI

### Integration Plan
- [🔄] Refactor NetworkFilterPanel to be agent-specific
- [🔄] Add filter section to each agent card in Remote-Agents UI
- [🔄] Ensure filter state is saved per agent
- [🔄] Implement filter synchronization between UI and agent status 

## Integration der Filterfunktionalität in Remote-Agents UI

### Vorbereitende Analyse
- [🔄] Identifizierung der relevanten Komponenten im Remote-Agents-Tab
- [🔄] Analyse der bestehenden Agent-Karten-Struktur
- [🔄] Festlegung optimaler Positionierung der Filter-UI innerhalb der Agent-Karte

### Refactoring der Filterkomponente
- [ ] Anpassung der NetworkFilterPanel-Komponente für agentenspezifische Verwendung
- [ ] Implementierung von Props für Agent-ID und Status-Weitergabe
- [ ] Optimierung der UI-Größe und Darstellung für kompakte Integration
- [ ] Erstellung eines Filter-Collapse-Panels pro Agent für erweiterbaren Filterbereich

### UI-Integration und Erweiterungen
- [ ] Integration des NetworkFilterPanels in die Agent-Detailansicht
- [ ] Erweiterung der Agent-Karte um Filterstatusanzeige (aktiver Filter)
- [ ] Hinzufügen von visueller Indikation für aktive Filter in der Agentenübersicht
- [ ] Anpassung der Filterfunktion an den aktuellen UI-Stil

### Datenspeicherung und Zustandsmanagement
- [ ] Implementierung eines agentenspezifischen Filter-Zustandsmanagements
- [ ] Persistenz der Filter-Einstellungen pro Agent im lokalen Speicher
- [ ] Synchronisation des Filter-Zustands mit dem Agentenstatus
- [ ] Erweiterung der Agentenverbindungsverwaltung für Filterübertragung

### Backend-Anpassungen
- [ ] Überprüfung der bestehenden API für agentenspezifische Filterverarbeitung
- [ ] Erweiterung der Agent-Status-API um Filterinformationen (falls erforderlich)
- [ ] Optimierung der Filter-Validierung für Echtzeitfeedback

### Benutzererfahrung und Dokumentation
- [ ] Hinzufügen von Tooltips und Hilfetexten zur Filter-Bedienung
- [ ] Erstellung einfacher Beispielfilter für häufige Anwendungsfälle
- [ ] Dokumentation der Filterfunktion in der Benutzeranleitung
- [ ] Erstellen von Fehlermeldungen bei ungültigen Filterausdrücken

### Tests und Qualitätssicherung
- [ ] Erstellen von Testfällen für die integrierte Filterfunktion
- [ ] Testen der UI auf verschiedenen Bildschirmgrößen
- [ ] Performance-Tests mit mehreren Agenten und aktiven Filtern
- [ ] Integration in bestehende Testsuite 

## Frontend-Refactoring: Feature-basierte Organisation

### Struktur und Setup
- [x] Erstellen der neuen Feature-basierten Verzeichnisstruktur im Frontend
- [x] Einrichten der `shared` Verzeichnisse für gemeinsam genutzte Komponenten
- [ ] Aktualisieren der Import-Pfade in Haupt-App-Dateien (App.jsx, index.jsx)
- [ ] Aktualisieren der Build-Konfiguration für die neue Struktur

### Feature: Remote Agents
- [x] Erstellen des `features/remoteAgents` Verzeichnisses mit Unterordnern
- [x] Migrieren der RemoteAgentsContainer.jsx in die neue Struktur
- [x] Extrahieren und Migrieren der Agent-Komponenten aus dem network-Verzeichnis
- [x] Erstellen einer neuen RemoteAgentsPage.jsx als Container-Komponente
- [x] Refaktorisieren der AgentCardWithFilter.jsx in kleinere, spezialisierte Komponenten
- [x] Migrieren und Refaktorisieren der Filter-Komponenten für Agent-spezifische Nutzung
- [x] Anpassen der API-Aufrufe und Services für die neue Struktur

### Feature: Network Capture
- [x] Erstellen des `features/networkCapture` Verzeichnisses mit Unterordnern
- [x] Migrieren der NetworkCapturePage.jsx in die neue Struktur
- [x] Migrieren und Refaktorisieren der NetworkCapturePanel.jsx
- [x] Erstellen spezialisierter Komponenten für Capture-Steuerung und -Anzeige
- [ ] Aktualisieren der Routen in der App für die neue Struktur

### Feature: Timeline
- [ ] Erstellen des `features/timeline` Verzeichnisses mit Unterordnern
- [ ] Migrieren existierender Timeline-Komponenten in die neue Struktur
- [ ] Entwickeln einer neuen TimelinePage.jsx als Container-Komponente
- [ ] Aktualisieren der Timeline-spezifischen Services und API-Aufrufe

### Feature: AI Integration
- [ ] Erstellen des `features/ai` Verzeichnisses mit Unterordnern
- [ ] Migrieren der KI-bezogenen Komponenten aus verschiedenen Verzeichnissen
- [ ] Entwickeln einer neuen AIAssistantPage.jsx als Container-Komponente
- [ ] Aktualisieren der KI-spezifischen Services und API-Aufrufe

### Shared-Komponenten
- [x] Identifizieren gemeinsam genutzter UI-Komponenten in der aktuellen Struktur
- [x] Migrieren dieser Komponenten in das `shared/components` Verzeichnis
- [x] Refaktorisieren gemeinsam genutzter Hooks in das `shared/hooks` Verzeichnis
- [x] Extrahieren und Migrieren gemeinsam genutzter Utilities
- [x] Aktualisieren aller Imports für die gemeinsam genutzten Komponenten

### Tests und Qualitätssicherung
- [ ] Erstellen von Komponententests für die migrierten Features
- [ ] Durchführen manueller Tests für alle migrierten Features
- [ ] Überprüfen aller Routen und Navigationen
- [ ] Sicherstellen der korrekten Funktionalität aller interaktiven Elemente
- [ ] Validieren der API-Integration und Datenflüsse

### Dokumentation
- [ ] Aktualisieren der Frontend-Dokumentation in ARCHITECTURE.md
- [ ] Erstellen von README-Dateien für jedes Feature-Verzeichnis
- [ ] Dokumentieren der Komponenten-API und Props
- [ ] Aktualisieren aller Code-Kommentare entsprechend der neuen Struktur
- [ ] Erstellen einer Migrationsdokumentation für zukünftige Entwickler

### Aufräumen
- [ ] Entfernen duplizierter Code nach der Migration
- [x] Entfernen nicht mehr benötigter Verzeichnisse und Dateien
- [ ] Optimieren der Import-Statements in allen Dateien
- [ ] Anpassen des Code-Styles für Konsistenz
- [ ] Durchführen einer finalen Code-Review des refaktorierten Codes

### Deployment und Integration
- [ ] Testen des Builds der refaktorierten Anwendung
- [ ] Überprüfen der Kompatibilität mit dem Backend
- [ ] Sicherstellen, dass die neue Struktur im Entwicklungsworkflow funktioniert
- [ ] Bereitstellen der aktualisierten Anwendung in der Testumgebung
- [ ] Validieren der End-to-End-Funktionalität nach dem Refactoring 

### Phase 1: Client-UI-Struktur vervollständigen
- [x] HTML-Templates aus Go-Code in `/client/ui/templates` extrahieren
- [x] CSS und JavaScript in `/client/ui/static` organisieren
- [x] Statische Assets in `/client/ui/public` verschieben
- [x] UI-Logik in `/client/internal/webui` überarbeiten

### Phase 5: Timeline-Feature
- [ ] Timeline-Komponenten entwickeln
- [ ] UI-Integration
- [ ] Event-Synchronisation 