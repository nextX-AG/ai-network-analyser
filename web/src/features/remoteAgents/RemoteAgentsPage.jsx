import React from 'react';
import { Container, Typography, Paper } from '@mui/material';
import RemoteAgentsContainer from './containers/RemoteAgentsContainer';

/**
 * RemoteAgentsPage - Hauptseite für die Verwaltung von Remote-Agenten
 * Diese Komponente dient als Container für die Remote-Agents-Funktionalität
 */
const RemoteAgentsPage = () => {
  return (
    <Container maxWidth="lg">
      <Paper elevation={0} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Remote-Agenten
        </Typography>
        <Typography variant="body1" paragraph>
          Hier können Sie die verbundenen Remote-Agenten verwalten, Netzwerkerfassungen starten und Filter anwenden.
        </Typography>
        
        <RemoteAgentsContainer />
      </Paper>
    </Container>
  );
};

export default RemoteAgentsPage; 