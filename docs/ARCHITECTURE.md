# KI-Netzwerk-Analyzer – Initial ARCHITECTURE.md

**Stand: 2025-05-20 (Initial Canvas / living document)**

---

## Executive Summary

Dies ist das initiale Architektur-Canvas für ein modernes, modular aufgebautes Framework zur intelligenten, annotierbaren Analyse von Netzwerkverkehr (PCAP/TCPDUMP) mit KI-Integration und Timeline-Event-Markierung. Fokus liegt auf Performance, Erweiterbarkeit und maximaler Entwicklerfreundlichkeit. Die Anbindung von AQEA-Objektmodell ist für spätere Versionen vorgesehen und wird durch Modularität sichergestellt.

---

## Architekturübersicht (Top Level)

```
[Netzwerk-Traffic (PCAP/tcpdump/live)]
      |
      v
 [Go-Backend]
   - Packet-Capture (gopacket)
   - Speech2Text (Whisper.cpp lokal, API optional)
   - KI-Annotation (OpenAI GPT-API, später lokale LLM)
   - Event/Timeline-Modul
   - REST/Websocket-API
      |
      v
[Frontend: Three.js/React]
   - Timeline & Graph-Visualisierung
   - Event-Annotation (Text/Sprache)
   - KI-Resultate/Info-Panels
   - Logbuch/Export
```

---

## Modularitäts- und Entwicklungsprinzipien

* **Maximale Erweiterbarkeit:** Jedes Modul (Capture, Speech2Text, KI, Event-Store, Frontend) unabhängig austausch-/erweiterbar
* **Schnelle Iteration:** MVP-First, dann gezielte Erweiterung
* **Technologiefokus:** Go für Backend/Parsing/Services, Three.js in React für Timeline- und Datenvisualisierung (Frontend)
* **Keine CORS-Probleme:** API & Frontend laufen lokal auf einer Instanz bzw. über Reverse-Proxy; keine komplizierten Multi-Origin-Setups
* **KI & Speech Modular:** Speech2Text über Whisper.cpp (lokal), KI-Analyse initial per Cloud-API (OpenAI), später Umstieg auf lokale LLMs
* **Späteres AQEA-Modul:** Alle Events/Objekte werden so gestaltet, dass sie später AQEA-fähig erweitert werden können (JSON/Protobuf/MsgPack)

---

## Komponenten im Detail

### 1. Go-Backend

* **gopacket-basiertes Packet-Capturing** (live & pcap-Import)
* **Speech2Text-Modul**

  * Whisper.cpp (lokal via CLI/Binary-Call, ggf. Modul-Schnittstelle für Alternativen)
  * API-Fallback: OpenAI Whisper oder andere Speech2Text-Anbieter
* **KI-Modul**

  * Cloudbasierte LLM (OpenAI GPT-API), später optionale lokale LLMs (Ollama/LM Studio)
  * Aufgaben: Eventbeschreibung, Protokollanalyse, Extraktion, Code-Snippets, Dokumentation
* **Event & Timeline-Modul**

  * Erzeugung, Speicherung & Management von Event-Markern (Text/Sprache/Timestamp)
  * Markierung, Notiz, Typisierung
  * Export als JSON/Markdown/CSV
* **REST/Websocket-API** (JSON)

  * Endpunkte: `/capture`, `/events`, `/speechmark`, `/ai/annotate`, `/export` etc.
* **Persistenz/DB**

  * Initial: SQLite (Dateibasiert, schnell, einfach)
  * Events, Sessions, User-Notizen, KI-Resultate

### 2. Frontend

* **Three.js + React/Next.js** (modernes SPA)
* **Timeline/Graph-Komponente**

  * Anzeige von Netzwerk-Events, KI-Annotations, User-Markern
  * Filter, Suche, Zoom, Notiz-Felder
* **Event-Annotation**

  * Markierung (Button, Hotkey, Sprache)
  * Sofortige Rückmeldung & Marker auf der Timeline
* **KI-Analyse-Anzeige**

  * Infopanels, automatisierte Notizen, Code-Vorschläge
* **Session/Logbuch/Export**

  * Markdown/CSV/JSON-Export
* **API-Anbindung**

  * REST/WS-Integration, keine CORS-Probleme durch gemeinsames Hosting/Proxy

---

## Schnittstellen und Kommunikation

* **Backend bietet REST-API (JSON) für Frontend**
* **Websockets** optional für Live-Events/Push (z.B. Live Capture, Marker)
* **Frontend nutzt ausschließlich diese API, kein direkter DB-Zugriff**
* **Speech2Text und KI können sowohl asynchron als auch synchron angesprochen werden (Batch & Live-Modus)**

---

## KI/Speech-Modul (initial)

* **Speech2Text:** Whisper.cpp lokal (CLI), alternativ API-basierte Engine (frei konfigurierbar, austauschbar)
* **KI/LLM:** OpenAI GPT-API für Event-Zuordnung, Protokoll-Analyse, automatisierte Notizen (Modul austauschbar, später lokale LLM)

---

## Persistenz

* **SQLite** als initiale Datenbank (einfach zu deployen, keine Setup-Hürden, performant für Einzel/MVP-Systeme)
* **Tabellen:** Events, Sessions, KI-Resultate, Usernotes
* **Migration auf MongoDB, Postgres etc. bei Bedarf später problemlos möglich (ORM)**

---

## Erweiterbarkeit & Integration

* **Spätere AQEA-Anbindung:** Alle Events/Objekte so modelliert, dass sie AQEA-konform erweitert werden können
* **Plugins:** Capture- und Event-Module als unabhängige Plugins denkbar (z.B. Protokoll-Decoder, Exportformate)
* **Backend kann später Microservices-Architektur erhalten (EventStore, Speech, KI als einzelne Services)**
* \*\*Frontend kann als Standalone-Spa oder als Modul in größere Plattformen eingebettet werden

---

## Initiale Roadmap

1. Go-Backend Scaffold (Capture, API, SQLite)
2. Three.js/React Frontend Scaffold (Timeline, Marker, API-Bindung)
3. Speech2Text-Integration (Whisper.cpp, REST-Endpoint)
4. KI-Integration (OpenAI GPT-API, Analyse-Endpunkt)
5. Event-Timeline mit Marker, Notizen, Export
6. Refactoring/Modularisierung für spätere AQEA/Plugin/Cloud-Erweiterung

---

## Anmerkungen

* MVP-orientiert starten, Features dann gezielt und modular erweitern
* Alles wird so dokumentiert, dass später Multi-User/Cloud und AQEA möglich sind
* **Kein Vendor-Lock-In, volle Kontrolle über Daten und Erweiterungen**

## Überblick
Der KI-Netzwerk-Analyzer ist ein modulares System zur Analyse von Netzwerkverkehr, das mit KI-Integration und Schwerpunkt auf Gateway-Analyse arbeitet. Die Anwendung kombiniert Netzwerkpaketkenntnisse mit moderner KI-Technologie und ermöglicht es Benutzern, Netzwerkereignisse zu markieren und mit Metadaten zu versehen, sei es manuell oder durch KI-Unterstützung.

## Hauptmodule und Komponenten

### Core Module
- **Paketerfassung (PCAP/TCPDUMP)**: Lesen und Analysieren von PCAP-Dateien, Live-Erfassung von Netzwerkpaketen
- **Gateway-Analyse**: Erkennung und Analyse von Gateway-relevantem Netzwerkverkehr
- **Echtzeit-Monitoring**: Live-Erfassung und Analyse von Netzwerkpaketen über WebSockets

### Remote-Capture-System
- **Capture-Agent**: Leichtgewichtiger Dienst für Remote-Geräte (z.B. Raspberry Pi) zur Netzwerkerfassung
- **Agent-API**: REST-Schnittstelle zur Konfiguration und Steuerung von Remote-Erfassungsgeräten
- **Stream-Protokoll**: Effizientes WebSocket-basiertes Protokoll zum Übertragen von Paketdaten
- **Multi-Agent-Management**: Verwaltung mehrerer Erfassungsgeräte an verschiedenen Netzwerkpunkten

### Erweiterte Funktionen
- **Speech2Text**: Transkription von Sprachnotizen für Ereignismarkierungen
- **KI-Annotation**: Automatische Analyse und Annotation von Netzwerkpaketen mit LLMs
- **Timeline-Visualisierung**: Graphische Darstellung von Netzwerkereignissen auf einer Zeitleiste
- **Ereignismarkierung**: Manuelle und automatische Markierung interessanter Netzwerkereignisse

## Architektonische Muster

### Backend (Go)
- **Modulare Struktur**: Klare Trennung von Zuständigkeiten durch strukturierte Pakete
- **Schichtarchitektur**: Handler -> Service -> Repository-Pattern
- **Dependency Injection**: Flexible Konfiguration und bessere Testbarkeit
- **Pub/Sub-Muster**: Für Echtzeit-Updates und Ereignisverarbeitung

### Frontend (React)
- **Komponenten-basiert**: Wiederverwendbare UI-Komponenten
- **Redux/Context**: Zentrales State-Management
- **WebSocket-Integration**: Echtzeit-Updates der Benutzeroberfläche

## Remote-Capture-Architektur

Die Remote-Capture-Funktionalität basiert auf einem Agent-Server-Modell:

1. **Capture-Agent (auf Remote-Gerät)**
   - Leichtgewichtiger Go-Dienst, der auf Edge-Geräten wie Raspberry Pi oder UP Board läuft
   - Direkter Zugriff auf Netzwerkschnittstellen über libpcap/gopacket
   - Integriertes Webinterface zur Konfiguration und Verwaltung
   - Automatische Erkennung und Registrierung beim Hauptserver
   - Spezialisierte Bridge-Unterstützung für MITM-Monitoring
   - REST-API für Konfiguration und Verwaltung
   - WebSocket-Endpunkt für Paket-Streaming

2. **Hauptanwendung (Server)**
   - Verwaltet Verbindungen zu mehreren Remote-Agents
   - Aggregiert und verarbeitet Daten von allen Erfassungspunkten
   - Bietet einheitliche UI für die Verwaltung aller Erfassungsgeräte

3. **Kommunikationsprotokoll**
   - REST-API für Konfiguration, Start/Stop und Status
   - WebSocket für effizientes Echtzeit-Streaming von Paketdaten
   - Optimierte Datenübertragung (Serialisierung, Kompression)
   - Authentifizierung über API-Keys

4. **Deployment-Optionen**
   - Standalone-Binary für Edge-Geräte (Go's Cross-Compilation)
   - Systemd-Service für automatischen Start und Überwachung
   - Konfigurierbare Bridge-Schnittstellen für MITM-Monitoring

### Agent Web-Interface

Der Remote-Agent verfügt über ein eingebautes Web-Interface, das folgende Funktionen bietet:
- Status-Übersicht (Paketzähler, Verbindungsstatus, Schnittstelleninformationen)
- Konfiguration (Server-URL, Schnittstellenauswahl, API-Key)
- Netzwerkschnittstellen-Übersicht mit Erkennung von Bridge-Interfaces
- Aktions-Buttons (Agent-Neustart, Server-Registrierung)

### Automatische Erkennung und Registrierung

Der Agent verfügt über eine automatische Registrierungsfunktion:
- Beim Start versucht der Agent, sich automatisch beim konfigurierten Server zu registrieren
- Erfasst und übermittelt Informationen zu allen verfügbaren Netzwerkschnittstellen
- Fallback auf manuelle Registrierung via Web-Interface

### Bridge-Optimierung für MITM-Monitoring

Spezielle Unterstützung für Bridge-Schnittstellen:
- Automatische Erkennung von Bridge-Interfaces im Betriebssystem
- Optimierte Paketerfassung mit erhöhten Buffer-Größen für Bridge-Traffic
- Konfiguration des Promisc-Modus und Immediate-Mode für bessere Leistung
- Dokumentierte Anleitung zur Einrichtung von Netzwerk-Bridges für effektives MITM-Monitoring

Diese Architektur ermöglicht ein skalierbares Netzwerk von Erfassungspunkten, die strategisch in einer Infrastruktur platziert werden können, während die zentrale Anwendung alle Daten aggregiert und analysiert.

## Datenbankstruktur

Die Datenbank besteht aus den folgenden Haupttabellen:
- **Packets**: Gespeicherte Paketinformationen für forensische Analyse
- **Events**: Markierte Netzwerkereignisse mit Metadaten
- **Annotations**: Benutzer- und KI-generierte Anmerkungen zu Ereignissen
- **GatewayInfo**: Informationen über identifizierte Gateway-Geräte

## API-Struktur

Die API ist RESTful mit den folgenden Hauptendpunkten:

- `/api/analyze`: PCAP-Dateianalyse
- `/api/live/start`, `/api/live/stop`: Steuerung der Live-Erfassung
- `/api/interfaces`: Verfügbare Netzwerkschnittstellen
- `/api/gateways`: Gateway-Informationen
- `/api/events/gateway`: Gateway-spezifische Ereignisse
- `/api/remote`: Verwaltung von Remote-Capture-Agents
- `/api/ws`: WebSocket-Endpunkt für Echtzeit-Updates

## Technologie-Stack

- **Backend**: Go mit folgenden Hauptbibliotheken:
  - `gopacket`: Paketerfassung und -analyse
  - `gorilla/mux`, `gorilla/websocket`: HTTP-Routing und WebSockets
  - `gorm` oder `sqlx`: Datenbankinteraktion
  
- **Frontend**:
  - React/TypeScript
  - Three.js für komplexe Visualisierungen
  - WebSockets für Echtzeit-Updates

- **KI/ML**:
  - Optional-Integration mit OpenAI API für erweiterte Analyse
  - Lokale Whisper-Integration für Spracherkennung

- **Datenbank**:
  - SQLite für einfache Bereitstellung und Einzelbenutzer-Modus

## Sicherheitskonzept

- Sichere Speicherung von API-Schlüsseln
- Zugriffssteuerung für sensible Operationen
- Validierung aller Eingaben
- Sichere WebSocket-Kommunikation
