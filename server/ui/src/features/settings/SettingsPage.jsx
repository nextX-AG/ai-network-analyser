import React from 'react';
import {
  Container,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
  Grid,
  TextField,
  Button,
  Divider
} from '@mui/material';

/**
 * SettingsPage - Hauptseite für die Anwendungseinstellungen
 * Hier können globale Einstellungen wie API-Keys und Systemkonfigurationen verwaltet werden
 */
const SettingsPage = () => {
  return (
    <Container maxWidth="lg">
      <Paper elevation={0} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Einstellungen
        </Typography>
        <Typography variant="body1" paragraph>
          Hier können Sie globale Einstellungen der Anwendung verwalten.
        </Typography>

        <Grid container spacing={3}>
          {/* AI Konfiguration */}
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  KI-Konfiguration
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Konfigurieren Sie hier Ihre KI-API-Zugangsdaten und Einstellungen.
                </Typography>
                <Box sx={{ mt: 2 }}>
                  <TextField
                    fullWidth
                    label="OpenAI API Key"
                    type="password"
                    placeholder="sk-..."
                    sx={{ mb: 2 }}
                  />
                  <TextField
                    fullWidth
                    label="Model"
                    defaultValue="gpt-4"
                    sx={{ mb: 2 }}
                  />
                  <Button variant="contained" color="primary">
                    Speichern
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* System Konfiguration */}
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  System-Konfiguration
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Allgemeine Systemeinstellungen und Konfigurationen.
                </Typography>
                <Box sx={{ mt: 2 }}>
                  <TextField
                    fullWidth
                    label="Standard-Interface"
                    placeholder="eth0"
                    sx={{ mb: 2 }}
                  />
                  <TextField
                    fullWidth
                    label="Capture Buffer Size (MB)"
                    type="number"
                    defaultValue={100}
                    sx={{ mb: 2 }}
                  />
                  <Button variant="contained" color="primary">
                    Speichern
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* Logging und Debug */}
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Logging & Debug
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Konfigurieren Sie die Logging- und Debug-Einstellungen.
                </Typography>
                <Box sx={{ mt: 2 }}>
                  <TextField
                    fullWidth
                    label="Log Level"
                    select
                    defaultValue="info"
                    SelectProps={{
                      native: true,
                    }}
                    sx={{ mb: 2 }}
                  >
                    <option value="debug">Debug</option>
                    <option value="info">Info</option>
                    <option value="warn">Warning</option>
                    <option value="error">Error</option>
                  </TextField>
                  <TextField
                    fullWidth
                    label="Log Datei"
                    defaultValue="/var/log/ai-network-analyser.log"
                    sx={{ mb: 2 }}
                  />
                  <Button variant="contained" color="primary">
                    Speichern
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Paper>
    </Container>
  );
};

export default SettingsPage; 