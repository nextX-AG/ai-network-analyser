# KI-Netzwerk-Analyzer

Ein intelligentes Werkzeug zur Analyse von Netzwerkverkehr mit spezialisiertem Fokus auf Gateway-Analyse und KI-Integration.

## Übersicht

Der KI-Netzwerk-Analyzer ist eine modulare Plattform zur intelligenten Analyse von Netzwerkverkehr. Mit einem besonderen Fokus auf Gateway-Analyse ermöglicht das Tool:

- Erfassung und Analyse von PCAP/TCPDUMP-Dateien
- Automatische Erkennung von Gateway-Geräten im Netzwerk
- Analyse von Gateway-bezogenem Verkehr (DHCP, DNS, ARP, NAT)
- Visualisierung von Netzwerkereignissen mit Gateway-Fokus
- Integration von KI für erweiterte Musteranalyse (in zukünftigen Versionen)
- Remote-Capture-Funktionalität für verteilte Netzwerkerfassung

## Funktionen des MVP

Die aktuelle Version (MVP) bietet folgende Kernfunktionen:

- Einlesen und Analysieren von PCAP-Dateien über die Web-Oberfläche
- Automatische Erkennung von Gateway-Verkehr
- Identifikation von DHCP-, DNS- und ARP-bezogenen Gateway-Interaktionen
- Zusammenfassung wichtiger Gateway-Aktivitäten
- Benutzerfreundliche Web-Oberfläche
- Echtzeit-Netzwerkerfassung und -Analyse

## Remote-Capture-System

Der KI-Netzwerk-Analyzer unterstützt ein verteiltes Erfassungssystem mit Remote-Capture-Agents, die auf verschiedenen Geräten im Netzwerk ausgeführt werden können.

### Remote-Agent-Installation

Der Remote-Capture-Agent kann auf jedem Linux-System installiert werden, insbesondere auf UP Boards mit Ubuntu:

1. Installieren Sie die erforderlichen Bibliotheken:
   ```bash
   sudo apt-get update
   sudo apt-get install -y golang-go libpcap-dev git
   ```

2. Klonen Sie das Repository:
   ```bash
   git clone https://github.com/nextX-AG/ai-network-analyser
   cd ai-network-analyser
   ```

3. Erstellen Sie den Agent:
   ```bash
   go build -o bin/agent cmd/agent/main.go
   ```

4. Konfigurieren Sie den Agent:
   ```bash
   cp configs/agent.example.json configs/agent.json
   ```

5. Bearbeiten Sie die Konfigurationsdatei und passen Sie die Parameter an, insbesondere:
   - `server_url`: URL des Hauptservers
   - `interface`: Name der Netzwerkschnittstelle für die Paketerfassung
   - `name`: Eindeutiger Name für den Agent
   - `api_key`: Authentifizierungsschlüssel (falls aktiviert)

### Agent starten

```bash
sudo ./bin/agent --config=configs/agent.json
```

Sudo-Rechte werden benötigt, um auf die Netzwerkschnittstellen zuzugreifen.

### Agent-Verwaltung im Hauptsystem

1. Starten Sie die Hauptanwendung:
   ```bash
   ./bin/analyzer --config=configs/config.example.json
   ```

2. Öffnen Sie die Web-Oberfläche unter http://localhost:9090
3. Wechseln Sie zum Tab "Remote-Agents"
4. Klicken Sie auf "Agents aktualisieren", um verfügbare Agents zu sehen
5. Wählen Sie einen Agent aus und starten Sie die Erfassung auf diesem Gerät

### Automatischer Start als Systemdienst

Um den Agent als Systemdienst einzurichten (für automatischen Start beim Booten):

1. Erstellen Sie eine Systemd-Servicedatei:
   ```bash
   sudo nano /etc/systemd/system/network-agent.service
   ```

2. Fügen Sie folgenden Inhalt ein:
   ```
   [Unit]
   Description=Network Analyzer Remote Agent
   After=network.target

   [Service]
   ExecStart=/pfad/zum/bin/agent --config=/pfad/zur/configs/agent.json
   WorkingDirectory=/pfad/zum/projektverzeichnis
   User=root
   Restart=always
   RestartSec=10

   [Install]
   WantedBy=multi-user.target
   ```

3. Aktivieren und starten Sie den Dienst:
   ```bash
   sudo systemctl enable network-agent
   sudo systemctl start network-agent
   ```

4. Überprüfen Sie den Status des Dienstes:
   ```bash
   sudo systemctl status network-agent
   ```

## Installation

### Voraussetzungen

- Go 1.19 oder höher
- libpcap-Entwicklungsbibliotheken (für Linux/Mac)
  - Ubuntu/Debian: `sudo apt-get install libpcap-dev`
  - macOS: `brew install libpcap`
  - Windows: [WinPcap Developer Pack](https://www.winpcap.org/devel.htm) oder [Npcap SDK](https://nmap.org/npcap/)

### Bauen aus dem Quellcode

1. Repository klonen:
   ```
   git clone https://github.com/sayedamirkarim/ki-network-analyzer
   cd ki-network-analyzer
   ```

2. Abhängigkeiten installieren:
   ```
   go mod tidy
   ```

3. Anwendung bauen:
   ```
   go build -o bin/analyzer cmd/server/main.go
   ```

## Verwendung

### Starten der Anwendung

```
./bin/analyzer --config=configs/config.example.json
```

### Befehlszeilenoptionen

- `--config`: Pfad zur Konfigurationsdatei (optional)
- `--pcap`: Pfad zu einer PCAP-Datei, die sofort analysiert werden soll (optional)
- `--listen`: Adresse und Port zum Lauschen, z.B. `--listen=0.0.0.0:8080` (überschreibt Konfiguration)
- `--debug`: Debug-Modus aktivieren
- `--live`: Aktiviert Live-Capture-Modus
- `--interface`: Netzwerkschnittstelle für Live-Capture

### Web-Oberfläche

Nach dem Start ist die Web-Oberfläche unter http://localhost:9090 erreichbar (abhängig von der Konfiguration).

1. Navigieren Sie zur Web-Oberfläche
2. Laden Sie eine PCAP-Datei hoch oder ziehen Sie sie per Drag & Drop
3. Die Datei wird automatisch analysiert und Gateway-relevante Informationen werden angezeigt

### Live-Capture

1. Wechseln Sie zum "Live-Capture"-Tab in der Web-Oberfläche
2. Wählen Sie eine Netzwerkschnittstelle aus der Dropdown-Liste
3. Klicken Sie auf "Capture starten", um die Echtzeit-Analyse zu beginnen
4. Beobachten Sie Gateway-Traffic in Echtzeit

## Gateway-Analyse-Funktionen

Das System analysiert folgende Gateway-relevante Protokolle und Aktivitäten:

- **DHCP**: Lease-Anfragen, Gateway-Informationen in DHCP-Antworten
- **DNS**: DNS-Anfragen und -Antworten durch Gateway oder DNS-Server
- **ARP**: Gateway-ARP-Ankündigungen, ARP-Auflösungen für Gateway-Adressen
- **NAT**: Netzwerk-Adressübersetzung (derzeit vereinfacht implementiert)

## API-Endpunkte

- `GET /api/health`: Statusüberwachung
- `POST /api/analyze`: PCAP-Datei hochladen und analysieren
- `GET /api/gateways`: Liste erkannter Gateways abrufen
- `GET /api/traffic/gateway`: Gateway-Verkehrsstatistiken
- `GET /api/events/gateway`: Gateway-relevante Ereignisse
- `GET /api/interfaces`: Verfügbare Netzwerkschnittstellen
- `POST /api/live/start`: Live-Erfassung starten
- `POST /api/live/stop`: Live-Erfassung stoppen

## Projektstruktur

```
ai-network-analyser/
├── cmd/                  # Hauptanwendungen
│   └── server/           # Web-Server-Implementierung
├── configs/              # Konfigurationsdateien
├── docs/                 # Dokumentation
├── internal/             # Interne Pakete
│   ├── api/              # API-Handler
│   ├── config/           # Konfigurationsstrukturen
│   ├── packet/           # Paketanalyse
│   └── storage/          # Datenspeicherung
├── pkg/                  # Wiederverwendbare Pakete
│   ├── models/           # Datenmodelle
│   └── utils/            # Hilfsfunktionen
└── web/                  # Web-Frontend
```

## Nächste Schritte

- Integration von KI-Funktionen zur Erkennung von Anomalien
- Timeline-basierte Visualisierung von Gateway-Ereignissen
- Netzwerkgraph-Visualisierung mit Gateway-Fokus
- Echtzeitanalyse von laufendem Netzwerkverkehr
- Speicherung historischer Daten in einer Datenbank
- Implementierung des Remote-Capture-Systems für verteilte Erfassung

## Mitwirken

Wir freuen uns über Mitwirkung am Projekt. Bitte lesen Sie dazu `CONTRIBUTING.md` und `CODING_CONVENTIONS.md`.

## Lizenz

[MIT-Lizenz](LICENSE)
