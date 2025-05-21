# Installation des KI-Netzwerk-Analyzer Remote Agents

Diese Anleitung beschreibt die Installation und Konfiguration des Remote-Agents auf einem UP Board oder einem anderen Linux-System.

## Systemanforderungen

- Linux-System (Ubuntu/Debian empfohlen)
- Go 1.16+ (falls Kompilierung aus Quellcode gewünscht)
- Administratorrechte (sudo)

## Binäre Installation

### 1. Vorgefertigte Binärdatei herunterladen

```bash
mkdir -p /opt/ki-network-analyzer
cd /opt/ki-network-analyzer
wget https://github.com/sayedamirkarim/ki-network-analyzer/releases/latest/download/agent-linux-amd64 -O agent
chmod +x agent
```

### 2. Konfigurationsverzeichnisse erstellen

```bash
mkdir -p /etc/ki-network-analyzer
mkdir -p /var/log/ki-network-analyzer
```

### 3. Basis-Konfiguration erstellen

```bash
cat > /etc/ki-network-analyzer/agent.json << EOL
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
  }
}
EOL
```

### 4. Systemd-Service installieren

```bash
cat > /etc/systemd/system/ki-network-analyzer-agent.service << EOL
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
systemctl daemon-reload
systemctl enable ki-network-analyzer-agent
systemctl start ki-network-analyzer-agent
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
apt-get update
apt-get install -y bridge-utils
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
systemctl restart networking
```

### 4. Agent-Konfiguration anpassen

Öffnen Sie die Agent-Konfigurationsdatei:

```bash
nano /etc/ki-network-analyzer/agent.json
```

Ändern Sie die Interface-Einstellung auf `br0`:

```json
"interface": "br0",
```

### 5. Agent neustarten

```bash
systemctl restart ki-network-analyzer-agent
```

## Fehlerbehebung

### Agent startet nicht

Überprüfen Sie die Logs mit:

```bash
journalctl -u ki-network-analyzer-agent -f
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