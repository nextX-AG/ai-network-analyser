<!-- Version: 0.1.0 | Last Updated: 2024-06-19 14:30:00 UTC -->


# Deployment und CI/CD-Prozess: KI-Netzwerk-Analyzer

## Übersicht

Dieses Dokument beschreibt den Entwicklungs-, Build-, Test- und Deployment-Prozess für den KI-Netzwerk-Analyzer. Es definiert die Branching-Strategie, CI/CD-Pipeline und Best Practices für die Bereitstellung der Anwendung.

## Branching-Strategie und Git-Workflow

### Branching-Modell

Wir verwenden GitHub Flow, eine vereinfachte Branching-Strategie, die sich gut für kontinuierliche Deployment-Umgebungen eignet:

```
main ──────●───────●───────●───────●───────●───────●───────●───────
            │       │       │       │       │       │
            │       │       │       │       │       │
feature/A ──┘       │       │       │       │       │
                    │       │       │       │       │
                    │       │       │       │       │
feature/B ──────────┘       │       │       │       │
                            │       │       │       │
                            │       │       │       │
fix/C ─────────────────────┘       │       │       │
                                    │       │       │
                                    │       │       │
feature/D ────────────────────────┘       │       │
                                          │       │
                                          │       │
docs/E ────────────────────────────────────┘       │
                                                    │
                                                    │
refactor/F ──────────────────────────────────────────┘
```

### Branch-Typen

- **main**: Produktionsreifer Code, immer stabil und bereit für Deployment
- **feature/\***: Entwicklung neuer Funktionalitäten
- **fix/\***: Behebung von Fehlern
- **docs/\***: Dokumentationsänderungen
- **refactor/\***: Code-Refactoring ohne Änderung der Funktionalität
- **perf/\***: Performance-Optimierungen

### Commit-Konventionen

Wir verwenden das [Conventional Commits](https://www.conventionalcommits.org/)-Format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Beispiele:
- `feat(timeline): Zoom-Funktionalität für die Timeline-Komponente`
- `fix(api): Fehlerbehandlung bei fehlerhaften Paketen verbessert`
- `docs(setup): Installationsanleitung für macOS hinzugefügt`
- `refactor(packet-parser): Code-Struktur für bessere Lesbarkeit optimiert`

### Pull Request-Prozess

1. Feature-/Fix-Branches werden von `main` abgezweigt
2. Nach Fertigstellung wird ein Pull Request erstellt
3. Code-Review durch mindestens einen anderen Entwickler
4. Automatisierte Tests müssen erfolgreich sein
5. Nach Genehmigung erfolgt der Merge in `main` (Squash oder Rebase)

## CI/CD-Pipeline

### Kontinuierliche Integration

Bei jedem Push und Pull Request wird folgende Pipeline ausgeführt:

```
┌───────────┐     ┌───────────┐     ┌──────────┐     ┌──────────────┐     ┌─────────────┐
│ Code-Push │────>│ Go Linting│────>│ Go Tests │────>│ Frontend Tests│────>│ Build & Tag │
└───────────┘     └───────────┘     └──────────┘     └──────────────┘     └─────────────┘
```

### Build-Prozess

#### Go-Backend-Build

```
┌──────────────┐     ┌───────────────┐     ┌─────────────────┐
│ Abhängigkeiten│────>│ Go Build mit  │────>│ Statische Analyse│
│ installieren  │     │ Version-Flags │     │ (go vet, etc.)   │
└──────────────┘     └───────────────┘     └─────────────────┘
        │                    │                      │
        └────────────────────┼──────────────────────┘
                             ▼
                     ┌───────────────┐
                     │ Backend-Binary │
                     └───────────────┘
```

#### Frontend-Build

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ NPM-Pakete   │────>│ TypeScript   │────>│ Webpack Build│────>│ Optimierung  │
│ installieren │     │ Kompilierung │     │ & Bundling   │     │ & Minifizierung│
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
                                                                     │
                                                                     ▼
                                                          ┌────────────────────┐
                                                          │ Frontend-Distribution│
                                                          └────────────────────┘
```

### Kontinuierliches Deployment

```
┌───────────┐     ┌────────────────┐     ┌──────────────┐     ┌──────────────┐
│ Entwickler│────>│ Pull Request   │────>│ Automatische │────>│ Release mit  │
│ Umgebung  │     │ & Code Review  │     │ Tests        │     │ GitHub Tags  │
└───────────┘     └────────────────┘     └──────────────┘     └──────────────┘
   Kontinuierlich    Nach jedem Feature    Bei jedem PR       Bei Versionserhöhung
```

## Build-Umgebungen

### Entwicklungsumgebung

- **Zweck**: Lokale Entwicklung und Tests
- **Aktualisierung**: Kontinuierlich (lokale Builds)
- **Konfiguration**: Development-spezifisch mit Hot-Reloading
- **Besonderheiten**: 
  - Backend: `go run cmd/server/main.go --config dev.json`
  - Frontend: `npm run dev` oder `yarn dev`

### Test-Umgebung (CI)

- **Zweck**: Automatisierte Tests für Pull Requests
- **Aktualisierung**: Bei jedem Push in PR-Branches
- **Konfiguration**: Test-spezifisch mit Mocks für externe Dienste
- **Besonderheiten**:
  - Vollständiger Test-Suite-Lauf
  - Code-Coverage-Berichte
  - Statische Code-Analyse

### Release-Umgebung

- **Zweck**: Bereitstellung von Releases
- **Aktualisierung**: Bei Tags (v0.1.0, v0.2.0, etc.)
- **Konfiguration**: Produktionsähnlich
- **Besonderheiten**:
  - Cross-Platform-Builds (Linux, macOS, Windows)
  - Signierte Binärdateien
  - Release-Artefakte (Binaries, Installationspakete)

## Build- und Deployment-Prozess

### Backend-Build (Go)

1. **Abhängigkeiten installieren**:
   ```bash
   go mod download
   ```

2. **Version einbinden**:
   ```bash
   VERSION=$(cat VERSION.txt)
   GIT_COMMIT=$(git rev-parse HEAD)
   BUILD_DATE=$(date -u '+%Y-%m-%d %H:%M:%S')
   ```

3. **Binary kompilieren**:
   ```bash
   go build -ldflags "-X github.com/username/ki-network-analyzer/pkg/version.Version=$VERSION -X github.com/username/ki-network-analyzer/pkg/version.CommitHash=$GIT_COMMIT -X github.com/username/ki-network-analyzer/pkg/version.BuildDate=$BUILD_DATE" -o bin/ki-network-analyzer cmd/server/main.go
   ```

4. **Cross-Compilation (optional)**:
   ```bash
   # Für Windows
   GOOS=windows GOARCH=amd64 go build -ldflags "..." -o bin/ki-network-analyzer.exe cmd/server/main.go
   
   # Für macOS
   GOOS=darwin GOARCH=amd64 go build -ldflags "..." -o bin/ki-network-analyzer-macos cmd/server/main.go
   ```

### Frontend-Build (React/TypeScript)

1. **Abhängigkeiten installieren**:
   ```bash
   cd web && npm install
   ```

2. **Produktionsbuild erstellen**:
   ```bash
   npm run build
   ```

3. **Statische Assets für Backend vorbereiten**:
   ```bash
   mkdir -p ../bin/web
   cp -r build/* ../bin/web/
   ```

### Deployment-Schritte

1. **Pre-Deployment-Checks**:
   - Überprüfung der Zielplattform-Kompatibilität
   - Validierung der Konfigurationsdateien
   - Sicherheitsüberprüfungen für API-Schlüssel und Secrets

2. **Deployment-Prozess**:
   - Binäry und statische Assets kopieren
   - Konfigurationen anpassen
   - Dienst starten oder neu starten

3. **Post-Deployment-Validierung**:
   - Healthcheck-Endpunkte überprüfen
   - Funktionale Tests durchführen
   - Performance-Monitoring initialisieren

### Paketierungsoptionen

1. **Standalone-Binary mit eingebetteten Assets**:
   - Go-Binary mit eingebetteten Frontend-Assets via `embed` Package
   - Einfache distribution.zip mit Binary und Konfigurationsvorlagen

2. **Docker-Container** (empfohlen):
   - Multi-Stage-Build für minimale Image-Größe
   - Alpine-basiertes Image für Produktionsumgebungen

3. **Installationspakete** (optional):
   - .deb-Pakete für Debian/Ubuntu
   - .rpm-Pakete für RHEL/Fedora/CentOS
   - macOS-Installer-Pakete

## Docker-basiertes Deployment

### Dockerfile

```dockerfile
# Build stage for backend
FROM golang:1.20-alpine AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION
ARG COMMIT_HASH
ARG BUILD_DATE

RUN go build -ldflags "-X github.com/username/ki-network-analyzer/pkg/version.Version=${VERSION} -X github.com/username/ki-network-analyzer/pkg/version.CommitHash=${COMMIT_HASH} -X github.com/username/ki-network-analyzer/pkg/version.BuildDate=${BUILD_DATE}" -o /ki-network-analyzer cmd/server/main.go

# Build stage for frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build

# Final stage
FROM alpine:3.17
WORKDIR /app

RUN apk --no-cache add ca-certificates libcap && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Copy binaries and assets
COPY --from=backend-builder /ki-network-analyzer .
COPY --from=frontend-builder /app/build ./web
COPY configs/config.docker.json ./config.json

# Set capabilities for network capturing (if needed)
RUN setcap cap_net_raw,cap_net_admin=eip /app/ki-network-analyzer

USER appuser
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q --spider http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/app/ki-network-analyzer", "--config", "/app/config.json"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  ki-network-analyzer:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        - VERSION=${VERSION:-0.1.0}
        - COMMIT_HASH=${COMMIT_HASH:-unknown}
        - BUILD_DATE=${BUILD_DATE:-unknown}
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./pcaps:/app/pcaps
      - ./config.json:/app/config.json
    restart: unless-stopped
    cap_add:
      - NET_ADMIN
      - NET_RAW
    networks:
      - ki-network

networks:
  ki-network:
    driver: bridge
```

## Systemd Service (Linux)

Für Linux-basierte Deployments ohne Docker kann ein systemd-Service verwendet werden:

```ini
[Unit]
Description=KI-Netzwerk-Analyzer - Intelligente Netzwerkverkehr-Analyse
After=network.target

[Service]
Type=simple
User=analyzer
Group=analyzer
WorkingDirectory=/opt/ki-network-analyzer
ExecStart=/opt/ki-network-analyzer/ki-network-analyzer --config /etc/ki-network-analyzer/config.json
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=ki-network-analyzer
Environment=PATH=/usr/bin:/bin:/usr/local/bin

# Security hardening
PrivateTmp=true
ProtectSystem=full
ReadWritePaths=/var/lib/ki-network-analyzer /var/log/ki-network-analyzer
NoNewPrivileges=true
ProtectHome=true
ProtectControlGroups=true
ProtectKernelModules=true
ProtectKernelTunables=true
CapabilityBoundingSet=CAP_NET_RAW CAP_NET_ADMIN
AmbientCapabilities=CAP_NET_RAW CAP_NET_ADMIN

[Install]
WantedBy=multi-user.target
```

## CI/CD mit GitHub Actions

### Workflow für Tests

```yaml
# .github/workflows/test.yml
name: Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.txt

  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: cd web && npm install
      
      - name: Run tests
        run: cd web && npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./web/coverage/lcov.info
```

### Workflow für Release-Builds

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build-and-release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Build frontend
        run: |
          cd web
          npm install
          npm run build
      
      - name: Set version variables
        id: vars
        run: |
          VERSION=${GITHUB_REF#refs/tags/v}
          COMMIT_HASH=$(git rev-parse HEAD)
          BUILD_DATE=$(date -u '+%Y-%m-%d %H:%M:%S')
          echo "VERSION=$VERSION" >> $GITHUB_ENV
          echo "COMMIT_HASH=$COMMIT_HASH" >> $GITHUB_ENV
          echo "BUILD_DATE=$BUILD_DATE" >> $GITHUB_ENV
      
      - name: Build binaries
        run: |
          mkdir -p bin/linux bin/macos bin/windows bin/web
          
          # Copy frontend build
          cp -r web/build/* bin/web/
          
          # Build for Linux
          GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/username/ki-network-analyzer/pkg/version.Version=${VERSION} -X github.com/username/ki-network-analyzer/pkg/version.CommitHash=${COMMIT_HASH} -X github.com/username/ki-network-analyzer/pkg/version.BuildDate=${BUILD_DATE}" -o bin/linux/ki-network-analyzer cmd/server/main.go
          
          # Build for macOS
          GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/username/ki-network-analyzer/pkg/version.Version=${VERSION} -X github.com/username/ki-network-analyzer/pkg/version.CommitHash=${COMMIT_HASH} -X github.com/username/ki-network-analyzer/pkg/version.BuildDate=${BUILD_DATE}" -o bin/macos/ki-network-analyzer cmd/server/main.go
          
          # Build for Windows
          GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/username/ki-network-analyzer/pkg/version.Version=${VERSION} -X github.com/username/ki-network-analyzer/pkg/version.CommitHash=${COMMIT_HASH} -X github.com/username/ki-network-analyzer/pkg/version.BuildDate=${BUILD_DATE}" -o bin/windows/ki-network-analyzer.exe cmd/server/main.go
      
      - name: Create distribution packages
        run: |
          # Create Linux package
          mkdir -p dist/ki-network-analyzer-${VERSION}-linux
          cp bin/linux/ki-network-analyzer dist/ki-network-analyzer-${VERSION}-linux/
          cp -r bin/web dist/ki-network-analyzer-${VERSION}-linux/
          cp configs/config.example.json dist/ki-network-analyzer-${VERSION}-linux/config.json
          tar -czf ki-network-analyzer-${VERSION}-linux.tar.gz -C dist ki-network-analyzer-${VERSION}-linux
          
          # Create macOS package
          mkdir -p dist/ki-network-analyzer-${VERSION}-macos
          cp bin/macos/ki-network-analyzer dist/ki-network-analyzer-${VERSION}-macos/
          cp -r bin/web dist/ki-network-analyzer-${VERSION}-macos/
          cp configs/config.example.json dist/ki-network-analyzer-${VERSION}-macos/config.json
          tar -czf ki-network-analyzer-${VERSION}-macos.tar.gz -C dist ki-network-analyzer-${VERSION}-macos
          
          # Create Windows package
          mkdir -p dist/ki-network-analyzer-${VERSION}-windows
          cp bin/windows/ki-network-analyzer.exe dist/ki-network-analyzer-${VERSION}-windows/
          cp -r bin/web dist/ki-network-analyzer-${VERSION}-windows/
          cp configs/config.example.json dist/ki-network-analyzer-${VERSION}-windows/config.json
          zip -r ki-network-analyzer-${VERSION}-windows.zip dist/ki-network-analyzer-${VERSION}-windows
      
      - name: Create Release
        id: create_release
        uses: softprops/action-gh-release@v1
        with:
          tag_name: ${{ github.ref }}
          name: Release ${{ env.VERSION }}
          draft: false
          prerelease: false
          files: |
            ki-network-analyzer-${{ env.VERSION }}-linux.tar.gz
            ki-network-analyzer-${{ env.VERSION }}-macos.tar.gz
            ki-network-analyzer-${{ env.VERSION }}-windows.zip
```

## Monitoring und Betrieb

### Healthchecks

Die Anwendung stellt einen `/api/health`-Endpunkt zur Verfügung, der den Status der Komponenten überwacht:

```json
{
  "status": "healthy",
  "version": "0.1.0",
  "components": {
    "database": "healthy",
    "ai_client": "healthy",
    "speech_service": "healthy"
  },
  "uptime": "3h 25m 12s"
}
```

### Logging

Die Anwendung verwendet strukturiertes JSON-Logging mit konfigurierbaren Log-Leveln:

```json
{
  "level": "info",
  "timestamp": "2024-06-19T14:25:30Z",
  "message": "Server started",
  "component": "http",
  "port": 8080,
  "version": "0.1.0"
}
```

### Backup und Wiederherstellung

Folgende Daten sollten regelmäßig gesichert werden:
- SQLite-Datenbank (`/var/lib/ki-network-analyzer/database.db`)
- Konfigurationsdateien (`/etc/ki-network-analyzer/config.json`)
- Benutzerdefinierte Analysen und Annotationen (`/var/lib/ki-network-analyzer/annotations/`)

## Lokale Entwicklungsumgebung einrichten

### Backend einrichten

```bash
# Repository klonen
git clone https://github.com/username/ki-network-analyzer.git
cd ki-network-analyzer

# Go-Abhängigkeiten installieren
go mod download

# Konfiguration für lokale Entwicklung vorbereiten
cp configs/config.example.json configs/config.local.json
# Bearbeite configs/config.local.json und passe die Einstellungen an

# Backend starten
go run cmd/server/main.go --config configs/config.local.json
```

### Frontend einrichten

```bash
# Frontend-Verzeichnis wechseln
cd web

# Abhängigkeiten installieren
npm install

# Entwicklungsserver starten (Hot-Reloading)
npm run dev
```

## Fazit

Der CI/CD-Prozess des KI-Netzwerk-Analyzers ist auf schnelle Iteration und hohe Qualitätssicherung ausgelegt. Durch die Kombination aus automatisierten Tests, reproduzierbaren Builds und einer klaren Deployment-Strategie wird ein effizienter und zuverlässiger Entwicklungsprozess gewährleistet. 