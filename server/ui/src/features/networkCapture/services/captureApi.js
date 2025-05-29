/**
 * Dienste für die Kommunikation mit der Netzwerkerfassung-API
 */

import api from '../../../services/api';
import { convertToBpfSyntax } from '../../../shared/utils/filterUtils';

/**
 * Holt den Status und die Interfaces vom Agenten
 * @param {string} agentId - ID des Agenten
 * @returns {Promise<Object>} Status und Schnittstellen des Agenten
 */
export const fetchAgentStatus = async (agentId) => {
  try {
    const response = await api.get(`/agents/${agentId}/status`);
    return response;
  } catch (error) {
    console.error('Fehler beim Laden des Agentenstatus:', error);
    throw error;
  }
};

/**
 * Startet die Paketerfassung
 * @param {string} agentId - ID des Agenten
 * @param {string} interfaceName - Name der Netzwerkschnittstelle
 * @param {string|Array} filter - Filter (String oder Array von Filter-Objekten)
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const startCapture = async (agentId, interfaceName, filter) => {
  try {
    const requestData = {
      interface: interfaceName
    };
    
    if (filter) {
      if (typeof filter === 'string') {
        requestData.filter = filter;
      } else if (Array.isArray(filter) && filter.length > 0) {
        requestData.filter = convertToBpfSyntax(filter);
      }
    }
    
    return await api.post(`/agents/${agentId}/capture/start`, requestData);
  } catch (error) {
    console.error('Fehler beim Starten der Erfassung:', error);
    throw error;
  }
};

/**
 * Stoppt die Paketerfassung
 * @param {string} agentId - ID des Agenten
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const stopCapture = async (agentId) => {
  try {
    return await api.post(`/agents/${agentId}/capture/stop`);
  } catch (error) {
    console.error('Fehler beim Stoppen der Erfassung:', error);
    throw error;
  }
};

/**
 * Setzt die Netzwerkschnittstelle
 * @param {string} agentId - ID des Agenten
 * @param {string} interfaceName - Name der Netzwerkschnittstelle
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const setInterface = async (agentId, interfaceName) => {
  try {
    return await api.post(`/agents/${agentId}/capture/set-interface`, {
      interface: interfaceName
    });
  } catch (error) {
    console.error('Fehler beim Setzen der Schnittstelle:', error);
    throw error;
  }
}; 