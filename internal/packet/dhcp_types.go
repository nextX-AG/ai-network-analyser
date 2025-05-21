package packet

// DHCP-Nachrichtentypen (fehlen in gopacket/layers)
const (
	DHCPMsgTypeDiscover = 1
	DHCPMsgTypeOffer    = 2
	DHCPMsgTypeRequest  = 3
	DHCPMsgTypeDecline  = 4
	DHCPMsgTypeACK      = 5
	DHCPMsgTypeNAK      = 6
	DHCPMsgTypeRelease  = 7
	DHCPMsgTypeInform   = 8
)
