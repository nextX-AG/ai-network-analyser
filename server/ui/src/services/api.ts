import axios from 'axios';

// Basis-URL für alle API-Anfragen
const baseURL = process.env.REACT_APP_API_URL || '/api';

// Erstelle eine Axios-Instanz mit der Basis-Konfiguration
const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor für einheitliche Fehlerbehandlung
api.interceptors.request.use(
  (config) => {
    // Hier können wir später Auth-Token hinzufügen
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response Interceptor für einheitliche Fehlerbehandlung
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    // Einheitliche Fehlerbehandlung
    if (error.response) {
      // Server antwortet mit Fehler
      console.error('API Error:', error.response.data);
    } else if (error.request) {
      // Keine Antwort vom Server
      console.error('Network Error:', error.request);
    } else {
      // Fehler beim Erstellen der Anfrage
      console.error('Request Error:', error.message);
    }
    return Promise.reject(error);
  }
);

export default api; 