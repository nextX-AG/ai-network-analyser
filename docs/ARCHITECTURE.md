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
