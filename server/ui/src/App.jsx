import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { Box, Container, AppBar, Toolbar, Typography, Tabs, Tab } from '@mui/material';
import NetworkCapturePage from './features/networkCapture/NetworkCapturePage';
import RemoteAgentsPage from './features/remoteAgents/RemoteAgentsPage';
import SettingsPage from './features/settings/SettingsPage';

const App = () => {
  return (
    <Router>
      <Box sx={{ flexGrow: 1 }}>
        <AppBar position="static">
          <Toolbar>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              KI-Netzwerk-Analyzer
            </Typography>
          </Toolbar>
          <Tabs value={false} centered>
            <Tab label="PCAP-Datei hochladen" component={Link} to="/" />
            <Tab label="Live-Capture" component={Link} to="/live" />
            <Tab label="Remote-Agents" component={Link} to="/agents" />
            <Tab label="Einstellungen" component={Link} to="/settings" />
          </Tabs>
        </AppBar>

        <Box sx={{ mt: 3 }}>
          <Routes>
            <Route path="/" element={<NetworkCapturePage />} />
            <Route path="/live" element={<NetworkCapturePage />} />
            <Route path="/agents" element={<RemoteAgentsPage />} />
            <Route path="/settings" element={<SettingsPage />} />
          </Routes>
        </Box>
      </Box>
    </Router>
  );
};

export default App; 