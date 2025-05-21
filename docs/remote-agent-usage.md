# Verwendung des Remote-Agent-Systems

Diese Anleitung beschreibt, wie Sie mit dem Hauptserver Remote-Agents verwalten und damit Netzwerkverkehr an verschiedenen Stellen im Netzwerk erfassen können.

## Voraussetzungen

- Laufender Hauptserver (KI-Netzwerk-Analyzer)
- Ein oder mehrere installierte Remote-Agents (siehe [Agent-Installationsanleitung](agent-installation.md))

## Hauptserver konfigurieren

### 1. Server für Remote-Agents konfigurieren

Stellen Sie sicher, dass der Hauptserver auf einer Adresse lauscht, die von den Remote-Agents erreicht werden kann:

```json
{
  "server": {
    "host": "0.0.0.0",  // Wichtig: Nicht 127.0.0.1 verwenden
    "port": 9090,
    "enable_websocket": true
  }
}
```

### 2. Hauptserver starten

```bash
./bin/analyzer --config=configs/config.json
```

Erfolgreiche Ausgabe:
```
Server gestartet auf 0.0.0.0:9090
```

## Remote-Agents verwalten

### Registrierte Agents anzeigen

Öffnen Sie im Webbrowser die Adresse des Hauptservers: `http://server-ip:9090`

Navigieren Sie zum Bereich "Remote Agents". Hier werden alle registrierten Agents mit Status, Name, IP-Adresse und Interface angezeigt.

### Agent registrieren

Die Registrierung erfolgt automatisch, wenn der Agent gestartet wird und die korrekte Server-URL konfiguriert hat.

Alternativ können Sie folgende Schritte durchführen:
1. Greifen Sie auf das Webinterface des Agents zu: `http://agent-ip:8090/admin`
2. Konfigurieren Sie die korrekte Server-URL
3. Klicken Sie auf "Bei Server registrieren"

### Agent-Status prüfen

Im Hauptserver:
1. Navigieren Sie zum Bereich "Remote Agents"
2. Prüfen Sie den Status jedes Agents (aktiv, inaktiv, erfassend)
3. Sehen Sie Details wie erfasste Pakete, aktive Schnittstelle und Netzwerkkonfiguration

## Netzwerkerfassung mit Remote-Agents

### Erfassung starten

1. Navigieren Sie im Hauptserver zum Bereich "Remote Agents"
2. Wählen Sie den gewünschten Agent aus
3. Klicken Sie auf "Erfassung starten"
4. Wählen Sie ggf. eine spezifische Schnittstelle und Filter aus

Hinweis: Bei Bridge-konfiguriertem Agent sollten Sie "br0" als Interface wählen.

### Erfassung in Echtzeit überwachen

Während der Erfassung werden die Pakete in Echtzeit über WebSocket an den Hauptserver übertragen:

1. Navigieren Sie zum Bereich "Live-Capture"
2. Wählen Sie im Dropdown-Menü den aktiven Remote-Agent aus
3. Die erfassten Pakete werden in Echtzeit angezeigt
4. Nutzen Sie die Filter, um bestimmte Pakettypen hervorzuheben

### Erfassung stoppen

1. Navigieren Sie im Hauptserver zum Bereich "Remote Agents"
2. Wählen Sie den aktiven Agent aus
3. Klicken Sie auf "Erfassung stoppen"

## Datenanalyse

Die vom Remote-Agent erfassten Daten werden genau wie lokal erfasste Daten analysiert:

1. Zeitbasierten Traffic auf der Timeline anzeigen
2. Gateway-bezogenen Verkehr hervorheben
3. DHCP-, DNS- und ARP-Verkehr analysieren
4. Netzwerkgraphen und Verbindungen visualisieren

### Gateway-Analyse

Remote-Agents sind besonders wertvoll für Gateway-Analysen:

1. Platzieren Sie einen Agent mit Bridge-Konfiguration zwischen Gateway und Netzwerk
2. Navigieren Sie zum Bereich "Gateway-Analyse"
3. Wählen Sie den Agent, der am Gateway platziert ist
4. Analysieren Sie die Gateway-Kommunikation und -Funktionen

## Multi-Agent-Erfassung

Für komplexe Netzwerke können Sie mehrere Agents gleichzeitig einsetzen:

1. Platzieren Sie Agents an strategischen Punkten im Netzwerk
2. Starten Sie die Erfassung auf mehreren Agents
3. Wechseln Sie zwischen den Agent-Streams oder nutzen Sie die aggregierte Ansicht
4. Korrelieren Sie Ereignisse über verschiedene Netzwerksegmente hinweg

## Fehlerbehebung

### Agent erscheint nicht in der Liste

1. Prüfen Sie, ob der Agent läuft: `sudo systemctl status ki-network-analyzer-agent`
2. Überprüfen Sie die Server-URL in der Agent-Konfiguration
3. Stellen Sie sicher, dass keine Firewall die Verbindung blockiert
4. Prüfen Sie die Netzwerkkonnektivität: `ping server-ip`

### Keine Pakete werden empfangen

1. Überprüfen Sie, ob der Agent aktiv erfasst
2. Prüfen Sie die Schnittstelle und den Promisc-Modus
3. Überprüfen Sie die WebSocket-Verbindung im Browser-Entwicklertool
4. Testen Sie die Verbindung mit einem einfachen Ping oder Traceroute

### Verbindungsverlust während der Erfassung

1. Überprüfen Sie die Netzwerkstabilität
2. Prüfen Sie die Server- und Agent-Logs
3. Starten Sie den Agent bei Bedarf neu: `sudo systemctl restart ki-network-analyzer-agent` 