/**
 * Hilfsfunktionen für die Filterverarbeitung
 */

/**
 * Konvertiert strukturierte Filter in BPF-Syntax
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