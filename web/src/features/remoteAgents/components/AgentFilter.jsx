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
  IconButton,
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
import {
  IP_FILTER_TYPES,
  PORT_FILTER_TYPES,
  PROTOCOL_FILTER_TYPES,
  MAC_FILTER_TYPES,
  COMMON_PORTS,
  LOGICAL_OPERATORS
} from '../../../shared/constants/filterConstants';
import { convertToBpfSyntax } from '../../../shared/utils/filterUtils';

/**
 * AgentFilter - Filterpanel für einen spezifischen Agenten
 * Basiert auf der NetworkFilterPanel-Komponente
 */
const AgentFilter = ({ onApplyFilter, currentFilter, agentId, onClose }) => {
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
    // Agenten-spezifische Filter
    const loadedFilters = localStorage.getItem(`savedFilters_agent_${agentId}`);
    if (loadedFilters) {
      setSavedFilters(JSON.parse(loadedFilters));
    } else {
      // Fallback auf globale Filter
      const globalFilters = localStorage.getItem('savedNetworkFilters');
      if (globalFilters) {
        setSavedFilters(JSON.parse(globalFilters));
      }
    }
  }, [agentId]);

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
    // Schließe das Panel nach dem Anwenden des Filters
    onClose?.();
  };

  const handleSaveFilter = () => {
    const newFilter = {
      id: Date.now(),
      name: newFilterName,
      filters: showAdvanced ? manualBpfFilter : activeFilters
    };
    
    const updatedFilters = [...savedFilters, newFilter];
    setSavedFilters(updatedFilters);
    
    // Speichere die Filter für diesen spezifischen Agenten
    localStorage.setItem(`savedFilters_agent_${agentId}`, JSON.stringify(updatedFilters));
    
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
    localStorage.setItem(`savedFilters_agent_${agentId}`, JSON.stringify(updatedFilters));
  };

  // Verwende die importierte Funktion zum Konvertieren von Filtern in BPF-Syntax
  const generateBpfFilter = () => {
    return convertToBpfSyntax(activeFilters);
  };

  return (
    <Box>
      <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 2 }}>
        <Typography variant="subtitle1" sx={{ flex: 1 }}>
          <FilterAltIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
          Paketfilter
        </Typography>
        <Button
          size="small"
          variant={!showAdvanced ? "contained" : "outlined"}
          onClick={() => setShowAdvanced(false)}
        >
          Standard
        </Button>
        <Button
          size="small"
          variant={showAdvanced ? "contained" : "outlined"}
          onClick={() => setShowAdvanced(true)}
        >
          BPF
        </Button>
      </Stack>

      {/* Gespeicherte Filter als horizontale Liste */}
      {savedFilters.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Typography variant="caption" display="block" gutterBottom>
            Gespeicherte Filter
          </Typography>
          <Stack direction="row" spacing={1} sx={{ flexWrap: 'wrap', gap: 1 }}>
            {savedFilters.map((filter) => (
              <Chip
                key={filter.id}
                label={filter.name}
                onClick={() => handleLoadFilter(filter)}
                onDelete={() => handleDeleteSavedFilter(filter.id)}
                size="small"
                color="secondary"
                variant="outlined"
              />
            ))}
          </Stack>
        </Box>
      )}

      {/* Fortgeschrittener BPF-Filter */}
      {showAdvanced ? (
        <Box>
          <TextField
            fullWidth
            label="BPF-Syntax Filter"
            placeholder="z.B. tcp port 80 or udp port 53"
            value={manualBpfFilter}
            onChange={(e) => setManualBpfFilter(e.target.value)}
            multiline
            rows={2}
            variant="outlined"
            size="small"
          />
        </Box>
      ) : (
        // Standard-Filter UI bleibt unverändert
        <Grid container spacing={2}>
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
        </Grid>
      )}

      {/* Aktionsbuttons */}
      <Stack direction="row" spacing={2} justifyContent="flex-end" sx={{ mt: 2 }}>
        <Button
          variant="outlined"
          size="small"
          onClick={() => setShowSaveDialog(true)}
          startIcon={<SaveIcon />}
          disabled={showAdvanced ? !manualBpfFilter : activeFilters.length === 0}
        >
          Speichern
        </Button>
        <Button
          variant="contained"
          size="small"
          onClick={handleApplyFilter}
          disabled={showAdvanced ? !manualBpfFilter : activeFilters.length === 0}
        >
          Filter anwenden
        </Button>
      </Stack>

      {/* Dialog zum Speichern von Filtern */}
      {showSaveDialog && (
        <Paper variant="outlined" sx={{ p: 2, mt: 2 }}>
          <Typography variant="subtitle2" gutterBottom>
            Filter speichern
          </Typography>
          <TextField
            fullWidth
            size="small"
            label="Filtername"
            value={newFilterName}
            onChange={(e) => setNewFilterName(e.target.value)}
            sx={{ mb: 2 }}
          />
          <Stack direction="row" spacing={2} justifyContent="flex-end">
            <Button size="small" onClick={() => setShowSaveDialog(false)}>
              Abbrechen
            </Button>
            <Button
              size="small"
              variant="contained"
              onClick={handleSaveFilter}
              disabled={!newFilterName}
            >
              Speichern
            </Button>
          </Stack>
        </Paper>
      )}
    </Box>
  );
};

export default AgentFilter; 