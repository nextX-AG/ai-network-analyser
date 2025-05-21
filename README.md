# KI-Netzwerk-Analyzer

Ein intelligentes Werkzeug zur Analyse von Netzwerkverkehr mit spezialisiertem Fokus auf Gateway-Analyse und KI-Integration.

## Übersicht

Der KI-Netzwerk-Analyzer ist eine modulare Plattform zur intelligenten Analyse von Netzwerkverkehr. Mit einem besonderen Fokus auf Gateway-Analyse ermöglicht das Tool:

- Erfassung und Analyse von PCAP/TCPDUMP-Dateien
- Automatische Erkennung von Gateway-Geräten im Netzwerk
- Analyse von Gateway-bezogenem Verkehr (DHCP, DNS, ARP, NAT)
- Visualisierung von Netzwerkereignissen mit Gateway-Fokus
- Integration von KI für erweiterte Musteranalyse (in zukünftigen Versionen)

## Funktionen des MVP

Die aktuelle Version (MVP) bietet folgende Kernfunktionen:

- Einlesen und Analysieren von PCAP-Dateien über die Web-Oberfläche
- Automatische Erkennung von Gateway-Verkehr
- Identifikation von DHCP-, DNS- und ARP-bezogenen Gateway-Interaktionen
- Zusammenfassung wichtiger Gateway-Aktivitäten
- Benutzerfreundliche Web-Oberfläche

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

### Web-Oberfläche

Nach dem Start ist die Web-Oberfläche unter http://localhost:9090 erreichbar (abhängig von der Konfiguration).

1. Navigieren Sie zur Web-Oberfläche
2. Laden Sie eine PCAP-Datei hoch oder ziehen Sie sie per Drag & Drop
3. Die Datei wird automatisch analysiert und Gateway-relevante Informationen werden angezeigt

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

## Mitwirken

Wir freuen uns über Mitwirkung am Projekt. Bitte lesen Sie dazu `CONTRIBUTING.md` und `CODING_CONVENTIONS.md`.

## Lizenz

[MIT-Lizenz](LICENSE)
