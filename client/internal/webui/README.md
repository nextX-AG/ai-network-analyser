# Client WebUI Module

Dieses Modul implementiert die leichtgewichtige Web-Benutzeroberfläche für den Network Analyzer Client.

## Struktur

```
webui/
├── webui.go          # Hauptimplementierung der Web-UI
└── README.md         # Diese Datei
```

## Features

- Administrationsschnittstelle für Client-Konfiguration
- Echtzeit-Statusanzeige
- Netzwerkschnittstellen-Management
- Konfigurationsverwaltung

## Technologien

- Go Templates für HTML-Rendering
- Statisches CSS/JS für UI-Funktionalität
- WebSocket für Echtzeit-Updates
- RESTful API für Konfiguration

## Verzeichnisstruktur

Die UI-Assets befinden sich in `client/ui/`:
- `templates/` - HTML-Templates
- `static/` - CSS, JavaScript, Bilder
- `public/` - Öffentlich zugängliche Dateien
