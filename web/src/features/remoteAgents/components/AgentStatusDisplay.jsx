import React from 'react';
import { Box, Typography } from '@mui/material';

/**
 * AgentStatusDisplay - Zeigt den Status eines Agenten an
 */
const AgentStatusDisplay = ({ status, packetCount, isCapturing, hasError }) => {
  return (
    <Box>
      <Typography variant="body2" color="text.secondary">
        Status:
      </Typography>
      <Typography 
        variant="body1" 
        color={
          isCapturing ? 'success.main' : (hasError ? 'error.main' : 'text.secondary')
        }
        fontWeight="bold"
      >
        {isCapturing ? 'Erfassung läuft' : (hasError ? 'Fehler' : 'Bereit')}
      </Typography>
      
      {isCapturing && packetCount !== undefined && (
        <Box sx={{ mt: 1 }}>
          <Typography variant="body2" color="text.secondary">
            Erfasste Pakete:
          </Typography>
          <Typography variant="body1" fontWeight="bold">
            {packetCount || 0}
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default AgentStatusDisplay; 