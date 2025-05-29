# AI Network Analyser - Server

## Übersicht

Die Server-Komponente des AI Network Analysers bietet:
- Zentrale Verwaltung aller Clients
- KI-basierte Netzwerkanalyse
- React-basierte Hauptanwendung
- Sprachverarbeitung und Timeline-Funktionen

## Struktur

```
server/
├── cmd/                  # Hauptanwendung
├── configs/              # Server-Konfigurationen
├── internal/             # Private Pakete
│   ├── ai/              # KI-Integration
│   │   ├── analyzer/    # Analyse-Logik
│   │   ├── models/      # KI-Modelle
│   │   └── services/    # KI-Dienste
│   ├── api/             # API-Definitionen
│   │   ├── handlers/    # HTTP-Handler
│   │   ├── middleware/  # HTTP-Middleware
│   │   └── routes/      # Routendefinitionen
│   ├── config/          # Konfigurationslogik
│   ├── packet/          # Paketverarbeitung
│   ├── speech/          # Sprachverarbeitung
│   └── timeline/        # Zeitleisten-Management
└── ui/                  # React Frontend
    ├── public/          # Statische Dateien
    └── src/             # React-Quellcode
        ├── components/  # UI-Komponenten
        ├── features/    # Feature-Module
        ├── services/    # API-Dienste
        └── utils/       # Hilfsfunktionen
```

## Features

### KI-Integration
- Netzwerkanalyse mittels Machine Learning
- Anomalie-Erkennung
- Automatische Klassifizierung

### API
- RESTful API für Client-Kommunikation
- WebSocket für Echtzeitdaten
- Authentifizierung und Autorisierung

### Frontend
- Moderne React-Anwendung
- Feature-basierte Architektur
- Responsive Design

### Timeline
- Zeitleisten-basierte Visualisierung
- Event-Synchronisation
- Filterung und Gruppierung

## Entwicklung

### Voraussetzungen
- Go 1.20 oder höher
- Node.js 18 oder höher
- SQLite3

### Backend Build
```bash
cd cmd
go build -o ../bin/server
```

### Frontend Build
```bash
cd ui
npm install
npm run build
```

### Entwicklungsserver starten
```bash
# Backend
cd cmd
go run main.go

# Frontend (separates Terminal)
cd ui
npm run dev
```

### Tests
```bash
# Backend Tests
go test ./...

# Frontend Tests
cd ui
npm test
``` 