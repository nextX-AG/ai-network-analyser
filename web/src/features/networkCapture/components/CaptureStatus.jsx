import React from 'react';
import {
  Box,
  Typography,
  Card,
} from '@mui/material';

/**
 * CaptureStatus - Zeigt den Status der Netzwerkerfassung und die Anzahl der erfassten Pakete an
 */
const CaptureStatus = ({ status, packetCount }) => {
  const statusText = 
    status === 'capturing' ? 'Erfassung läuft' : 
    status === 'error' ? 'Fehler' : 'Bereit';
  
  const statusColor = 
    status === 'capturing' ? 'success.main' : 
    status === 'error' ? 'error.main' : 'text.secondary';
  
  return (
    <Box>
      <Box display="flex" alignItems="center">
        <Typography variant="body1" sx={{ mr: 1 }}>
          Status:
        </Typography>
        <Typography 
          variant="body1" 
          color={statusColor}
          fontWeight="bold"
        >
          {statusText}
        </Typography>
      </Box>
      
      {status === 'capturing' && (
        <Card variant="outlined" sx={{ mt: 1, p: 1, bgcolor: 'action.hover' }}>
          <Box display="flex" justifyContent="space-between" px={2}>
            <Typography variant="body1">
              Erfasste Pakete:
            </Typography>
            <Typography variant="body1" fontWeight="bold">
              {packetCount}
            </Typography>
          </Box>
        </Card>
      )}
    </Box>
  );
};

export default CaptureStatus; 