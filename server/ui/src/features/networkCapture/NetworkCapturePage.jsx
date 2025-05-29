import React, { useState, useEffect } from 'react';
import { Box, Container, Typography, Paper } from '@mui/material';
import NetworkCapturePanel from './components/NetworkCapturePanel';

/**
 * NetworkCapturePage - Hauptseite für die Netzwerkerfassung
 * Diese Komponente dient als Container für die Netzwerkerfassungsfunktionalität
 */
const NetworkCapturePage = () => {
  const [agentUrl, setAgentUrl] = useState('');

  // Lade Agent-URL aus dem lokalen Speicher oder der Konfiguration
  useEffect(() => {
    const storedUrl = localStorage.getItem('agentUrl');
    if (storedUrl) {
      setAgentUrl(storedUrl);
    } else {
      // Fallback auf lokale URL (für Entwicklung)
      setAgentUrl('http://localhost:8090');
    }
  }, []);

  return (
    <Container maxWidth="lg">
      <Paper elevation={0} sx={{ p: 3, mt: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Netzwerkerfassung
        </Typography>
        <Typography variant="body1" paragraph>
          Hier können Sie Netzwerkverkehr erfassen und filtern. Wählen Sie eine Schnittstelle und optional einen Filter, um bestimmte Pakete zu erfassen.
        </Typography>
        
        <Box mt={4}>
          <NetworkCapturePanel agentUrl={agentUrl} />
        </Box>
      </Paper>
    </Container>
  );
};

export default NetworkCapturePage; 