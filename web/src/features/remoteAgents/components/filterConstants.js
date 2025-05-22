/**
 * Konstanten für die Filter-Komponenten
 */

// IP-Filter-Typen
export const IP_FILTER_TYPES = [
  { value: 'src', label: 'Quell-IP' },
  { value: 'dst', label: 'Ziel-IP' }
];

// Port-Filter-Typen
export const PORT_FILTER_TYPES = [
  { value: 'src', label: 'Quellport' },
  { value: 'dst', label: 'Zielport' }
];

// Protokoll-Filter-Typen
export const PROTOCOL_FILTER_TYPES = [
  { value: 'tcp', label: 'TCP' },
  { value: 'udp', label: 'UDP' },
  { value: 'icmp', label: 'ICMP' },
  { value: 'arp', label: 'ARP' },
  { value: 'ip', label: 'IP' },
  { value: 'ipv6', label: 'IPv6' },
  { value: 'http', label: 'HTTP' },
  { value: 'https', label: 'HTTPS' },
  { value: 'dns', label: 'DNS' }
];

// MAC-Filter-Typen
export const MAC_FILTER_TYPES = [
  { value: 'src', label: 'Quell-MAC' },
  { value: 'dst', label: 'Ziel-MAC' }
];

// Logische Operatoren
export const LOGICAL_OPERATORS = [
  { value: 'and', label: 'UND' },
  { value: 'or', label: 'ODER' }
];

// Häufig verwendete Ports
export const COMMON_PORTS = [
  { value: '80', label: 'HTTP' },
  { value: '443', label: 'HTTPS' },
  { value: '53', label: 'DNS' },
  { value: '22', label: 'SSH' },
  { value: '21', label: 'FTP' },
  { value: '25', label: 'SMTP' },
  { value: '110', label: 'POP3' },
  { value: '143', label: 'IMAP' },
  { value: '3306', label: 'MySQL' },
  { value: '5432', label: 'PostgreSQL' },
  { value: '1433', label: 'MS SQL' },
  { value: '27017', label: 'MongoDB' },
  { value: '6379', label: 'Redis' },
  { value: '8080', label: 'Alternative HTTP' },
  { value: '8443', label: 'Alternative HTTPS' }
]; 