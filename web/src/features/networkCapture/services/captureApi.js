/**
 * Dienste für die Kommunikation mit der Netzwerkerfassung-API
 */

import { convertToBpfSyntax } from '../../../shared/utils/filterUtils';

/**
 * Holt den Status und die Interfaces vom Agenten
 * @param {string} agentUrl - URL des Agenten
 * @returns {Promise<Object>} Status und Schnittstellen des Agenten
 */
export const fetchAgentStatus = async (agentUrl) => {
  try {
    const response = await fetch(`${agentUrl}/status`);
    const data = await response.json();
    
    if (data.success) {
      return data.data;
    } else {
      throw new Error(data.error || 'Fehler beim Laden des Agentenstatus');
    }
  } catch (error) {
    console.error('Fehler beim Laden des Agentenstatus:', error);
    throw error;
  }
};

/**
 * Startet die Paketerfassung
 * @param {string} agentUrl - URL des Agenten
 * @param {string} interfaceName - Name der Netzwerkschnittstelle
 * @param {string|Array} filter - Filter (String oder Array von Filter-Objekten)
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const startCapture = async (agentUrl, interfaceName, filter) => {
  try {
    // Bereite Anfrage vor
    const requestData = {
      interface: interfaceName
    };
    
    // Füge Filter hinzu, wenn vorhanden
    if (filter) {
      if (typeof filter === 'string') {
        requestData.filter = filter;
      } else if (Array.isArray(filter) && filter.length > 0) {
        // Konvertiere komplexe Filter in BPF-Syntax
        const bpfFilter = convertToBpfSyntax(filter);
        requestData.filter = bpfFilter;
      }
    }
    
    const response = await fetch(`${agentUrl}/capture/start`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestData),
    });
    
    return await response.json();
  } catch (error) {
    console.error('Fehler beim Starten der Erfassung:', error);
    return { success: false, error: 'Verbindung zum Agenten fehlgeschlagen' };
  }
};

/**
 * Stoppt die Paketerfassung
 * @param {string} agentUrl - URL des Agenten
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const stopCapture = async (agentUrl) => {
  try {
    const response = await fetch(`${agentUrl}/capture/stop`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    return await response.json();
  } catch (error) {
    console.error('Fehler beim Stoppen der Erfassung:', error);
    return { success: false, error: 'Verbindung zum Agenten fehlgeschlagen' };
  }
};

/**
 * Setzt die Netzwerkschnittstelle
 * @param {string} agentUrl - URL des Agenten
 * @param {string} interfaceName - Name der Netzwerkschnittstelle
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const setInterface = async (agentUrl, interfaceName) => {
  try {
    const response = await fetch(`${agentUrl}/capture/set-interface`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ interface: interfaceName }),
    });
    
    return await response.json();
  } catch (error) {
    console.error('Fehler beim Setzen der Schnittstelle:', error);
    return { success: false, error: 'Verbindung zum Agenten fehlgeschlagen' };
  }
}; 