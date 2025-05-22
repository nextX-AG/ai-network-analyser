import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Divider,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Stack,
  TextField,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';
import FilterAltIcon from '@mui/icons-material/FilterAlt';
import SaveIcon from '@mui/icons-material/Save';

// Importieren der Filterkonstanten aus der remoteAgents-Komponente
// In einer vollständigen Implementierung würden diese in ein shared-Verzeichnis verschoben
const IP_FILTER_TYPES = [
  { value: 'src', label: 'Quell-IP' },
  { value: 'dst', label: 'Ziel-IP' }
];

const PORT_FILTER_TYPES = [
  { value: 'src', label: 'Quellport' },
  { value: 'dst', label: 'Zielport' }
];

const PROTOCOL_FILTER_TYPES = [
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

const MAC_FILTER_TYPES = [
  { value: 'src', label: 'Quell-MAC' },
  { value: 'dst', label: 'Ziel-MAC' }
];

const LOGICAL_OPERATORS = [
  { value: 'and', label: 'UND' },
  { value: 'or', label: 'ODER' }
];

const COMMON_PORTS = [
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

/**
 * CaptureFilter - Komponente zur Filterung der Netzwerkerfassung
 */
const CaptureFilter = ({ onApplyFilter, currentFilter }) => {
  const [expanded, setExpanded] = useState(false);
  const [activeFilters, setActiveFilters] = useState([]);
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [manualBpfFilter, setManualBpfFilter] = useState('');
  const [savedFilters, setSavedFilters] = useState([]);
  const [newFilterName, setNewFilterName] = useState('');
  const [showSaveDialog, setShowSaveDialog] = useState(false);
  
  // Form states
  const [filterType, setFilterType] = useState('ip');
  const [ipFilter, setIpFilter] = useState({
    type: 'src',
    value: '',
    operator: 'equals'
  });
  const [portFilter, setPortFilter] = useState({
    type: 'src',
    value: '',
    operator: 'equals'
  });
  const [protocolFilter, setProtocolFilter] = useState({
    type: 'tcp',
    operator: 'equals'
  });
  const [macFilter, setMacFilter] = useState({
    type: 'src',
    value: '',
    operator: 'equals'
  });
  const [logicalOperator, setLogicalOperator] = useState('and');

  // Lade gespeicherte Filter beim ersten Rendern
  useEffect(() => {
    const loadedFilters = localStorage.getItem('savedCaptureFilters');
    if (loadedFilters) {
      setSavedFilters(JSON.parse(loadedFilters));
    }
  }, []);

  // Aktualisiere die aktiven Filter, wenn currentFilter sich ändert
  useEffect(() => {
    if (currentFilter) {
      // Wenn ein manueller BPF-Filter ist
      if (typeof currentFilter === 'string') {
        setManualBpfFilter(currentFilter);
        setShowAdvanced(true);
      } 
      // Wenn ein strukturierter Filter ist
      else if (Array.isArray(currentFilter)) {
        setActiveFilters(currentFilter);
      }
    }
  }, [currentFilter]);

  const handleAccordionChange = () => {
    setExpanded(!expanded);
  };

  const handleFilterTypeChange = (event) => {
    setFilterType(event.target.value);
  };

  const handleAddFilter = () => {
    let newFilter;
    
    switch(filterType) {
      case 'ip':
        newFilter = {
          id: Date.now(),
          type: 'ip',
          subType: ipFilter.type,
          value: ipFilter.value,
          operator: ipFilter.operator,
          display: `${ipFilter.type} IP ${ipFilter.operator} ${ipFilter.value}`
        };
        break;
      case 'port':
        newFilter = {
          id: Date.now(),
          type: 'port',
          subType: portFilter.type,
          value: portFilter.value,
          operator: portFilter.operator,
          display: `${portFilter.type} port ${portFilter.operator} ${portFilter.value}`
        };
        break;
      case 'protocol':
        newFilter = {
          id: Date.now(),
          type: 'protocol',
          subType: protocolFilter.type,
          operator: protocolFilter.operator,
          display: `protocol ${protocolFilter.operator} ${protocolFilter.type}`
        };
        break;
      case 'mac':
        newFilter = {
          id: Date.now(),
          type: 'mac',
          subType: macFilter.type,
          value: macFilter.value,
          operator: macFilter.operator,
          display: `${macFilter.type} MAC ${macFilter.operator} ${macFilter.value}`
        };
        break;
      default:
        return;
    }
    
    // Füge logischen Operator hinzu, wenn dies nicht der erste Filter ist
    if (activeFilters.length > 0) {
      newFilter.logicalOperator = logicalOperator;
    }
    
    setActiveFilters([...activeFilters, newFilter]);
    
    // Leere das aktuelle Filterfeld
    if (filterType === 'ip') {
      setIpFilter({ ...ipFilter, value: '' });
    } else if (filterType === 'port') {
      setPortFilter({ ...portFilter, value: '' });
    } else if (filterType === 'mac') {
      setMacFilter({ ...macFilter, value: '' });
    }
  };

  const handleRemoveFilter = (id) => {
    setActiveFilters(activeFilters.filter(filter => filter.id !== id));
  };

  const handleApplyFilter = () => {
    // Prüfe, ob wir einen manuellen BPF-Filter oder strukturierte Filter verwenden
    if (showAdvanced && manualBpfFilter) {
      onApplyFilter(manualBpfFilter);
    } else {
      onApplyFilter(activeFilters);
    }
  };

  const handleSaveFilter = () => {
    const newFilter = {
      id: Date.now(),
      name: newFilterName,
      filters: showAdvanced ? manualBpfFilter : activeFilters
    };
    
    const updatedFilters = [...savedFilters, newFilter];
    setSavedFilters(updatedFilters);
    localStorage.setItem('savedCaptureFilters', JSON.stringify(updatedFilters));
    
    setNewFilterName('');
    setShowSaveDialog(false);
  };

  const handleLoadFilter = (filterToLoad) => {
    if (typeof filterToLoad.filters === 'string') {
      setManualBpfFilter(filterToLoad.filters);
      setShowAdvanced(true);
      setActiveFilters([]);
    } else {
      setActiveFilters(filterToLoad.filters);
      setShowAdvanced(false);
      setManualBpfFilter('');
    }
  };

  const handleDeleteSavedFilter = (id) => {
    const updatedFilters = savedFilters.filter(filter => filter.id !== id);
    setSavedFilters(updatedFilters);
    localStorage.setItem('savedCaptureFilters', JSON.stringify(updatedFilters));
  };

  // Konvertiert die aktiven Filter in BPF-Syntax
  const generateBpfFilter = () => {
    if (activeFilters.length === 0) return '';
    
    return activeFilters.map((filter, index) => {
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

  return (
    <Accordion expanded={expanded} onChange={handleAccordionChange}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Stack direction="row" spacing={1} alignItems="center">
          <FilterAltIcon />
          <Typography variant="h6">Netzwerk-Filter</Typography>
          {activeFilters.length > 0 && (
            <Chip 
              label={`${activeFilters.length} Filter aktiv`} 
              color="primary" 
              size="small" 
            />
          )}
        </Stack>
      </AccordionSummary>
      <AccordionDetails>
        <Card variant="outlined">
          <CardContent>
            <Grid container spacing={2}>
              {/* Umschalter zwischen einfachen und fortgeschrittenen Filtern */}
              <Grid item xs={12}>
                <FormControl fullWidth>
                  <Stack direction="row" spacing={2} justifyContent="space-between">
                    <Button 
                      variant={!showAdvanced ? "contained" : "outlined"} 
                      onClick={() => setShowAdvanced(false)}
                    >
                      Einfacher Filter
                    </Button>
                    <Button 
                      variant={showAdvanced ? "contained" : "outlined"} 
                      onClick={() => setShowAdvanced(true)}
                    >
                      BPF-Syntax (Fortgeschritten)
                    </Button>
                  </Stack>
                </FormControl>
              </Grid>

              {/* Gespeicherte Filter */}
              {savedFilters.length > 0 && (
                <Grid item xs={12}>
                  <Typography variant="subtitle1">Gespeicherte Filter</Typography>
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 1 }}>
                    {savedFilters.map((filter) => (
                      <Chip
                        key={filter.id}
                        label={filter.name}
                        onClick={() => handleLoadFilter(filter)}
                        onDelete={() => handleDeleteSavedFilter(filter.id)}
                        color="secondary"
                        deleteIcon={<DeleteIcon />}
                      />
                    ))}
                  </Box>
                </Grid>
              )}

              {showAdvanced ? (
                // Fortgeschrittener BPF-Filter
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    multiline
                    rows={3}
                    label="BPF Filter Syntax"
                    placeholder="z.B. tcp port 80 or udp port 53"
                    value={manualBpfFilter}
                    onChange={(e) => setManualBpfFilter(e.target.value)}
                    helperText="Berkeley Packet Filter Syntax"
                    variant="outlined"
                  />
                </Grid>
              ) : (
                // Einfache Filter-Oberfläche
                <>
                  <Grid item xs={12} sm={3}>
                    <FormControl fullWidth>
                      <InputLabel>Filtertyp</InputLabel>
                      <Select
                        value={filterType}
                        onChange={handleFilterTypeChange}
                        label="Filtertyp"
                      >
                        <MenuItem value="ip">IP-Adresse</MenuItem>
                        <MenuItem value="port">Port</MenuItem>
                        <MenuItem value="protocol">Protokoll</MenuItem>
                        <MenuItem value="mac">MAC-Adresse</MenuItem>
                      </Select>
                    </FormControl>
                  </Grid>

                  {/* Dynamische Filter-Optionen basierend auf dem ausgewählten Filtertyp */}
                  {filterType === 'ip' && (
                    <>
                      <Grid item xs={12} sm={3}>
                        <FormControl fullWidth>
                          <InputLabel>Typ</InputLabel>
                          <Select
                            value={ipFilter.type}
                            onChange={(e) => setIpFilter({...ipFilter, type: e.target.value})}
                            label="Typ"
                          >
                            {IP_FILTER_TYPES.map(type => (
                              <MenuItem key={type.value} value={type.value}>{type.label}</MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                      </Grid>
                      <Grid item xs={12} sm={4}>
                        <TextField
                          fullWidth
                          label="IP-Adresse"
                          placeholder="z.B. 192.168.1.1"
                          value={ipFilter.value}
                          onChange={(e) => setIpFilter({...ipFilter, value: e.target.value})}
                        />
                      </Grid>
                    </>
                  )}

                  {filterType === 'port' && (
                    <>
                      <Grid item xs={12} sm={3}>
                        <FormControl fullWidth>
                          <InputLabel>Typ</InputLabel>
                          <Select
                            value={portFilter.type}
                            onChange={(e) => setPortFilter({...portFilter, type: e.target.value})}
                            label="Typ"
                          >
                            {PORT_FILTER_TYPES.map(type => (
                              <MenuItem key={type.value} value={type.value}>{type.label}</MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                      </Grid>
                      <Grid item xs={12} sm={4}>
                        <FormControl fullWidth>
                          <InputLabel>Port</InputLabel>
                          <Select
                            value={portFilter.value}
                            onChange={(e) => setPortFilter({...portFilter, value: e.target.value})}
                            label="Port"
                          >
                            {COMMON_PORTS.map(port => (
                              <MenuItem key={port.value} value={port.value}>
                                {port.label} ({port.value})
                              </MenuItem>
                            ))}
                            <MenuItem value="custom">
                              <em>Benutzerdefiniert</em>
                            </MenuItem>
                          </Select>
                        </FormControl>
                      </Grid>
                      {portFilter.value === 'custom' && (
                        <Grid item xs={12} sm={2}>
                          <TextField
                            fullWidth
                            label="Port"
                            type="number"
                            placeholder="z.B. 8080"
                            onChange={(e) => setPortFilter({...portFilter, value: e.target.value})}
                          />
                        </Grid>
                      )}
                    </>
                  )}

                  {filterType === 'protocol' && (
                    <Grid item xs={12} sm={6}>
                      <FormControl fullWidth>
                        <InputLabel>Protokoll</InputLabel>
                        <Select
                          value={protocolFilter.type}
                          onChange={(e) => setProtocolFilter({...protocolFilter, type: e.target.value})}
                          label="Protokoll"
                        >
                          {PROTOCOL_FILTER_TYPES.map(type => (
                            <MenuItem key={type.value} value={type.value}>{type.label}</MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                    </Grid>
                  )}

                  {filterType === 'mac' && (
                    <>
                      <Grid item xs={12} sm={3}>
                        <FormControl fullWidth>
                          <InputLabel>Typ</InputLabel>
                          <Select
                            value={macFilter.type}
                            onChange={(e) => setMacFilter({...macFilter, type: e.target.value})}
                            label="Typ"
                          >
                            {MAC_FILTER_TYPES.map(type => (
                              <MenuItem key={type.value} value={type.value}>{type.label}</MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                      </Grid>
                      <Grid item xs={12} sm={4}>
                        <TextField
                          fullWidth
                          label="MAC-Adresse"
                          placeholder="z.B. 00:11:22:33:44:55"
                          value={macFilter.value}
                          onChange={(e) => setMacFilter({...macFilter, value: e.target.value})}
                        />
                      </Grid>
                    </>
                  )}

                  <Grid item xs={12} sm={2}>
                    <Button
                      fullWidth
                      variant="contained"
                      startIcon={<AddIcon />}
                      onClick={handleAddFilter}
                      disabled={
                        (filterType === 'ip' && !ipFilter.value) ||
                        (filterType === 'port' && (!portFilter.value || portFilter.value === 'custom' && !portFilter.customValue)) ||
                        (filterType === 'mac' && !macFilter.value)
                      }
                    >
                      Hinzufügen
                    </Button>
                  </Grid>

                  {activeFilters.length > 0 && (
                    <Grid item xs={12}>
                      <Typography variant="subtitle1">Aktive Filter</Typography>
                      <Paper variant="outlined" sx={{ p: 2, mt: 1 }}>
                        <Stack spacing={1}>
                          {activeFilters.map((filter, index) => (
                            <Box key={filter.id} sx={{ display: 'flex', alignItems: 'center' }}>
                              {index > 0 && (
                                <FormControl size="small" sx={{ width: 100, mr: 1 }}>
                                  <Select
                                    value={filter.logicalOperator || 'and'}
                                    onChange={(e) => {
                                      const updatedFilters = [...activeFilters];
                                      updatedFilters[index].logicalOperator = e.target.value;
                                      setActiveFilters(updatedFilters);
                                    }}
                                    displayEmpty
                                  >
                                    {LOGICAL_OPERATORS.map(op => (
                                      <MenuItem key={op.value} value={op.value}>{op.label}</MenuItem>
                                    ))}
                                  </Select>
                                </FormControl>
                              )}
                              <Chip 
                                label={filter.display} 
                                onDelete={() => handleRemoveFilter(filter.id)}
                                color="primary"
                                variant="outlined"
                                deleteIcon={<DeleteIcon />}
                              />
                            </Box>
                          ))}
                        </Stack>
                      </Paper>
                    </Grid>
                  )}

                  {activeFilters.length > 0 && (
                    <Grid item xs={12}>
                      <Typography variant="subtitle2">Generierter BPF-Filter:</Typography>
                      <Paper variant="outlined" sx={{ p: 1, mt: 1, bgcolor: 'background.default' }}>
                        <code>{generateBpfFilter()}</code>
                      </Paper>
                    </Grid>
                  )}
                </>
              )}

              {/* Aktionsbereich */}
              <Grid item xs={12}>
                <Divider sx={{ my: 1 }} />
                <Stack direction="row" spacing={2} justifyContent="space-between">
                  <Button
                    color="secondary"
                    startIcon={<SaveIcon />}
                    onClick={() => setShowSaveDialog(true)}
                    disabled={(showAdvanced && !manualBpfFilter) || (!showAdvanced && activeFilters.length === 0)}
                  >
                    Filter speichern
                  </Button>
                  <Button 
                    variant="contained" 
                    onClick={handleApplyFilter}
                    disabled={(showAdvanced && !manualBpfFilter) || (!showAdvanced && activeFilters.length === 0)}
                  >
                    Filter anwenden
                  </Button>
                </Stack>
              </Grid>

              {/* Dialog zum Speichern des Filters */}
              {showSaveDialog && (
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Typography variant="subtitle1">Filter speichern</Typography>
                    <TextField
                      fullWidth
                      label="Filtername"
                      value={newFilterName}
                      onChange={(e) => setNewFilterName(e.target.value)}
                      sx={{ mt: 1 }}
                    />
                    <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2 }}>
                      <Button onClick={() => setShowSaveDialog(false)}>Abbrechen</Button>
                      <Button 
                        variant="contained"
                        onClick={handleSaveFilter}
                        disabled={!newFilterName}
                      >
                        Speichern
                      </Button>
                    </Stack>
                  </Paper>
                </Grid>
              )}
            </Grid>
          </CardContent>
        </Card>
      </AccordionDetails>
    </Accordion>
  );
};

export default CaptureFilter; 