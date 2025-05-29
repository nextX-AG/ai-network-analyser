# AI Network Analyser - Client

## Übersicht

Die Client-Komponente des AI Network Analysers ist verantwortlich für:
- Paketerfassung auf Remote-Systemen
- Lokale Analyse von Netzwerkdaten
- Leichtgewichtige Web-UI für Monitoring und Konfiguration

## Struktur

```
client/
├── cmd/                  # Hauptanwendung
├── configs/              # Client-Konfigurationen
├── internal/             # Private Pakete
│   ├── capture/         # Paketerfassung
│   ├── service/         # Client-Dienste
│   └── webui/          # Web-UI Backend
└── ui/                  # Frontend
    ├── public/          # Öffentliche Dateien
    ├── templates/       # HTML-Templates
    └── static/          # Statische Assets
        ├── css/         # Stylesheets
        ├── js/          # JavaScript
        └── images/      # Bilder
```

## Komponenten

### Capture
- Implementiert die Paketerfassung mittels gopacket
- Unterstützt verschiedene Netzwerk-Interfaces
- Bietet Filterung und Vorverarbeitung

### Service
- Verwaltet den Client-Lebenszyklus
- Handhabt die Kommunikation mit dem Server
- Implementiert lokale Datenverarbeitung

### WebUI
- Leichtgewichtige Web-Oberfläche für Monitoring
- Konfigurationsmöglichkeiten
- Statusanzeigen und Statistiken

## Entwicklung

### Voraussetzungen
- Go 1.20 oder höher
- libpcap-dev (für Paketerfassung)
- Root-Rechte für Netzwerkzugriff

### Build
```bash
cd cmd
go build -o ../bin/client
```

### Konfiguration
Die Konfiguration erfolgt über:
1. Kommandozeilenparameter
2. Konfigurationsdatei (configs/client.yaml)
3. Umgebungsvariablen

### Tests
```bash
go test ./...
``` 