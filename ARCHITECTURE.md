# AI Network Analyser - Architektur und Projektstruktur

## Überblick

Der AI Network Analyser ist eine verteilte Anwendung zur Netzwerkanalyse, bestehend aus mehreren Hauptkomponenten:
- Agent: Für die Paketerfassung und lokale Analyse
- Server: Für die zentrale Verwaltung und Visualisierung
- AI-Module: Für die intelligente Analyse von Netzwerkdaten
- Speech-Module: Für die Spracherkennung und -verarbeitung

## Projektstruktur

```
ai-network-analyser/
├── agent/                    # Remote-Agent Komponente
│   ├── cmd/                  # Hauptanwendung
│   ├── configs/              # Agent-Konfigurationen
│   ├── internal/             # Private Pakete
│   │   ├── capture/         # Paketerfassung
│   │   ├── service/         # Agent-Dienste
│   │   └── webui/          # Agent-UI Backend
│   └── ui/                  # Frontend
│       ├── public/          # Statische Dateien
│       ├── templates/       # HTML-Templates
│       └── static/          # Assets (CSS, JS, Images)
│
├── server/                   # Server Komponente
│   ├── cmd/                 # Hauptanwendung
│   ├── configs/             # Server-Konfigurationen
│   ├── internal/            # Private Pakete
│   │   ├── ai/             # KI-Integrationen
│   │   │   ├── analyzer/   # KI-basierte Analyse
│   │   │   ├── models/     # KI-Modelle
│   │   │   └── services/   # KI-Dienste
│   │   ├── api/            # API-Definitionen
│   │   │   ├── handlers/   # HTTP-Handler
│   │   │   ├── middleware/ # HTTP-Middleware
│   │   │   └── routes/     # Routendefinitionen
│   │   ├── config/         # Konfigurationslogik
│   │   ├── packet/         # Paketverarbeitung
│   │   │   ├── analyzer/   # Paketanalyse-Logik
│   │   │   └── processor/  # Paketverarbeitung
│   │   ├── speech/         # Sprachverarbeitung
│   │   │   ├── recognition/# Spracherkennung (Whisper)
│   │   │   └── synthesis/  # Sprachsynthese
│   │   └── timeline/       # Zeitleisten-Management
│   │       ├── events/     # Event-Handling
│   │       └── sync/       # Zeitsynchronisation
│   ├── ui/                 # React Frontend
│   │   ├── public/         # Statische Dateien
│   │   └── src/            # React-Quellcode
│   │       ├── components/ # UI-Komponenten
│   │       │   ├── ai/     # KI-bezogene Komponenten
│   │       │   ├── common/ # Gemeinsame Komponenten
│   │       │   ├── filters/# Filter-Komponenten
│   │       │   └── layout/ # Layout-Komponenten
│   │       ├── features/   # Feature-Module
│   │       │   ├── ai/     # KI-Features
│   │       │   ├── capture/# Netzwerkerfassung
│   │       │   ├── agents/ # Remote-Agents
│   │       │   ├── speech/ # Sprachverarbeitung
│   │       │   └── timeline/# Zeitleiste
│   │       ├── services/   # API-Dienste
│   │       └── utils/      # Hilfsfunktionen
│
├── pkg/                      # Gemeinsame Pakete
│   ├── ai/                  # KI-Utilities
│   ├── common/              # Geteilte Funktionen
│   ├── models/              # Datenmodelle
│   ├── protocol/            # Protokolldefinitionen
│   └── version/             # Versionierung
│
├── docs/                     # Dokumentation
├── scripts/                  # Build-Skripte
└── configs/                  # Globale Konfigurationen

```

## Architekturprinzipien

1. **Klare Komponentenaufteilung**
   - Strikte Trennung zwischen Agent und Server
   - Modulare KI- und Sprachverarbeitungskomponenten
   - Jede Komponente ist eigenständig lauffähig

2. **Konsistente UI-Struktur**
   - Agent: Leichtgewichtige Template-basierte UI
   - Server: React-basierte Hauptanwendung mit Feature-Modulen
   - Integrierte KI- und Sprachverarbeitungs-UI

3. **Modulare Codeorganisation**
   - Gemeinsamer Code in `pkg/`
   - Komponentenspezifischer Code in jeweiligen Verzeichnissen
   - Spezialisierte Module für KI und Sprachverarbeitung

4. **KI-Integration**
   - Lokale LLM-Integration für Offline-Fähigkeit
   - OpenAI API-Integration für erweiterte Funktionen
   - Modulares KI-System für einfache Erweiterbarkeit

5. **Sprachverarbeitung**
   - Whisper.cpp Integration für lokale Spracherkennung
   - Optionale Cloud-API-Anbindung
   - Asynchrone Verarbeitung für bessere Performance

6. **Konfigurationsmanagement**
   - Globale Konfigurationen in `configs/`
   - Komponentenspezifische Konfigurationen lokal
   - Flexible API-Schlüssel-Verwaltung

## Entwicklungsrichtlinien

1. **Code-Organisation**
   - Neue Features in entsprechende Komponente einordnen
   - Gemeinsamen Code in `pkg/` auslagern
   - UI-Komponenten in jeweiligen `ui/` Verzeichnissen
   - KI- und Sprachmodule klar strukturieren

2. **Dependency Management**
   - Externe Abhängigkeiten in `go.mod` verwalten
   - Frontend-Dependencies in jeweiligen `package.json`
   - KI- und Sprachmodell-Versionen dokumentieren

3. **Dokumentation**
   - Technische Dokumentation in `docs/`
   - README.md für Komponenten
   - Inline-Dokumentation für Funktionen
   - API-Dokumentation für KI- und Sprachschnittstellen

## Build und Deployment

1. **Build-Prozess**
   - Separate Builds für Agent und Server
   - Frontend-Builds in jeweiligen UI-Verzeichnissen
   - KI-Modell-Management und Versionierung
   - Sprachmodell-Integration

2. **Deployment**
   - Containerisierte Deployments
   - Konfiguration über Umgebungsvariablen
   - Skalierbare KI- und Sprachverarbeitung

## Entwicklungsworkflow

1. **Feature-Entwicklung**
   - Branch von `main` erstellen
   - Feature implementieren
   - Tests schreiben
   - Pull Request erstellen

2. **Code Review**
   - Automatische Tests
   - Peer Review
   - Merge nach Bestätigung

## Nächste Schritte

1. **Restrukturierung**
   - [ ] Backup erstellen
   - [ ] Neue Verzeichnisstruktur anlegen
   - [ ] Code migrieren
   - [ ] Tests anpassen
   - [ ] Dokumentation aktualisieren

2. **KI-Integration**
   - [ ] LLM-Integration implementieren
   - [ ] KI-Analysemodule entwickeln
   - [ ] UI-Komponenten erstellen
   - [ ] API-Schnittstellen definieren

3. **Sprachverarbeitung**
   - [ ] Whisper.cpp einbinden
   - [ ] Spracherkennungsmodule implementieren
   - [ ] UI-Integration entwickeln
   - [ ] API-Endpunkte erstellen

4. **Cleanup**
   - [ ] Alte Verzeichnisse entfernen
   - [ ] Imports aktualisieren
   - [ ] Dependencies bereinigen

5. **Validierung**
   - [ ] Builds testen
   - [ ] Integration testen
   - [ ] UI-Tests durchführen
   - [ ] KI- und Sprachmodule validieren 