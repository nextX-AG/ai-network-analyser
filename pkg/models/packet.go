package models

import (
	"net"
	"time"
)

// PacketInfo repräsentiert die wesentlichen Informationen eines Netzwerkpakets
type PacketInfo struct {
	Timestamp       time.Time `json:"timestamp"`
	SourceIP        net.IP    `json:"source_ip"`
	DestinationIP   net.IP    `json:"destination_ip"`
	SourcePort      uint16    `json:"source_port,omitempty"`
	DestinationPort uint16    `json:"destination_port,omitempty"`
	Protocol        string    `json:"protocol"`
	Length          uint32    `json:"length"`
	TTL             uint8     `json:"ttl,omitempty"`

	// Gateway-relevante Informationen
	IsGatewayTraffic bool      `json:"is_gateway_traffic"`
	GatewayIP        net.IP    `json:"gateway_ip,omitempty"`
	NATInfo          *NATInfo  `json:"nat_info,omitempty"`
	DNSInfo          *DNSInfo  `json:"dns_info,omitempty"`
	DHCPInfo         *DHCPInfo `json:"dhcp_info,omitempty"`
	ARPInfo          *ARPInfo  `json:"arp_info,omitempty"`

	// Rohpaketdaten für detaillierte Analyse
	RawData []byte `json:"-"`
}

// NATInfo enthält Informationen zu NAT-Übersetzungen
type NATInfo struct {
	OriginalSourceIP        net.IP `json:"original_source_ip,omitempty"`
	OriginalDestinationIP   net.IP `json:"original_destination_ip,omitempty"`
	OriginalSourcePort      uint16 `json:"original_source_port,omitempty"`
	OriginalDestinationPort uint16 `json:"original_destination_port,omitempty"`
	TranslationType         string `json:"translation_type,omitempty"` // SNAT, DNAT, PAT, etc.
}

// DNSInfo enthält DNS-spezifische Informationen
type DNSInfo struct {
	Queries  []DNSQuery  `json:"queries,omitempty"`
	Answers  []DNSAnswer `json:"answers,omitempty"`
	IsQuery  bool        `json:"is_query"`
	IsAnswer bool        `json:"is_answer"`
}

// DNSQuery repräsentiert eine DNS-Anfrage
type DNSQuery struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
}

// DNSAnswer repräsentiert eine DNS-Antwort
type DNSAnswer struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
	TTL   uint32 `json:"ttl"`
	Data  string `json:"data"`
}

// DHCPInfo enthält DHCP-spezifische Informationen
type DHCPInfo struct {
	MessageType    string   `json:"message_type,omitempty"` // DISCOVER, OFFER, REQUEST, ACK
	ClientIP       net.IP   `json:"client_ip,omitempty"`
	YourIP         net.IP   `json:"your_ip,omitempty"`
	ServerIP       net.IP   `json:"server_ip,omitempty"`
	GatewayIP      net.IP   `json:"gateway_ip,omitempty"`
	ClientMAC      string   `json:"client_mac,omitempty"`
	ServerHostname string   `json:"server_hostname,omitempty"`
	DNSServers     []net.IP `json:"dns_servers,omitempty"`
	LeaseTime      uint32   `json:"lease_time,omitempty"`
}

// ARPInfo enthält ARP-spezifische Informationen
type ARPInfo struct {
	Operation    string `json:"operation,omitempty"` // REQUEST, REPLY
	SenderMAC    string `json:"sender_mac,omitempty"`
	SenderIP     net.IP `json:"sender_ip,omitempty"`
	TargetMAC    string `json:"target_mac,omitempty"`
	TargetIP     net.IP `json:"target_ip,omitempty"`
	IsGratuitous bool   `json:"is_gratuitous,omitempty"`
}

// PacketSummary enthält eine kompakte Zusammenfassung des Pakets
type PacketSummary struct {
	Timestamp        time.Time `json:"timestamp"`
	SourceIP         string    `json:"source_ip"`
	DestinationIP    string    `json:"destination_ip"`
	Protocol         string    `json:"protocol"`
	Length           uint32    `json:"length"`
	IsGatewayTraffic bool      `json:"is_gateway_traffic"`
	Summary          string    `json:"summary"`
	EventType        string    `json:"event_type,omitempty"` // Normal, Warning, Error
}

// GatewayEvent repräsentiert ein Gateway-relevantes Ereignis
type GatewayEvent struct {
	Timestamp      time.Time   `json:"timestamp"`
	EventType      string      `json:"event_type"` // "dhcp", "dns", "nat", "arp"
	Description    string      `json:"description"`
	Severity       string      `json:"severity"`                  // "info", "warning", "error"
	RelatedPackets []uint64    `json:"related_packets,omitempty"` // Paket-IDs
	GatewayIP      string      `json:"gateway_ip,omitempty"`
	ClientIP       string      `json:"client_ip,omitempty"`
	Data           interface{} `json:"data,omitempty"` // Typspezifische Daten
}
