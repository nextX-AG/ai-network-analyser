import React, { useState, useEffect } from 'react';
import { Box, Typography, Alert, CircularProgress } from '@mui/material';
import AgentCard from '../components/AgentCard';
import { fetchAgents, startCapture, stopCapture, setAgentInterface } from '../services/agentApi';

/**
 * RemoteAgentsContainer - Container-Komponente für Remote-Agenten
 * Verwaltet den Zustand und die Logik für Remote-Agenten
 */
const RemoteAgentsContainer = () => {
  const [agents, setAgents] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  
  // Lade Agenten beim Komponenten-Mount
  useEffect(() => {
    const loadAgents = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        const data = await fetchAgents();
        setAgents(data || []);
      } catch (err) {
        console.error('Fehler beim Laden der Agenten:', err);
        setError(err.message);
      } finally {
        setIsLoading(false);
      }
    };
    
    loadAgents();
    
    // Regelmäßiges Polling für Statusupdates
    const intervalId = setInterval(loadAgents, 10000);
    
    return () => clearInterval(intervalId);
  }, []);
  
  // Starte die Paketerfassung auf einem Agenten
  const handleStartCapture = async (agentId, interfaceName, filter) => {
    try {
      const agent = agents.find(a => a.id === agentId);
      if (!agent) return;
      
      const result = await startCapture(agent, interfaceName, filter);
      
      if (result.success) {
        // Agenten-Status aktualisieren
        setAgents(agents.map(a => {
          if (a.id === agentId) {
            return {
              ...a,
              status: 'capturing',
              packets_captured: 0,
              active_filter: filter || null,
              error: '',
            };
          }
          return a;
        }));
      } else {
        // Fehler setzen
        setAgents(agents.map(a => {
          if (a.id === agentId) {
            return {
              ...a,
              status: 'error',
              error: result.error || 'Fehler beim Starten der Erfassung',
            };
          }
          return a;
        }));
      }
    } catch (error) {
      console.error(`Fehler beim Starten der Erfassung für Agent ${agentId}:`, error);
      
      // Fehler setzen
      setAgents(agents.map(a => {
        if (a.id === agentId) {
          return {
            ...a,
            status: 'error',
            error: 'Verbindung zum Agenten fehlgeschlagen',
          };
        }
        return a;
      }));
    }
  };
  
  // Stoppe die Paketerfassung auf einem Agenten
  const handleStopCapture = async (agentId) => {
    try {
      const agent = agents.find(a => a.id === agentId);
      if (!agent) return;
      
      const result = await stopCapture(agent);
      
      if (result.success) {
        // Agenten-Status aktualisieren
        setAgents(agents.map(a => {
          if (a.id === agentId) {
            return {
              ...a,
              status: 'idle',
              active_filter: null,
              error: '',
            };
          }
          return a;
        }));
      } else {
        // Fehler setzen
        setAgents(agents.map(a => {
          if (a.id === agentId) {
            return {
              ...a,
              error: result.error || 'Fehler beim Stoppen der Erfassung',
            };
          }
          return a;
        }));
      }
    } catch (error) {
      console.error(`Fehler beim Stoppen der Erfassung für Agent ${agentId}:`, error);
      
      // Fehler setzen
      setAgents(agents.map(a => {
        if (a.id === agentId) {
          return {
            ...a,
            status: 'error',
            error: 'Verbindung zum Agenten fehlgeschlagen',
          };
        }
        return a;
      }));
    }
  };
  
  // Setze die Schnittstelle für einen Agenten
  const handleSetInterface = async (agentId, interfaceName) => {
    try {
      const agent = agents.find(a => a.id === agentId);
      if (!agent) return;
      
      const result = await setAgentInterface(agent, interfaceName);
      
      if (result.success) {
        // Agenten-Status aktualisieren
        setAgents(agents.map(a => {
          if (a.id === agentId) {
            return {
              ...a,
              interface: interfaceName,
              error: '',
            };
          }
          return a;
        }));
      } else {
        // Fehler setzen
        setAgents(agents.map(a => {
          if (a.id === agentId) {
            return {
              ...a,
              error: result.error || 'Fehler beim Setzen der Schnittstelle',
            };
          }
          return a;
        }));
      }
    } catch (error) {
      console.error(`Fehler beim Setzen der Schnittstelle für Agent ${agentId}:`, error);
      
      // Fehler setzen
      setAgents(agents.map(a => {
        if (a.id === agentId) {
          return {
            ...a,
            status: 'error',
            error: 'Verbindung zum Agenten fehlgeschlagen',
          };
        }
        return a;
      }));
    }
  };
  
  return (
    <>
      {isLoading && agents.length === 0 && (
        <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
          <CircularProgress />
        </Box>
      )}
      
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}
      
      {agents.length === 0 && !isLoading && (
        <Alert severity="info" sx={{ mb: 3 }}>
          Keine Remote-Agenten verfügbar. Bitte stellen Sie sicher, dass Agenten registriert und verbunden sind.
        </Alert>
      )}
      
      {agents.map((agent) => (
        <AgentCard
          key={agent.id}
          agent={agent}
          onStartCapture={handleStartCapture}
          onStopCapture={handleStopCapture}
          onSetInterface={handleSetInterface}
        />
      ))}
    </>
  );
};

export default RemoteAgentsContainer; 