import api from './api';

export interface NetworkInterface {
  name: string;
  index: number;
  macAddress: string;
  ipAddresses: string[];
  isUp: boolean;
  isLoopback: boolean;
}

export interface CaptureOptions {
  interface: string;
  filter?: string;
}

const networkService = {
  // Netzwerkschnittstellen abrufen
  getInterfaces: () => {
    return api.get<NetworkInterface[]>('/interfaces');
  },

  // PCAP-Datei analysieren
  analyzePcap: (file: File) => {
    const formData = new FormData();
    formData.append('pcap', file);
    return api.post('/analyze', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },

  // Live-Capture starten
  startLiveCapture: (options: CaptureOptions) => {
    return api.post('/live/start', options);
  },

  // Live-Capture stoppen
  stopLiveCapture: () => {
    return api.post('/live/stop');
  },

  // Gateway-Informationen abrufen
  getGateways: () => {
    return api.get('/gateways');
  },

  // Gateway-Traffic-Statistiken abrufen
  getGatewayTraffic: () => {
    return api.get('/traffic/gateway');
  },

  // Gateway-Ereignisse abrufen
  getGatewayEvents: () => {
    return api.get('/events/gateway');
  },
};

export default networkService; 