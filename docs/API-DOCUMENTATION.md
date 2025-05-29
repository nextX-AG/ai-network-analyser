# AI Network Analyser - API Documentation

## Übersicht

Die API des AI Network Analysers ist über den Basis-Pfad `/api` erreichbar. Alle Endpunkte erwarten und liefern JSON-Daten, sofern nicht anders angegeben.

## Authentifizierung

Aktuell ist keine Authentifizierung implementiert. API-Key-Unterstützung ist für zukünftige Versionen geplant.

## Allgemeine Antwortstruktur

Alle API-Antworten folgen diesem Format:
```json
{
  "success": true|false,
  "message": "Optionale Nachricht",
  "data": { /* Optionale Daten */ },
  "error": "Optionale Fehlermeldung"
}
```

## Endpunkte

### System & Status

#### GET /api/health
Prüft den Systemstatus und liefert Versionsinformationen.

**Antwort:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": "1.0.0",
    "build_date": "2025-05-29",
    "commit": "abc123",
    "uptime": "1h2m3s",
    "components": {
      "server": "healthy",
      "storage": "healthy"
    }
  }
}
```

### Netzwerkschnittstellen

#### GET /api/interfaces
Listet alle verfügbaren Netzwerkschnittstellen.

**Antwort:**
```json
{
  "success": true,
  "data": [
    {
      "name": "eth0",
      "index": 1,
      "mac_address": "00:11:22:33:44:55",
      "ip_addresses": ["192.168.1.100"],
      "is_up": true,
      "is_loopback": false
    }
  ]
}
```

### Paketerfassung

#### POST /api/analyze
Analysiert eine hochgeladene PCAP-Datei.

**Request:** Multipart-Formular mit `pcap`-Datei

**Antwort:**
```json
{
  "success": true,
  "message": "PCAP-Datei erfolgreich analysiert",
  "data": {
    "total_packets": 1000,
    "gateway_packets": 500,
    "gateway_percentage": 50.0,
    "sample_packets": [/* ... */]
  }
}
```

#### POST /api/live/start
Startet die Live-Erfassung.

**Request:**
```json
{
  "interface": "eth0",
  "filter": "port 80"
}
```

#### POST /api/live/stop
Stoppt die aktive Live-Erfassung.

### Remote Clients

#### GET /api/agents
Listet alle registrierten Remote-Clients.

**Antwort:**
```json
{
  "success": true,
  "data": [
    {
      "name": "client-1",
      "url": "http://192.168.1.100:8080",
      "status": "online",
      "last_seen": "2025-05-29T11:40:50Z",
      "interfaces": ["eth0", "wlan0"],
      "active_interface": "eth0",
      "packets_captured": 1000,
      "version": "1.0.0",
      "os": "linux",
      "hostname": "client-1"
    }
  ]
}
```

#### POST /api/agents/register
Registriert einen neuen Remote-Client.

**Request:**
```json
{
  "name": "client-1",
  "url": "http://192.168.1.100:8080",
  "interfaces": ["eth0", "wlan0"],
  "interface_details": [/* ... */],
  "version": "1.0.0",
  "os": "linux",
  "hostname": "client-1"
}
```

#### POST /api/agents/unregister
Meldet einen Remote-Client ab.

**Request:**
```json
{
  "name": "client-1"
}
```

#### POST /api/agents/heartbeat
Aktualisiert den Status eines Remote-Clients.

**Request:**
```json
{
  "name": "client-1",
  "status": "capturing",
  "packets_captured": 1000,
  "interface": "eth0"
}
```

#### POST /api/agents/capture/start
Startet die Paketerfassung auf einem Remote-Client.

**Request:**
```json
{
  "name": "client-1",
  "interface": "eth0",
  "filter": "port 80"
}
```

#### POST /api/agents/capture/stop
Stoppt die Paketerfassung auf einem Remote-Client.

**Request:**
```json
{
  "name": "client-1"
}
```

#### POST /api/agents/set-interface
Setzt die aktive Netzwerkschnittstelle eines Remote-Clients.

**Request:**
```json
{
  "name": "client-1",
  "interface": "eth0"
}
```

### Gateway-Analyse

#### GET /api/gateways
Listet erkannte Gateways.

**Antwort:**
```json
{
  "success": true,
  "data": [
    {
      "ip": "192.168.1.1",
      "mac": "11:22:33:44:55:66",
      "is_active": true,
      "role": "default_gateway",
      "services": ["DHCP", "DNS"]
    }
  ]
}
```

#### GET /api/traffic/gateway
Liefert Gateway-Traffic-Statistiken.

**Antwort:**
```json
{
  "success": true,
  "data": {
    "total_packets": 1000,
    "gateway_packets": 650,
    "gateway_percentage": 65.0,
    "protocols": {
      "DNS": 150,
      "HTTP": 200,
      "DHCP": 50,
      "ARP": 100,
      "Other": 150
    }
  }
}
```

#### GET /api/events/gateway
Liefert Gateway-relevante Ereignisse.

**Antwort:**
```json
{
  "success": true,
  "data": [
    {
      "timestamp": "2025-05-29T11:35:50Z",
      "event_type": "dhcp",
      "description": "DHCP-Lease-Erneuerung",
      "severity": "info",
      "gateway_ip": "192.168.1.1",
      "client_ip": "192.168.1.100"
    }
  ]
}
```

### WebSocket-Verbindungen

#### WebSocket /api/ws
Endpunkt für Echtzeit-Updates.

**Nachrichten-Typen:**
- `packet`: Neue Paketinformationen
- `stats`: Statistische Updates
- `system`: Systemnachrichten

## Fehlerbehandlung

Alle Fehlerantworten folgen diesem Format:
```json
{
  "success": false,
  "error": "Beschreibung des Fehlers"
}
```

Häufige HTTP-Statuscodes:
- 200: Erfolgreiche Anfrage
- 400: Ungültige Anfrage
- 404: Ressource nicht gefunden
- 409: Konflikt (z.B. bereits laufende Capture)
- 500: Interner Serverfehler 