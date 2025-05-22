import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardActions,
  CardHeader,
  Collapse,
  Chip,
  Divider,
  IconButton,
  Typography,
  CircularProgress,
  Grid,
  Badge,
  Alert,
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import StopIcon from '@mui/icons-material/Stop';
import FilterAltIcon from '@mui/icons-material/FilterAlt';
import AgentFilter from './AgentFilter';
import AgentStatusDisplay from './AgentStatusDisplay';

/**
 * AgentCard - Komponente zur Darstellung und Verwaltung eines einzelnen Agenten
 */
const AgentCard = ({ agent, onStartCapture, onStopCapture, onSetInterface }) => {
  const [expanded, setExpanded] = useState(false);
  const [filterExpanded, setFilterExpanded] = useState(false);
  const [activeFilter, setActiveFilter] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  
  // Status des Agenten
  const isCapturing = agent.status === 'capturing';
  const hasError = agent.status === 'error';
  const isConnected = agent.status !== 'disconnected';
  
  // Filter-Anzeige basierend auf dem Agentenstatus
  useEffect(() => {
    if (agent.active_filter) {
      setActiveFilter(agent.active_filter);
    }
  }, [agent.active_filter]);
  
  // Handler für den Filter
  const handleApplyFilter = (filter) => {
    setActiveFilter(filter);
    
    // Die Hauptanwendung soll den Filter anwenden
    if (isCapturing) {
      // Wenn der Agent bereits erfasst, muss die Erfassung neu gestartet werden
      onStopCapture(agent.id);
      // Kurz warten und dann mit dem neuen Filter starten
      setTimeout(() => {
        onStartCapture(agent.id, agent.interface, filter);
      }, 500);
    }
  };
  
  // Anzeige für aktiven Filter
  const renderActiveFilter = () => {
    if (!activeFilter) return null;
    
    return (
      <Box sx={{ mt: 1 }}>
        <Chip
          icon={<FilterAltIcon />}
          label={typeof activeFilter === 'string' 
            ? `Filter: ${activeFilter.length > 30 ? activeFilter.substring(0, 30) + '...' : activeFilter}` 
            : `${activeFilter.length || 0} Filter aktiv`}
          color="primary"
          variant="outlined"
          size="small"
        />
      </Box>
    );
  };
  
  return (
    <Card variant="outlined" sx={{ mb: 2 }}>
      <CardHeader
        title={
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Typography variant="h6" component="div">
              {agent.name}
            </Typography>
            {activeFilter && (
              <Badge 
                badgeContent="" 
                color="primary" 
                variant="dot" 
                sx={{ ml: 1 }}
              />
            )}
          </Box>
        }
        subheader={
          <Box>
            <Typography variant="body2" color="text.secondary">
              {agent.url} • {agent.interface || 'Keine Schnittstelle ausgewählt'}
            </Typography>
            {renderActiveFilter()}
          </Box>
        }
        action={
          <Box sx={{ display: 'flex' }}>
            <IconButton
              onClick={() => setFilterExpanded(!filterExpanded)}
              aria-expanded={filterExpanded}
              aria-label="Filter anzeigen"
              size="small"
              color={activeFilter ? "primary" : "default"}
            >
              <FilterAltIcon />
            </IconButton>
            <IconButton
              onClick={() => setExpanded(!expanded)}
              aria-expanded={expanded}
              aria-label="Details anzeigen"
              size="small"
            >
              <ExpandMoreIcon />
            </IconButton>
          </Box>
        }
      />
      
      {/* Filter-Erweiterung */}
      <Collapse in={filterExpanded} timeout="auto" unmountOnExit>
        <CardContent sx={{ pt: 0 }}>
          <AgentFilter 
            onApplyFilter={handleApplyFilter} 
            currentFilter={activeFilter}
            agentId={agent.id}
          />
        </CardContent>
      </Collapse>
      
      {/* Status und Aktionen */}
      <CardContent sx={{ pt: 0 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <AgentStatusDisplay 
            status={agent.status} 
            packetCount={agent.packets_captured} 
            isCapturing={isCapturing}
            hasError={hasError}
          />
          
          <Box>
            {isCapturing ? (
              <Button
                variant="contained"
                color="error"
                startIcon={<StopIcon />}
                onClick={() => onStopCapture(agent.id)}
                disabled={isLoading || !isConnected}
                size="small"
              >
                {isLoading ? <CircularProgress size={24} /> : 'Stoppen'}
              </Button>
            ) : (
              <Button
                variant="contained"
                color="primary"
                startIcon={<PlayArrowIcon />}
                onClick={() => onStartCapture(agent.id, agent.interface, activeFilter)}
                disabled={isLoading || !isConnected || !agent.interface}
                size="small"
              >
                {isLoading ? <CircularProgress size={24} /> : 'Starten'}
              </Button>
            )}
          </Box>
        </Box>
        
        {hasError && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {agent.error || 'Ein Fehler ist aufgetreten.'}
          </Alert>
        )}
      </CardContent>
      
      {/* Erweiterte Details */}
      <Collapse in={expanded} timeout="auto" unmountOnExit>
        <CardContent sx={{ pt: 0 }}>
          <Divider sx={{ my: 2 }} />
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="text.secondary">
                Hostname:
              </Typography>
              <Typography variant="body1">
                {agent.hostname || 'Unbekannt'}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="text.secondary">
                Betriebssystem:
              </Typography>
              <Typography variant="body1">
                {agent.os || 'Unbekannt'}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="text.secondary">
                Uptime:
              </Typography>
              <Typography variant="body1">
                {agent.uptime || 'Unbekannt'}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="text.secondary">
                Letzter Heartbeat:
              </Typography>
              <Typography variant="body1">
                {agent.last_heartbeat 
                  ? new Date(agent.last_heartbeat).toLocaleString() 
                  : 'Unbekannt'}
              </Typography>
            </Grid>
            
            {/* Weitere Informationen */}
            <Grid item xs={12}>
              <Typography variant="body2" color="text.secondary">
                Verfügbare Schnittstellen:
              </Typography>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 1 }}>
                {agent.interfaces && agent.interfaces.map((iface) => (
                  <Chip
                    key={iface}
                    label={iface}
                    variant={agent.interface === iface ? 'filled' : 'outlined'}
                    color={agent.interface === iface ? 'primary' : 'default'}
                    onClick={() => onSetInterface(agent.id, iface)}
                    size="small"
                  />
                ))}
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Collapse>
    </Card>
  );
};

export default AgentCard; 