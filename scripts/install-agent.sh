#!/bin/bash

# KI-Netzwerk-Analyzer Agent Installationsskript
# --------------------------------------------

set -e  # Skript bei Fehlern beenden

# Farben für Ausgabe
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}KI-Netzwerk-Analyzer Remote-Agent Installation${NC}"
echo "---------------------------------------------"
echo

# Prüfen, ob Script mit Root-Rechten läuft
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Bitte führen Sie dieses Skript als root aus (sudo).${NC}"
  exit 1
fi

# Standardwerte
SERVER_URL="http://localhost:9090"
AGENT_NAME=$(hostname)
INTERFACE="eth0"
INSTALL_DIR="/opt/ki-network-analyzer"
CONFIG_DIR="/etc/ki-network-analyzer"
LOG_DIR="/var/log/ki-network-analyzer"

# Parameter verarbeiten
while [ $# -gt 0 ]; do
  case "$1" in
    --server-url=*)
      SERVER_URL="${1#*=}"
      ;;
    --name=*)
      AGENT_NAME="${1#*=}"
      ;;
    --interface=*)
      INTERFACE="${1#*=}"
      ;;
    --help)
      echo "Verwendung: $0 [Optionen]"
      echo "Optionen:"
      echo "  --server-url=URL    URL des Hauptservers (default: http://localhost:9090)"
      echo "  --name=NAME         Name des Agents (default: Hostname)"
      echo "  --interface=IFACE   Zu verwendende Netzwerkschnittstelle (default: eth0)"
      echo "  --help              Diese Hilfe anzeigen"
      exit 0
      ;;
    *)
      echo -e "${RED}Unbekannte Option: $1${NC}"
      echo "Verwenden Sie --help für Hilfe."
      exit 1
      ;;
  esac
  shift
done

echo -e "${GREEN}Konfiguration:${NC}"
echo "Server URL:   $SERVER_URL"
echo "Agent Name:   $AGENT_NAME"
echo "Interface:    $INTERFACE"
echo "Install Dir:  $INSTALL_DIR"
echo

# Abhängigkeiten installieren
echo -e "${YELLOW}Abhängigkeiten installieren...${NC}"
apt-get update
apt-get install -y golang-go libpcap-dev git build-essential

# Go-Version prüfen
GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
GO_VERSION_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_VERSION_MINOR=$(echo $GO_VERSION | cut -d. -f2)

if [ "$GO_VERSION_MAJOR" -lt 1 ] || ([ "$GO_VERSION_MAJOR" -eq 1 ] && [ "$GO_VERSION_MINOR" -lt 16 ]); then
  echo -e "${RED}Go Version $GO_VERSION ist zu alt. Version 1.16 oder höher wird benötigt.${NC}"
  exit 1
fi

echo -e "${GREEN}Go $GO_VERSION gefunden.${NC}"

# Verzeichnisse erstellen
echo -e "${YELLOW}Verzeichnisse vorbereiten...${NC}"
mkdir -p $INSTALL_DIR
mkdir -p $CONFIG_DIR
mkdir -p $LOG_DIR

# Repository klonen
echo -e "${YELLOW}Repository klonen...${NC}"
cd /tmp
if [ -d "ai-network-analyser" ]; then
  rm -rf ai-network-analyser
fi

git clone https://github.com/nextX-AG/ai-network-analyser.git
cd ai-network-analyser

# Agent kompilieren
echo -e "${YELLOW}Agent kompilieren...${NC}"
go build -o agent cmd/agent/main.go
if [ ! -f "agent" ]; then
  echo -e "${RED}Kompilierung fehlgeschlagen!${NC}"
  exit 1
fi

# Agent installieren
echo -e "${YELLOW}Agent installieren...${NC}"
cp agent $INSTALL_DIR/
chmod +x $INSTALL_DIR/agent

# Konfiguration erstellen
echo -e "${YELLOW}Konfiguration erstellen...${NC}"
cat > $CONFIG_DIR/agent.json << EOL
{
  "agent": {
    "listen": "0.0.0.0:8090",
    "server_url": "$SERVER_URL",
    "interface": "$INTERFACE",
    "name": "$AGENT_NAME",
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

# Systemd Service erstellen
echo -e "${YELLOW}Systemd Service erstellen...${NC}"
cat > /etc/systemd/system/ki-network-analyzer-agent.service << EOL
[Unit]
Description=KI-Netzwerk-Analyzer Remote Agent
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/agent --config $CONFIG_DIR/agent.json
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
echo -e "${YELLOW}Service aktivieren und starten...${NC}"
systemctl daemon-reload
systemctl enable ki-network-analyzer-agent
systemctl start ki-network-analyzer-agent

# Status überprüfen
sleep 2
if systemctl is-active --quiet ki-network-analyzer-agent; then
  echo -e "${GREEN}Agent wurde erfolgreich installiert und gestartet!${NC}"
else
  echo -e "${RED}Agent konnte nicht gestartet werden. Bitte prüfen Sie den Status mit: systemctl status ki-network-analyzer-agent${NC}"
  exit 1
fi

# Webinterface-URL anzeigen
IP_ADDR=$(hostname -I | awk '{print $1}')
echo
echo -e "${GREEN}Installation abgeschlossen!${NC}"
echo -e "Agent Webinterface ist verfügbar unter: ${YELLOW}http://$IP_ADDR:8090/admin${NC}"
echo -e "Sie können die Konfiguration dort anpassen und den Agent bei Bedarf neu starten."
echo
echo -e "Wenn Sie eine Bridge für MITM-Monitoring einrichten möchten, können Sie das Webinterface verwenden oder die Anleitung in der Dokumentation befolgen."
echo 