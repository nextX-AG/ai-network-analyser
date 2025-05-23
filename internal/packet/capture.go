package packet

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"github.com/sayedamirkarim/ki-network-analyzer/internal/config"
	"github.com/sayedamirkarim/ki-network-analyzer/pkg/models"
)

// Capturer ist die Schnittstelle für Paketerfassung
type Capturer interface {
	OpenPcapFile(path string) error
	OpenLiveCapture(interfaceName string) error
	StartCapture(ctx context.Context) (<-chan *models.PacketInfo, <-chan error)
	Close() error
}

// PcapCapturer ist die Implementierung der Capturer-Schnittstelle
type PcapCapturer struct {
	config      *config.CaptureConfig
	gwConfig    *config.GatewayConfig
	handle      *pcap.Handle
	packetChan  chan *models.PacketInfo
	errorChan   chan error
	gatewayInfo *GatewayDetector
}

// GatewayDetector enthält Informationen über das erkannte Gateway
type GatewayDetector struct {
	knownGateways map[string]bool // IP-Adressen als Strings
	gatewayIP     net.IP
	gatewayMAC    net.HardwareAddr
	localNets     []*net.IPNet
	dhcpServers   map[string]bool   // DHCP-Server IPs
	dnsServers    map[string]bool   // DNS-Server IPs
	arpTable      map[string]string // IP zu MAC
}

// NewPcapCapturer erstellt einen neuen PcapCapturer
func NewPcapCapturer(cfg *config.Config) *PcapCapturer {
	gwDetector := &GatewayDetector{
		knownGateways: make(map[string]bool),
		dhcpServers:   make(map[string]bool),
		dnsServers:    make(map[string]bool),
		arpTable:      make(map[string]string),
	}

	// Bekannte Gateways hinzufügen
	for _, gw := range cfg.Gateway.KnownGateways {
		gwDetector.knownGateways[gw] = true
	}

	// Lokale Netzwerke erkennen
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				gwDetector.localNets = append(gwDetector.localNets, ipnet)
			}
		}
	}

	// Default-Gateway ermitteln
	defaultGW, _ := getDefaultGateway()
	if defaultGW != nil {
		gwDetector.gatewayIP = defaultGW
	}

	return &PcapCapturer{
		config:      &cfg.Capture,
		gwConfig:    &cfg.Gateway,
		packetChan:  make(chan *models.PacketInfo, 1000),
		errorChan:   make(chan error, 10),
		gatewayInfo: gwDetector,
	}
}

// OpenPcapFile öffnet eine PCAP-Datei zum Lesen
func (c *PcapCapturer) OpenPcapFile(path string) error {
	var err error
	c.handle, err = pcap.OpenOffline(path)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen der PCAP-Datei: %w", err)
	}

	if c.config.Filter != "" {
		if err := c.handle.SetBPFFilter(c.config.Filter); err != nil {
			return fmt.Errorf("Fehler beim Setzen des BPF-Filters: %w", err)
		}
	}

	return nil
}

// OpenLiveCapture öffnet eine Live-Netzwerkschnittstelle für die Erfassung
func (c *PcapCapturer) OpenLiveCapture(interfaceName string) error {
	var err error

	// Prüfen, ob es sich um eine Bridge-Schnittstelle handelt
	isBridge := false
	if _, err := os.Stat(fmt.Sprintf("/sys/class/net/%s/bridge", interfaceName)); err == nil {
		isBridge = true
		fmt.Printf("Bridge-Interface erkannt: %s\n", interfaceName)

		// Für Bridge-Interfaces: Promisc-Modus erzwingen und BufferSize erhöhen
		c.config.PromiscMode = true

		// Buffer-Größe für Bridge-Interfaces erhöhen
		if c.config.BufferSize < 8*1024*1024 {
			c.config.BufferSize = 8 * 1024 * 1024 // 8MB für Bridge-Interfaces
			fmt.Printf("Buffer-Größe für Bridge-Interface auf %d Bytes erhöht\n", c.config.BufferSize)
		}
	}

	// Konfigurieren der pcap-Bibliothek für Live-Capture
	inactive, err := pcap.NewInactiveHandle(interfaceName)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen des inaktiven Handles: %w", err)
	}
	defer inactive.CleanUp()

	// SnapLen setzen (maximale Paketgröße)
	if err := inactive.SetSnapLen(c.config.SnapLen); err != nil {
		return fmt.Errorf("Fehler beim Setzen von SnapLen: %w", err)
	}

	// Promisc-Modus setzen
	if err := inactive.SetPromisc(c.config.PromiscMode); err != nil {
		return fmt.Errorf("Fehler beim Setzen des Promisc-Modus: %w", err)
	}

	// Timeout setzen (BlockForever = -1)
	if err := inactive.SetTimeout(pcap.BlockForever); err != nil {
		return fmt.Errorf("Fehler beim Setzen des Timeouts: %w", err)
	}

	// Buffer-Größe setzen
	if err := inactive.SetBufferSize(c.config.BufferSize); err != nil {
		return fmt.Errorf("Fehler beim Setzen der Buffer-Größe: %w", err)
	}

	// Für Bridge-Interfaces: Immediate-Modus aktivieren (falls verfügbar)
	if isBridge {
		// Immediate-Modus ist plattformspezifisch, daher mit Fehlerbehandlung
		if err := inactive.SetImmediateMode(true); err != nil {
			// Fehler nur loggen, nicht abbrechen
			fmt.Printf("Warnung: Immediate-Modus konnte nicht aktiviert werden: %v\n", err)
		}
	}

	// Handle aktivieren
	c.handle, err = inactive.Activate()
	if err != nil {
		return fmt.Errorf("Fehler beim Aktivieren des Handles: %w", err)
	}

	// BPF-Filter setzen, falls konfiguriert
	if c.config.Filter != "" {
		if err := c.handle.SetBPFFilter(c.config.Filter); err != nil {
			return fmt.Errorf("Fehler beim Setzen des BPF-Filters: %w", err)
		}
	}

	fmt.Printf("Live-Capture auf Interface %s gestartet (Promisc: %v, SnapLen: %d, BufferSize: %d)\n",
		interfaceName, c.config.PromiscMode, c.config.SnapLen, c.config.BufferSize)

	return nil
}

// StartCapture startet die Erfassung von Paketen
func (c *PcapCapturer) StartCapture(ctx context.Context) (<-chan *models.PacketInfo, <-chan error) {
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	packetSource.DecodeOptions.Lazy = true
	packetSource.DecodeOptions.NoCopy = true

	// Debug-Ausgabe
	fmt.Printf("DEBUG: Starte Paketerfassung mit Linktyp: %v\n", c.handle.LinkType())

	// Konfigurierter Filter aus der Konfiguration ausgeben
	if c.config.Filter != "" {
		fmt.Printf("DEBUG: Konfigurierter BPF-Filter: %s\n", c.config.Filter)
	} else {
		fmt.Printf("DEBUG: Kein BPF-Filter konfiguriert\n")
	}

	// Debug-Zähler
	var packetCount uint64 = 0
	lastLogTime := time.Now()

	go func() {
		defer close(c.packetChan)
		defer close(c.errorChan)

		fmt.Println("DEBUG: Paketerfassungs-Goroutine gestartet")

		for {
			select {
			case <-ctx.Done():
				fmt.Println("DEBUG: Paketerfassung durch Kontext beendet")
				return
			case packet, ok := <-packetSource.Packets():
				if !ok {
					fmt.Println("DEBUG: Paketquelle geschlossen")
					return
				}

				// Paketzähler erhöhen
				packetCount++

				// Debug-Log alle 10 Pakete oder alle 5 Sekunden
				if packetCount%10 == 0 || time.Since(lastLogTime) > 5*time.Second {
					fmt.Printf("DEBUG: %d Pakete erfasst, letztes Paket: %d Bytes\n",
						packetCount, packet.Metadata().Length)
					lastLogTime = time.Now()
				}

				packetInfo, err := c.analyzePacket(packet)
				if err != nil {
					select {
					case c.errorChan <- err:
					default:
						// Errorkanal voll - ignorieren
					}
					continue
				}

				if packetInfo != nil {
					select {
					case c.packetChan <- packetInfo:
						// Debug-Info alle 50 Pakete
						if packetCount%50 == 0 {
							fmt.Printf("DEBUG: Paket an Kanal gesendet: %s -> %s (%s)\n",
								packetInfo.SourceIP, packetInfo.DestinationIP, packetInfo.Protocol)
						}
					default:
						// Kanal voll - ignorieren
						fmt.Println("DEBUG: Paketkanal ist voll - Paket verworfen")
					}
				} else {
					if packetCount%100 == 0 {
						fmt.Println("DEBUG: Paket wurde analysiert, aber nicht weitergeleitet (nil)")
					}
				}
			}
		}
	}()

	return c.packetChan, c.errorChan
}

// Close schließt den Capturer
func (c *PcapCapturer) Close() error {
	if c.handle != nil {
		c.handle.Close()
	}
	return nil
}

// analyzePacket analysiert ein einzelnes Paket mit Gateway-Fokus
func (c *PcapCapturer) analyzePacket(packet gopacket.Packet) (*models.PacketInfo, error) {
	info := &models.PacketInfo{
		Timestamp: packet.Metadata().Timestamp,
		Length:    uint32(packet.Metadata().Length),
		Protocol:  "Unknown",
	}

	// Link Layer (z.B. Ethernet)
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer != nil {
		eth, _ := ethLayer.(*layers.Ethernet)

		// ARP-Analyse
		if eth.EthernetType == layers.EthernetTypeARP {
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer != nil {
				arp, _ := arpLayer.(*layers.ARP)
				return c.analyzeARPPacket(packet, arp, info)
			}
		}
	}

	// Netzwerk Layer (IPv4/IPv6)
	var srcIP, dstIP net.IP
	var ipLayer gopacket.Layer

	// IPv4
	ipLayer = packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		srcIP = ip.SrcIP
		dstIP = ip.DstIP
		info.SourceIP = srcIP
		info.DestinationIP = dstIP
		info.TTL = ip.TTL

		// ICMP-Analyse
		if ip.Protocol == layers.IPProtocolICMPv4 {
			icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
			if icmpLayer != nil {
				// Wir extrahieren ICMP-Typ und -Code nicht, aber könnten das später hinzufügen
				// icmp, _ := icmpLayer.(*layers.ICMPv4)
				info.Protocol = "ICMP"

				// Prüfen, ob Gateway involviert ist
				isGatewayTraffic := c.isGatewayTraffic(srcIP, dstIP)
				info.IsGatewayTraffic = isGatewayTraffic

				if isGatewayTraffic {
					// Identifizieren, welche IP das Gateway ist
					if c.isGatewayIP(srcIP) {
						info.GatewayIP = srcIP
					} else if c.isGatewayIP(dstIP) {
						info.GatewayIP = dstIP
					}
				}

				return info, nil
			}
		}
	} else {
		// IPv6
		ipLayer = packet.Layer(layers.LayerTypeIPv6)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv6)
			srcIP = ip.SrcIP
			dstIP = ip.DstIP
			info.SourceIP = srcIP
			info.DestinationIP = dstIP
			info.TTL = ip.HopLimit
		} else {
			// Weder IPv4 noch IPv6 - vermutlich ARP oder anderes Link-Layer-Protokoll
			return info, nil
		}
	}

	// Transport Layer (TCP/UDP)
	var srcPort, dstPort uint16

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort = uint16(tcp.SrcPort)
		dstPort = uint16(tcp.DstPort)
		info.SourcePort = srcPort
		info.DestinationPort = dstPort
		info.Protocol = "TCP"
	} else {
		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			srcPort = uint16(udp.SrcPort)
			dstPort = uint16(udp.DstPort)
			info.SourcePort = srcPort
			info.DestinationPort = dstPort
			info.Protocol = "UDP"

			// DNS-Analyse (Port 53)
			if udp.SrcPort == 53 || udp.DstPort == 53 {
				dnsLayer := packet.Layer(layers.LayerTypeDNS)
				if dnsLayer != nil {
					dns, _ := dnsLayer.(*layers.DNS)
					return c.analyzeDNSPacket(packet, dns, info)
				}
			}

			// DHCP-Analyse (Port 67/68)
			if (udp.SrcPort == 67 && udp.DstPort == 68) || (udp.SrcPort == 68 && udp.DstPort == 67) {
				dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4)
				if dhcpLayer != nil {
					dhcp, _ := dhcpLayer.(*layers.DHCPv4)
					return c.analyzeDHCPPacket(packet, dhcp, info)
				}
			}
		}
	}

	// Gateway-Traffic erkennen
	isGatewayTraffic := c.isGatewayTraffic(srcIP, dstIP)
	info.IsGatewayTraffic = isGatewayTraffic

	if isGatewayTraffic {
		// Identifizieren, welche IP das Gateway ist
		if c.isGatewayIP(srcIP) {
			info.GatewayIP = srcIP
		} else if c.isGatewayIP(dstIP) {
			info.GatewayIP = dstIP
		}
	}

	return info, nil
}

// analyzeARPPacket analysiert ein ARP-Paket mit Fokus auf Gateway-Erkennung
func (c *PcapCapturer) analyzeARPPacket(packet gopacket.Packet, arp *layers.ARP, info *models.PacketInfo) (*models.PacketInfo, error) {
	info.Protocol = "ARP"

	// ARP-spezifische Informationen
	senderIP := net.IP(arp.SourceProtAddress)
	targetIP := net.IP(arp.DstProtAddress)
	senderMAC := net.HardwareAddr(arp.SourceHwAddress)
	targetMAC := net.HardwareAddr(arp.DstHwAddress)

	info.SourceIP = senderIP
	info.DestinationIP = targetIP

	// ARP-Info erstellen
	arpInfo := &models.ARPInfo{
		SenderIP:  senderIP,
		TargetIP:  targetIP,
		SenderMAC: senderMAC.String(),
		TargetMAC: targetMAC.String(),
	}

	// Operation bestimmen
	if arp.Operation == layers.ARPRequest {
		arpInfo.Operation = "REQUEST"
	} else if arp.Operation == layers.ARPReply {
		arpInfo.Operation = "REPLY"

		// ARP-Tabelle aktualisieren
		c.gatewayInfo.arpTable[senderIP.String()] = senderMAC.String()

		// Prüfen, ob dieses Gerät ein Gateway ist
		if c.isGatewayIP(senderIP) {
			c.gatewayInfo.gatewayIP = senderIP
			c.gatewayInfo.gatewayMAC = senderMAC
		}
	}

	// Gratuitous ARP erkennen (gleiche Quell- und Ziel-IP)
	if senderIP.Equal(targetIP) {
		arpInfo.IsGratuitous = true
	}

	// Prüfen, ob Gateway involviert ist
	info.IsGatewayTraffic = c.isGatewayIP(senderIP) || c.isGatewayIP(targetIP)
	if info.IsGatewayTraffic {
		if c.isGatewayIP(senderIP) {
			info.GatewayIP = senderIP
		} else if c.isGatewayIP(targetIP) {
			info.GatewayIP = targetIP
		}
	}

	info.ARPInfo = arpInfo
	return info, nil
}

// analyzeDNSPacket analysiert ein DNS-Paket mit Fokus auf Gateway-Erkennung
func (c *PcapCapturer) analyzeDNSPacket(packet gopacket.Packet, dns *layers.DNS, info *models.PacketInfo) (*models.PacketInfo, error) {
	info.Protocol = "DNS"

	// DNS-Server-IP merken
	if dns.QR {
		// Es ist eine Antwort, Quell-IP ist ein DNS-Server
		c.gatewayInfo.dnsServers[info.SourceIP.String()] = true
	}

	// DNS-Info erstellen
	dnsInfo := &models.DNSInfo{
		IsQuery:  !dns.QR,
		IsAnswer: dns.QR,
	}

	// Abfragen extrahieren
	for _, question := range dns.Questions {
		query := models.DNSQuery{
			Name:  string(question.Name),
			Type:  question.Type.String(),
			Class: question.Class.String(),
		}
		dnsInfo.Queries = append(dnsInfo.Queries, query)
	}

	// Antworten extrahieren
	for _, answer := range dns.Answers {
		dnsAnswer := models.DNSAnswer{
			Name:  string(answer.Name),
			Type:  answer.Type.String(),
			Class: answer.Class.String(),
			TTL:   answer.TTL,
		}

		// Verschiedene Record-Typen
		switch answer.Type {
		case layers.DNSTypeA:
			dnsAnswer.Data = net.IP(answer.IP).String()
		case layers.DNSTypeAAAA:
			dnsAnswer.Data = net.IP(answer.IP).String()
		case layers.DNSTypeMX:
			dnsAnswer.Data = fmt.Sprintf("%d %s", answer.MX.Preference, string(answer.MX.Name))
		case layers.DNSTypeNS:
			dnsAnswer.Data = string(answer.NS)
		case layers.DNSTypeCNAME:
			dnsAnswer.Data = string(answer.CNAME)
		case layers.DNSTypePTR:
			dnsAnswer.Data = string(answer.PTR)
		case layers.DNSTypeTXT:
			for _, txt := range answer.TXTs {
				dnsAnswer.Data += string(txt) + " "
			}
		default:
			dnsAnswer.Data = "Unsupported record type"
		}

		dnsInfo.Answers = append(dnsInfo.Answers, dnsAnswer)
	}

	// Prüfen, ob Gateway involviert ist
	isGatewayTraffic := c.isGatewayIP(info.SourceIP) || c.isGatewayIP(info.DestinationIP)
	info.IsGatewayTraffic = isGatewayTraffic

	if isGatewayTraffic {
		if c.isGatewayIP(info.SourceIP) {
			info.GatewayIP = info.SourceIP
		} else if c.isGatewayIP(info.DestinationIP) {
			info.GatewayIP = info.DestinationIP
		}
	}

	info.DNSInfo = dnsInfo
	return info, nil
}

// analyzeDHCPPacket analysiert ein DHCP-Paket mit Fokus auf Gateway-Erkennung
func (c *PcapCapturer) analyzeDHCPPacket(packet gopacket.Packet, dhcp *layers.DHCPv4, info *models.PacketInfo) (*models.PacketInfo, error) {
	info.Protocol = "DHCP"

	// DHCP-Info erstellen
	dhcpInfo := &models.DHCPInfo{
		ClientIP:  dhcp.ClientIP,
		YourIP:    dhcp.YourClientIP,
		ServerIP:  dhcp.NextServerIP,
		ClientMAC: dhcp.ClientHWAddr.String(),
	}

	// DHCP-Optionen auswerten
	for _, option := range dhcp.Options {
		switch option.Type {
		case layers.DHCPOptMessageType:
			// Nachrichtentyp ermitteln
			if len(option.Data) > 0 {
				switch option.Data[0] {
				case byte(DHCPMsgTypeDiscover):
					dhcpInfo.MessageType = "DISCOVER"
				case byte(DHCPMsgTypeOffer):
					dhcpInfo.MessageType = "OFFER"
				case byte(DHCPMsgTypeRequest):
					dhcpInfo.MessageType = "REQUEST"
				case byte(DHCPMsgTypeACK):
					dhcpInfo.MessageType = "ACK"
				}
			}
		case layers.DHCPOptRouter:
			// Gateway-Information
			if len(option.Data) >= 4 {
				dhcpInfo.GatewayIP = net.IP(option.Data[:4])

				// Gateway-Detektion aktualisieren
				if c.gwConfig.DetectGateways {
					c.gatewayInfo.knownGateways[dhcpInfo.GatewayIP.String()] = true
					c.gatewayInfo.gatewayIP = dhcpInfo.GatewayIP
				}
			}
		case layers.DHCPOptServerID:
			// DHCP-Server-IP
			if len(option.Data) >= 4 {
				serverIP := net.IP(option.Data[:4])
				c.gatewayInfo.dhcpServers[serverIP.String()] = true
			}
		case layers.DHCPOptDNS:
			// DNS-Server
			for i := 0; i < len(option.Data); i += 4 {
				if i+4 <= len(option.Data) {
					dnsServer := net.IP(option.Data[i : i+4])
					dhcpInfo.DNSServers = append(dhcpInfo.DNSServers, dnsServer)
					c.gatewayInfo.dnsServers[dnsServer.String()] = true
				}
			}
		case layers.DHCPOptLeaseTime:
			// Lease-Zeit
			if len(option.Data) >= 4 {
				dhcpInfo.LeaseTime = uint32(option.Data[0])<<24 |
					uint32(option.Data[1])<<16 |
					uint32(option.Data[2])<<8 |
					uint32(option.Data[3])
			}
		case layers.DHCPOptHostname:
			// Hostname
			dhcpInfo.ServerHostname = string(option.Data)
		}
	}

	// DHCP-Server als Gateway-Kandidat hinzufügen
	if c.gwConfig.DetectGateways && dhcpInfo.ServerIP != nil && !dhcpInfo.ServerIP.IsUnspecified() {
		c.gatewayInfo.knownGateways[dhcpInfo.ServerIP.String()] = true
	}

	// Prüfen, ob Gateway involviert ist
	info.IsGatewayTraffic = true // DHCP ist fast immer Gateway-relevant

	// Wenn wir Gateway kennen, setzen wir es
	if dhcpInfo.GatewayIP != nil && !dhcpInfo.GatewayIP.IsUnspecified() {
		info.GatewayIP = dhcpInfo.GatewayIP
	} else if c.isGatewayIP(info.SourceIP) {
		info.GatewayIP = info.SourceIP
	} else if c.isGatewayIP(info.DestinationIP) {
		info.GatewayIP = info.DestinationIP
	}

	info.DHCPInfo = dhcpInfo
	return info, nil
}

// isGatewayIP prüft, ob eine IP-Adresse ein Gateway ist
func (c *PcapCapturer) isGatewayIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	// Bekannte Gateways prüfen
	if c.gatewayInfo.knownGateways[ip.String()] {
		return true
	}

	// Erkanntes Gateway prüfen
	if c.gatewayInfo.gatewayIP != nil && ip.Equal(c.gatewayInfo.gatewayIP) {
		return true
	}

	// DHCP-Server sind oft Gateways
	if c.gatewayInfo.dhcpServers[ip.String()] {
		return true
	}

	return false
}

// isGatewayTraffic prüft, ob ein Paket mit Gateway-Traffic zu tun hat
func (c *PcapCapturer) isGatewayTraffic(srcIP, dstIP net.IP) bool {
	if c.isGatewayIP(srcIP) || c.isGatewayIP(dstIP) {
		return true
	}

	// Prüfen, ob eine der IPs extern ist (also nicht im lokalen Netz)
	srcIsLocal := false
	dstIsLocal := false

	for _, localNet := range c.gatewayInfo.localNets {
		if localNet.Contains(srcIP) {
			srcIsLocal = true
		}
		if localNet.Contains(dstIP) {
			dstIsLocal = true
		}
	}

	// Wenn eine IP lokal und die andere nicht lokal ist,
	// dann ist es wahrscheinlich Gateway-Traffic
	return srcIsLocal != dstIsLocal
}

// getDefaultGateway versucht, das Standard-Gateway zu ermitteln
func getDefaultGateway() (net.IP, error) {
	// Hinweis: Diese Funktion ist plattformunabhängig und nicht vollständig
	// Eine vollständige Implementierung würde OS-spezifischen Code erfordern

	// Auf Unix/Linux könnte man `netstat -rn` oder `ip route` parsen
	// Auf Windows könnte man WMI oder ähnliches verwenden

	// Diese vereinfachte Version gibt nil zurück, da wir Gateway-Erkennung
	// hauptsächlich über DHCP und ARP durchführen
	return nil, fmt.Errorf("Plattformspezifische Gateway-Erkennung nicht implementiert")
}

// UpdateInterface aktualisiert die in der Konfiguration verwendete Schnittstelle
func (c *PcapCapturer) UpdateInterface(interfaceName string) {
	if c.config != nil {
		c.config.Interface = interfaceName
		fmt.Printf("Schnittstelle in der Konfiguration aktualisiert auf: %s\n", interfaceName)
	}
}
