/**
 * Dienste für die Kommunikation mit der Remote-Agents API
 */

/**
 * Konvertiert komplexe Filter in BPF-Syntax
 * @param {Array} filters - Array von Filter-Objekten
 * @returns {string} BPF-Syntax als String
 */
export const convertToBpfSyntax = (filters) => {
  if (!filters || filters.length === 0) return '';
  
  return filters.map((filter, index) => {
    let bpfPart = '';
    
    // Füge den logischen Operator hinzu, wenn nicht der erste Filter
    if (index > 0) {
      bpfPart += filter.logicalOperator === 'and' ? ' and ' : ' or ';
    }
    
    switch (filter.type) {
      case 'ip':
        bpfPart += `${filter.subType === 'src' ? 'src' : 'dst'} host ${filter.value}`;
        break;
      case 'port':
        bpfPart += `${filter.subType === 'src' ? 'src' : 'dst'} port ${filter.value}`;
        break;
      case 'protocol':
        bpfPart += filter.subType.toLowerCase();
        break;
      case 'mac':
        bpfPart += `ether ${filter.subType === 'src' ? 'src' : 'dst'} ${filter.value}`;
        break;
      default:
        return '';
    }
    
    return bpfPart;
  }).join('');
};

/**
 * Ruft alle Agenten vom Server ab
 * @returns {Promise<Array>} Liste der Agenten
 */
export const fetchAgents = async () => {
  try {
    // API-Aufruf, um alle Agenten zu laden
    const response = await fetch('/api/agents');
    if (!response.ok) {
      throw new Error(`Fehler beim Laden der Agenten: ${response.statusText}`);
    }
    
    const data = await response.json();
    if (data.success) {
      return data.data || [];
    } else {
      throw new Error(data.error || 'Fehler beim Laden der Agenten');
    }
  } catch (err) {
    console.error('Fehler beim Laden der Agenten:', err);
    throw err;
  }
};

/**
 * Startet die Paketerfassung auf einem Agenten
 * @param {Object} agent - Agent-Objekt
 * @param {string} interfaceName - Name der Netzwerkschnittstelle
 * @param {string|Array} filter - Filter (String oder Array von Filter-Objekten)
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const startCapture = async (agent, interfaceName, filter) => {
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
    
    // Sende Anfrage an den Agenten
    const response = await fetch(`${agent.url}/capture/start`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestData),
    });
    
    return await response.json();
  } catch (error) {
    console.error(`Fehler beim Starten der Erfassung für Agent ${agent.id}:`, error);
    return { success: false, error: 'Verbindung zum Agenten fehlgeschlagen' };
  }
};

/**
 * Stoppt die Paketerfassung auf einem Agenten
 * @param {Object} agent - Agent-Objekt
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const stopCapture = async (agent) => {
  try {
    // Sende Anfrage an den Agenten
    const response = await fetch(`${agent.url}/capture/stop`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    return await response.json();
  } catch (error) {
    console.error(`Fehler beim Stoppen der Erfassung für Agent ${agent.id}:`, error);
    return { success: false, error: 'Verbindung zum Agenten fehlgeschlagen' };
  }
};

/**
 * Setzt die Netzwerkschnittstelle für einen Agenten
 * @param {Object} agent - Agent-Objekt
 * @param {string} interfaceName - Name der Netzwerkschnittstelle
 * @returns {Promise<Object>} Ergebnis der Operation
 */
export const setAgentInterface = async (agent, interfaceName) => {
  try {
    // Sende Anfrage an den Agenten
    const response = await fetch(`${agent.url}/capture/set-interface`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ interface: interfaceName }),
    });
    
    return await response.json();
  } catch (error) {
    console.error(`Fehler beim Setzen der Schnittstelle für Agent ${agent.id}:`, error);
    return { success: false, error: 'Verbindung zum Agenten fehlgeschlagen' };
  }
}; 