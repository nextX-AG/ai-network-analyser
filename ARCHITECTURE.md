# KI-Netzwerk-Analyzer Architecture

This document provides a high-level overview of the KI-Netzwerk-Analyzer architecture. For more detailed documentation, please see the [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Overview

The KI-Netzwerk-Analyzer is a modular system for intelligent network traffic analysis with AI integration and a focus on gateway analysis. The application combines network packet knowledge with modern AI technology and allows users to mark network events and annotate them with metadata, either manually or through AI assistance.

## System Architecture

```
[Network Traffic (PCAP/tcpdump/live)]
      |
      v
 [Go Backend]
   - Packet Capture (gopacket)
   - Speech2Text (Whisper.cpp local, API optional)
   - AI Annotation (OpenAI GPT API, later local LLM)
   - Event/Timeline Module
   - REST/Websocket API
      |
      v
[Frontend: Three.js/React]
   - Timeline & Graph Visualization
   - Event Annotation (Text/Speech)
   - AI Results/Info Panels
   - Logbook/Export
```

## Main Components

### Core Modules
- **Packet Capture (PCAP/TCPDUMP)**: Reading and analyzing PCAP files, live capture of network packets
- **Gateway Analysis**: Detection and analysis of gateway-relevant network traffic
- **Real-time Monitoring**: Live capturing and analyzing network packets via WebSockets

### Remote Capture System
- **Capture Agent**: Lightweight service for remote devices (e.g., Raspberry Pi) for network capture
- **Agent API**: REST interface for configuration and control of remote capture devices
- **Stream Protocol**: Efficient WebSocket-based protocol for transmitting packet data
- **Multi-Agent Management**: Management of multiple capture devices at different network points

### Advanced Features
- **Speech2Text**: Transcription of voice notes for event markers
- **AI Annotation**: Automatic analysis and annotation of network packets with LLMs
- **Timeline Visualization**: Graphical representation of network events on a timeline
- **Event Marking**: Manual and automatic marking of interesting network events

## Technology Stack

- **Backend**: Go with the following main libraries:
  - `gopacket`: Packet capture and analysis
  - `gorilla/mux`, `gorilla/websocket`: HTTP routing and WebSockets
  - `gorm` or `sqlx`: Database interaction
  
- **Frontend**:
  - React/TypeScript
  - Three.js for complex visualizations
  - WebSockets for real-time updates

- **AI/ML**:
  - Optional integration with OpenAI API for advanced analysis
  - Local Whisper integration for speech recognition

- **Database**:
  - SQLite for simple deployment and single-user mode

## Project Structure

```
ai-network-analyser/
├── bin/                  # Binaries and executables
├── cmd/                  # Main applications
│   ├── agent/            # Remote capture agent implementation
│   └── server/           # Web server implementation
├── configs/              # Configuration files
├── data/                 # Data storage
├── docs/                 # Documentation
├── internal/             # Internal packages
│   ├── agent/            # Agent-specific code
│   ├── ai/               # AI integration
│   ├── api/              # API handlers
│   ├── config/           # Configuration structures
│   ├── packet/           # Packet analysis
│   ├── speech/           # Speech recognition
│   ├── storage/          # Data storage
│   └── timeline/         # Timeline and event management
├── pkg/                  # Reusable packages
│   ├── models/           # Data models
│   ├── protocol/         # Protocol definitions
│   └── utils/            # Utility functions
├── web/                  # Web frontend
    ├── public/           # Static files
    └── src/              # React source code
```

## Future Extensions

- Integration of local LLMs as an alternative to OpenAI
- Advanced protocol analysis for specific application protocols
- Multi-agent capture with synchronized timeline
- Agent group management and organization
- Collaborative analysis with multiple users
- Cloud-based deployment option
- Mobile view for frontend
- Integration with external security tools
- Full implementation of AQEA object model 

## Frontend-Architektur

Die Webanwendung folgt einer strukturierten, feature-basierten Architektur mit React, die klare Trennung von Zuständigkeiten und eine modulare Bauweise ermöglicht.

### Verzeichnisstruktur

```
web/
│
├── public/              # Statische Assets und HTML-Template
│
└── src/                 # Quellcode der Webanwendung
    ├── assets/          # Bilder, Fonts und andere statische Ressourcen
    │
    ├── features/        # Feature-Module (Feature-basierte Organisation)
    │   ├── remoteAgents/  # Remote-Agents Funktionalität
    │   │   ├── components/  # Remote-Agent-spezifische Komponenten
    │   │   │   ├── AgentCard.jsx        # Einzelne Agent-Karte
    │   │   │   ├── AgentFilter.jsx      # Filterkomponente für Agenten
    │   │   │   └── AgentStatusDisplay.jsx # Statusanzeige eines Agenten
    │   │   ├── containers/
    │   │   │   └── RemoteAgentsContainer.jsx # Container für die Agent-Verwaltung
    │   │   ├── hooks/
    │   │   │   └── useAgentControl.js   # Custom Hook für Agent-Steuerung
    │   │   ├── services/
    │   │   │   └── agentApi.js          # API-Dienste für Agenten
    │   │   └── RemoteAgentsPage.jsx     # Hauptseite für Remote-Agents
    │   │
    │   ├── networkCapture/  # Netzwerk-Capture Funktionalität
    │   │   ├── components/
    │   │   │   ├── CapturePanel.jsx     # Panel für Capture-Steuerung
    │   │   │   └── CaptureStatus.jsx    # Statusanzeige für aktive Captures
    │   │   ├── services/
    │   │   │   └── captureApi.js        # API-Dienste für Capture-Funktionen
    │   │   └── NetworkCapturePage.jsx   # Hauptseite für Netzwerk-Capture
    │   │
    │   ├── timeline/     # Timeline-Funktionalität
    │   │   ├── components/
    │   │   │   ├── TimelineControls.jsx # Steuerelemente für Timeline
    │   │   │   └── TimelineView.jsx     # Visualisierung der Timeline
    │   │   ├── services/
    │   │   │   └── timelineApi.js       # API-Dienste für Timeline-Daten
    │   │   └── TimelinePage.jsx         # Hauptseite für Timeline
    │   │
    │   └── ai/           # KI-Funktionalität
    │       ├── components/
    │       │   ├── AIAnalysisPanel.jsx  # Panel für KI-Analyse
    │       │   └── AISettings.jsx       # Einstellungen für KI-Funktionen
    │       ├── services/
    │       │   └── aiApi.js             # API-Dienste für KI-Funktionen
    │       └── AIAssistantPage.jsx      # Hauptseite für KI-Assistent
    │
    ├── shared/           # Gemeinsam genutzte Komponenten und Utilities
    │   ├── components/   # Wiederverwendbare UI-Komponenten
    │   │   ├── Button.jsx
    │   │   ├── Card.jsx
    │   │   ├── Dialog.jsx
    │   │   └── ...
    │   │
    │   ├── hooks/        # Gemeinsam genutzte Custom Hooks
    │   │   ├── useApi.js
    │   │   ├── useWebSocket.js
    │   │   └── ...
    │   │
    │   └── utils/        # Hilfsfunktionen und gemeinsam genutzte Logik
    │       ├── dateUtils.js
    │       ├── formatters.js
    │       └── ...
    │
    ├── context/          # Globaler Zustandsmanagement
    │   ├── AuthContext.js
    │   └── AppContext.js
    │
    ├── types/            # TypeScript-Typendefinitionen
    │
    ├── App.jsx           # Hauptanwendungskomponente
    └── index.jsx         # Einstiegspunkt der Anwendung
```

### Vorteile der Feature-basierten Organisation

Die feature-basierte Struktur bietet mehrere Vorteile:

1. **Kohäsion:** Zusammengehörige Komponenten sind nach Funktionalität organisiert, nicht nach Komponentenart
2. **Wartbarkeit:** Einfachere Navigation und Verständnis des Codes, da alle zu einem Feature gehörenden Dateien in einem Verzeichnis liegen
3. **Skalierbarkeit:** Neue Features können hinzugefügt werden, ohne die Struktur bestehender Features zu beeinflussen
4. **Teamarbeit:** Mehrere Entwickler können parallel an verschiedenen Features arbeiten mit minimalen Konflikten
5. **Wiederverwendbarkeit:** Gemeinsam genutzte Komponenten sind im `shared`-Verzeichnis organisiert

### Vergleich mit vorheriger Struktur

Die vorherige Struktur war nach Komponententyp (components, pages, services) organisiert, was bei wachsender Anwendung zu Unübersichtlichkeit führt. Durch die Umstellung auf eine feature-basierte Organisation wird die Codebase auch bei Erweiterung übersichtlich und modularer bleiben.

### Migrationsansatz

Die Migration von der bestehenden Struktur zur feature-basierten Organisation sollte schrittweise erfolgen:

1. Erstellung der neuen Verzeichnisstruktur
2. Migration eines Features nach dem anderen (z.B. zuerst Remote-Agents)
3. Anpassung der Imports und Routing
4. Refactoring von gemeinsam genutzten Komponenten in das `shared`-Verzeichnis
5. Tests zur Sicherstellung der Funktionalität

### Best Practices

- **Feature-Grenzen:** Klare Definition, was zu einem Feature gehört und was gemeinsam genutzt wird
- **Zirkuläre Abhängigkeiten vermeiden:** Features sollten nicht direkt voneinander abhängen, sondern über `shared` oder globalen State kommunizieren
- **Lazy-Loading:** Features können bei Bedarf dynamisch geladen werden für verbesserte Performance
- **Testing:** Jedes Feature sollte eigene Tests haben, die in seinem Verzeichnis liegen 