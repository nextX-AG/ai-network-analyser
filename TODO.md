# KI-Netzwerk-Analyzer TODO List

This document provides a high-level overview of the current tasks for the KI-Netzwerk-Analyzer project. For more detailed tasks and progress tracking, please see the [docs/TODO.md](docs/TODO.md).

## Current Focus Areas

### High Priority
- [x] Fix Remote Agent interface persistence issue - after restart selected interface is not preserved
- [x] Fix Remote Agent UI interface display - active interface is not shown in status
- [ ] Implement server-side network interface selection for agents
- [ ] Docker configuration for development environment
- [ ] SQLite integration for data persistence
- [ ] Implement optimizations for large PCAP files
- [ ] Complete the authentication and security concept for remote agents
- [ ] Implement Speech2Text module with Whisper.cpp integration

### Medium Priority
- [ ] Timeline visualization with Three.js
- [ ] AI annotation module with OpenAI GPT API integration
- [ ] Implement event and timeline module
- [ ] Extend test coverage for critical components
- [ ] Create Docker images for remote agents

### Low Priority
- [ ] Multi-agent capture synchronization
- [ ] Mobile view for frontend
- [ ] Prepare for AQEA compatibility
- [ ] CI/CD pipeline for automated tests and builds
- [ ] Cloud-based deployment option

## Completed Tasks
- [x] Project structure definition and initialization
- [x] Basic Go backend implementation
- [x] Integration of packet capture with gopacket
- [x] Basic React/Three.js frontend scaffold
- [x] Remote capture system for distributed capture
- [x] Web interface for agents
- [x] Automatic detection and registration of agents
- [x] Bridge optimization for MITM monitoring
- [x] Gateway detection and analysis implementation
- [x] REST API endpoints for gateway information
- [x] Systemd service templates for easy deployment
- [x] Fix Admin UI Route registration - made Admin UI accessible
- [x] Fix configuration file permissions - resolved read-only filesystem issue for configuration
- [x] Fix Remote Agent interface persistence and status display
- [x] Fix packet capturing permissions - ensured agent runs as root with proper capabilities

## Current Remote Agent Improvements
- [x] Added multiple configuration paths to handle read-only filesystems
- [x] Improved error handling in configuration saving
- [x] Updated installation script to use writable configuration paths
- [x] Added permission checks and fixes for configuration files
- [x] Fix interface persistence after agent restart
- [x] Fix active interface display in status UI
- [x] Added UpdateInterface method to PcapCapturer to ensure configuration is updated
- [x] Improved saveConfig function to try multiple paths if one fails
- [x] Ensured restart handler saves configuration before restarting
- [x] Added root permission check and explicit capability requirements in agent
- [x] Enhanced systemd service to ensure proper network capture permissions

## Server-Side Network Interface Selection Implementation
- [ ] Fix agent registration to use actual routable IP address instead of 0.0.0.0
- [ ] Enhance agent to send complete list of available network interfaces with details (IP, MAC, bridge status)
- [ ] Implement server-side API endpoint to select and activate interfaces on agents
- [ ] Update UI to display all available interfaces for each agent
- [ ] Add interface selection controls in Remote-Agents UI
- [ ] Implement WebSocket protocol for real-time capture status updates
- [ ] Add error handling for unreachable interfaces

## Current Agent Issues to Fix
- [x] Fix packet counter display in UI when packets are captured (Agent shows captured packets but UI doesn't)
- [x] Fix heartbeat mechanism to include captured packet count in status updates
- [x] Implement workaround for UI updating with real-time packet counts via polling
- [x] Ensure interface configuration is correctly persisted between agent restarts
- [x] Fix Server-URL configuration persistence and prioritization of saved values
- [x] Add detailed logging for agent configuration saving/loading process
- [x] Fix CORS issues with Agent API to allow cross-origin access from main server UI
- [ ] Implement proper error handling for WebSocket communication failures
- [ ] Add server-side packet counter validation against agent-reported values

## Next Action Items

1. Complete server-side network interface selection
2. Complete SQLite integration for data persistence
3. Implement the Speech2Text module
4. Begin development of the Three.js timeline visualization
5. Set up Docker configuration for development
6. Begin AI integration for packet analysis 