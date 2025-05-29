import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Select,
  Typography,
  Alert,
} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import StopIcon from '@mui/icons-material/Stop';
import CaptureFilter from './CaptureFilter';
import CaptureStatus from './CaptureStatus';
import { 
  fetchAgentStatus, 
  startCapture as apiStartCapture, 
  stopCapture as apiStopCapture,
  setInterface as apiSetInterface 
} from '../services/captureApi';

/**
 * NetworkCapturePanel - Hauptkomponente für die Netzwerkerfassung
 * Ermöglicht das Starten und Stoppen der Paketerfassung sowie die Filterung
 */
const NetworkCapturePanel = ({ agentUrl }) => {
  const [interfaces, setInterfaces] = useState([]);
  const [selectedInterface, setSelectedInterface] = useState('');
  const [captureStatus, setCaptureStatus] = useState('idle'); // 'idle', 'capturing', 'error'
  const [errorMessage, setErrorMessage] = useState('');
  const [activeFilter, setActiveFilter] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [packetCount, setPacketCount] = useState(0);
  
  // Lade verfügbare Netzwerkschnittstellen vom Agent
  useEffect(() => {
    if (!agentUrl) return;
    
    const fetchStatus = async () => {
      try {
        setIsLoading(true);
        const agentStatus = await fetchAgentStatus(agentUrl);
        
        setSelectedInterface(agentStatus.interface || '');
        setCaptureStatus(agentStatus.status);
        setPacketCount(agentStatus.packets_captured || 0);
        
        if (agentStatus.interfaces) {
          setInterfaces(agentStatus.interfaces);
        }
      } catch (error) {
        console.error('Fehler beim Laden der Schnittstellen:', error);
        setErrorMessage('Die Schnittstellen konnten nicht vom Agenten geladen werden.');
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchStatus();
    
    // Regelmäßiges Polling für Statusupdates
    const intervalId = setInterval(fetchStatus, 3000);
    
    return () => clearInterval(intervalId);
  }, [agentUrl]);
  
  // Starte die Paketerfassung
  const handleStartCapture = async () => {
    if (!agentUrl || !selectedInterface) return;
    
    try {
      setIsLoading(true);
      
      const result = await apiStartCapture(agentUrl, selectedInterface, activeFilter);
      
      if (result.success) {
        setCaptureStatus('capturing');
        setErrorMessage('');
      } else {
        setCaptureStatus('error');
        setErrorMessage(result.error || 'Fehler beim Starten der Erfassung');
      }
    } catch (error) {
      console.error('Fehler beim Starten der Erfassung:', error);
      setCaptureStatus('error');
      setErrorMessage('Die Paketerfassung konnte nicht gestartet werden.');
    } finally {
      setIsLoading(false);
    }
  };
  
  // Stoppe die Paketerfassung
  const handleStopCapture = async () => {
    if (!agentUrl) return;
    
    try {
      setIsLoading(true);
      const result = await apiStopCapture(agentUrl);
      
      if (result.success) {
        setCaptureStatus('idle');
        setErrorMessage('');
      } else {
        setErrorMessage(result.error || 'Fehler beim Stoppen der Erfassung');
      }
    } catch (error) {
      console.error('Fehler beim Stoppen der Erfassung:', error);
      setErrorMessage('Die Paketerfassung konnte nicht gestoppt werden.');
    } finally {
      setIsLoading(false);
    }
  };
  
  // Setze die Netzwerkschnittstelle
  const handleSetInterface = async (interfaceName) => {
    if (!agentUrl) return;
    
    try {
      setIsLoading(true);
      const result = await apiSetInterface(agentUrl, interfaceName);
      
      if (result.success) {
        setSelectedInterface(interfaceName);
        setErrorMessage('');
      } else {
        setErrorMessage(result.error || 'Fehler beim Setzen der Schnittstelle');
      }
    } catch (error) {
      console.error('Fehler beim Setzen der Schnittstelle:', error);
      setErrorMessage('Die Netzwerkschnittstelle konnte nicht gesetzt werden.');
    } finally {
      setIsLoading(false);
    }
  };
  
  // Anwenden des Filters
  const handleApplyFilter = (filter) => {
    setActiveFilter(filter);
  };
  
  return (
    <Box>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <CaptureFilter 
            onApplyFilter={handleApplyFilter} 
            currentFilter={activeFilter} 
          />
        </Grid>
        
        <Grid item xs={12}>
          <Card variant="outlined">
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Netzwerkerfassung
              </Typography>
              
              {errorMessage && (
                <Alert severity="error" sx={{ mb: 2 }}>
                  {errorMessage}
                </Alert>
              )}
              
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} md={6}>
                  <FormControl fullWidth disabled={captureStatus === 'capturing' || isLoading}>
                    <InputLabel>Netzwerkschnittstelle</InputLabel>
                    <Select
                      value={selectedInterface}
                      onChange={(e) => handleSetInterface(e.target.value)}
                      label="Netzwerkschnittstelle"
                    >
                      {interfaces.map((iface) => (
                        <MenuItem key={iface.name} value={iface.name}>
                          {iface.name} {iface.ips && iface.ips.length > 0 && `(${iface.ips[0]})`}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>
                
                <Grid item xs={12} md={3}>
                  <CaptureStatus 
                    status={captureStatus} 
                    packetCount={packetCount} 
                  />
                </Grid>
                
                <Grid item xs={12} md={3}>
                  {captureStatus === 'capturing' ? (
                    <Button
                      fullWidth
                      variant="contained"
                      color="error"
                      startIcon={<StopIcon />}
                      onClick={handleStopCapture}
                      disabled={isLoading}
                    >
                      {isLoading ? <CircularProgress size={24} /> : 'Erfassung stoppen'}
                    </Button>
                  ) : (
                    <Button
                      fullWidth
                      variant="contained"
                      color="primary"
                      startIcon={<PlayArrowIcon />}
                      onClick={handleStartCapture}
                      disabled={!selectedInterface || isLoading}
                    >
                      {isLoading ? <CircularProgress size={24} /> : 'Erfassung starten'}
                    </Button>
                  )}
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default NetworkCapturePanel; 