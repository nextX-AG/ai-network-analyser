# Installation des KI-Netzwerk-Analyzer Remote Agents

Diese Anleitung beschreibt die Installation und Konfiguration des Remote-Agents auf einem UP Board oder einem anderen Linux-System.

## Systemanforderungen

- Linux-System (Ubuntu/Debian empfohlen)
- Go 1.16+ (für Kompilierung aus Quellcode)
- Administratorrechte (sudo)
- git

## Installation aus Quellcode

Da aktuell noch keine vorgefertigten Binaries als Release verfügbar sind, muss der Agent aus dem Quellcode kompiliert werden.

### 1. Go installieren (falls noch nicht geschehen)

```bash
sudo apt-get update
sudo apt-get install -y golang-go
```

Überprüfen Sie die Go-Version:

```bash
go version
```

Die Version sollte 1.16 oder höher sein.

### 2. Repository klonen

```bash
cd /opt
sudo git clone https://github.com/nextX-AG/ai-network-analyser.git ki-network-analyzer
cd ki-network-analyzer
```

### 3. Agent kompilieren

```bash
sudo go build -o agent cmd/agent/main.go
```

### 4. Berechtigungen setzen

```bash
sudo chmod +x agent
```

## Konfiguration

### 1. Konfigurationsverzeichnisse erstellen

```bash
sudo mkdir -p /etc/ki-network-analyzer
sudo mkdir -p /var/log/ki-network-analyzer
```

### 2. Basis-Konfiguration erstellen

```bash
sudo cat > /etc/ki-network-analyzer/agent.json << EOL
{
  "agent": {
    "listen": "0.0.0.0:8090",
    "server_url": "http://your-main-server:9090",
    "interface": "eth0",
    "name": "up-board-agent",
    "api_key": ""
  },
  "capture": {
    "promisc_mode": true,
    "snap_len": 65535,
    "buffer_size": 4194304
  },
  "gateway": {
    "detect_gateways": true,
    "track_dhcp": true,
    "track_dns": true,
    "track_arp": true
  }
}
EOL
```

Passen Sie die folgenden Werte an:
- `server_url`: IP-Adresse oder Hostname des Hauptservers
- `interface`: Netzwerkschnittstelle, auf der der Agent lauschen soll
- `name`: Ein eindeutiger Name für diesen Agent

### 3. Abhängigkeiten für Packet-Capturing installieren

```bash
sudo apt-get install -y libpcap-dev
```

### 4. Systemd-Service installieren

```bash
sudo cat > /etc/systemd/system/ki-network-analyzer-agent.service << EOL
[Unit]
Description=KI-Netzwerk-Analyzer Remote Agent
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/ki-network-analyzer
ExecStart=/opt/ki-network-analyzer/agent --config /etc/ki-network-analyzer/agent.json
Restart=always
RestartSec=5
AmbientCapabilities=CAP_NET_RAW CAP_NET_ADMIN
ProtectSystem=full
ProtectHome=read-only
PrivateTmp=true
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
EOL

# Service aktivieren und starten
sudo systemctl daemon-reload
sudo systemctl enable ki-network-analyzer-agent
sudo systemctl start ki-network-analyzer-agent
```

## Webinterface

Nach der Installation ist das Webinterface des Agents unter http://ip-des-up-boards:8090/admin verfügbar. Dort können Sie:

- Den Namen des Agents konfigurieren
- Die Server-URL ändern
- Die Erfassungsschnittstelle auswählen
- Den Agent neu starten
- Die Bridge-Konfiguration prüfen

## Bridge-Konfiguration für MITM-Monitoring

Für Man-in-the-Middle (MITM) Monitoring kann eine Netzwerk-Bridge konfiguriert werden. Dies ermöglicht die Analyse des gesamten Netzwerkverkehrs zwischen zwei Netzwerksegmenten.

### 1. Benötigte Pakete installieren

```bash
sudo apt-get update
sudo apt-get install -y bridge-utils
```

### 2. Bridge-Schnittstelle einrichten

Fügen Sie folgende Konfiguration in `/etc/network/interfaces` hinzu:

```
# Bridge-Schnittstelle
auto br0
iface br0 inet static
    address 192.168.1.254
    netmask 255.255.255.0
    network 192.168.1.0
    broadcast 192.168.1.255
    bridge_ports eth0 eth1
    bridge_stp off
    bridge_fd 0
    bridge_maxwait 0
```

Passen Sie die IP-Adressen und Schnittstellen (eth0, eth1) nach Bedarf an.

### 3. Netzwerkkonfiguration neuladen

```bash
sudo systemctl restart networking
```

### 4. Agent-Konfiguration anpassen

Öffnen Sie die Agent-Konfigurationsdatei:

```bash
sudo nano /etc/ki-network-analyzer/agent.json
```

Ändern Sie die Interface-Einstellung auf `br0`:

```json
"interface": "br0",
```

### 5. Agent neustarten

```bash
sudo systemctl restart ki-network-analyzer-agent
```

## Fehlerbehebung

### Agent startet nicht

Überprüfen Sie die Logs mit:

```bash
sudo journalctl -u ki-network-analyzer-agent -f
```

### Kompilierungsfehler

Stellen Sie sicher, dass alle Abhängigkeiten installiert sind:

```bash
sudo apt-get install -y golang-go libpcap-dev build-essential
```

### Keine Pakete werden erfasst

- Prüfen Sie, ob die angegebene Schnittstelle existiert (`ip a`)
- Stellen Sie sicher, dass der Agent mit Root-Rechten läuft
- Überprüfen Sie, ob der Promisc-Modus aktiviert ist (`ip link show`)
- Bei Bridge-Konfigurationen: Überprüfen Sie, ob beide Bridge-Ports korrekt verbunden sind

### Keine Verbindung zum Hauptserver

- Prüfen Sie, ob die Server-URL korrekt ist
- Stellen Sie sicher, dass der Hauptserver läuft und vom Agent erreichbar ist
- Überprüfen Sie eventuelle Firewall-Einstellungen
- Prüfen Sie den API-Schlüssel, falls konfiguriert

## Automatisches Update

Zur Zeit wird das automatische Update noch entwickelt. Bitte führen Sie Updates manuell durch, indem Sie die neueste Agent-Binärdatei herunterladen und den Dienst neu starten. 