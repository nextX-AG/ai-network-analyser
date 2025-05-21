<!-- Version: 0.1.0 | Last Updated: 2024-06-19 14:30:00 UTC -->


# Versionierungssystem: KI-Netzwerk-Analyzer

## Übersicht

Dieses Dokument beschreibt das Versionierungssystem für den KI-Netzwerk-Analyzer. Es definiert die Versionierungsstrategie sowohl für die Codebasis als auch für die Dokumentation und stellt sicher, dass alle Änderungen nachvollziehbar und konsistent sind.

## Versionierungskonzept

### Semantische Versionierung

Wir folgen dem Prinzip der [Semantischen Versionierung](https://semver.org/lang/de/) im Format:

- **Major.Minor.Patch** (z.B. 1.2.3)
  - **Major**: Inkompatible API-Änderungen, grundlegende Umstrukturierungen
  - **Minor**: Abwärtskompatible Funktionserweiterungen
  - **Patch**: Bugfixes und kleinere Verbesserungen

Die Versionen werden nach folgenden Regeln erhöht:
- **Major**: Bei API-Breaking-Changes oder signifikanten Architekturänderungen
- **Minor**: Bei neuen Features, die die bestehende API nicht ändern
- **Patch**: Bei Bugfixes, Optimierungen oder kleineren Dokumentationsupdates

### Entwicklungsversionen

Während der frühen Entwicklungsphase (vor 1.0.0) können Änderungen wie folgt behandelt werden:
- **0.x.y**: Die Minor-Version (x) kann Breaking Changes enthalten
- **0.0.z**: Experimentelle Versionen in der frühen Entwicklungsphase

## Code-Versionierung

### Zentrale Versionsverwaltung

Die aktuelle Version wird zentral in zwei Dateien verwaltet:
- `VERSION.txt`: Enthält nur die Versionsnummer (z.B. "0.1.0")
- `pkg/version/version.go`: Go-Konstanten für programmatischen Zugriff

Beispiel für `pkg/version/version.go`:
```go
package version

// Version-Informationen
const (
    // Semantische Versionskomponenten
    Major      = 0
    Minor      = 1
    Patch      = 0
    
    // Vollständige Version als String
    Version    = "0.1.0"
    
    // Build-Informationen (werden vom Build-System gesetzt)
    BuildDate  = ""
    CommitHash = ""
)
```

### Build-Metadaten

Jeder Build erhält automatisch Metadaten:
- Git-Commit-Hash
- Build-Zeitstempel
- Bei Releases: Tag-Informationen

Diese werden zur Laufzeit über Flags an den Go-Compiler übergeben:
```bash
go build -ldflags "-X github.com/username/ki-network-analyzer/pkg/version.BuildDate=$(date -u '+%Y-%m-%d %H:%M:%S') -X github.com/username/ki-network-analyzer/pkg/version.CommitHash=$(git rev-parse HEAD)"
```

## Dokumentations-Versionierung

### Automatische Versionierung

Alle Markdown-Dokumente im `docs/`-Verzeichnis werden automatisch versioniert:

1. Jedes Dokument enthält einen Header mit Versionsinformationen:
   ```markdown
   <!-- Version: 0.1.0 | Last Updated: 2024-06-19 14:30:00 UTC -->
   ```

2. Ein Git Pre-Commit Hook aktualisiert diese Header automatisch:
   - Erkennt geänderte Markdown-Dateien
   - Aktualisiert die Zeitstempel
   - Synchronisiert mit der aktuellen Projektversion

### Git Pre-Commit Hook

Der folgende Pre-Commit Hook wird für die automatische Dokumentationsversionierung verwendet:

```bash
#!/bin/bash

# Finde geänderte Markdown-Dateien im docs/-Verzeichnis
docs_changed=$(git diff --cached --name-only --diff-filter=ACM | grep '^docs/.*\.md$')

if [ -n "$docs_changed" ]; then
  # Aktuelle Version aus VERSION.txt lesen
  if [ -f VERSION.txt ]; then
    version=$(cat VERSION.txt)
  else
    version="0.1.0"
  fi
  
  # Zeitstempel im UTC-Format generieren
  timestamp=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
  
  # Header in allen geänderten Dokumenten aktualisieren
  for file in $docs_changed; do
    # Entferne vorhandenen Versions-Header (falls vorhanden)
    sed -i '/^<!-- Version: .* | Last Updated: .* -->/d' $file
    # Füge neuen Versions-Header hinzu
    sed -i "1s/^/<!-- Version: $version | Last Updated: $timestamp -->\n\n/" $file
    git add $file
  done
fi
```

> **Hinweis:** Für macOS benötigt der `sed`-Befehl ein leeres Argument nach dem `-i`-Flag:
> ```bash
> sed -i '' '/^<!-- Version: .* | Last Updated: .* -->/d' $file
> sed -i '' "1s/^/<!-- Version: $version | Last Updated: $timestamp -->\n\n/" $file
> ```

## Release-Prozess

### Release-Workflow

1. **Vorbereitung**:
   - Aktualisierung der VERSION.txt mit der neuen Versionsnummer
   - Aktualisierung von pkg/version/version.go
   - Aktualisierung des CHANGELOG.md

2. **Release-Commit**:
   ```bash
   # Aktualisiere Versionsdateien
   echo "0.2.0" > VERSION.txt
   # Aktualisiere version.go manuell oder per Skript
   
   # Commit und Tag
   git add VERSION.txt pkg/version/version.go CHANGELOG.md
   git commit -m "Release v0.2.0"
   git tag -a v0.2.0 -m "Version 0.2.0"
   git push origin main --tags
   ```

3. **CI/CD**:
   - Der Tag löst automatisch den Build-Prozess aus
   - Artefakte werden erstellt und signiert
   - Release-Notes werden generiert

### Changelog-Verwaltung

Das CHANGELOG.md wird für jede Version aktualisiert und folgt diesem Format:

```markdown
# Changelog

## [0.2.0] - 2024-06-20

### Hinzugefügt
- Feature A: Beschreibung des Features
- Feature B: Beschreibung des Features

### Geändert
- Komponente X: Beschreibung der Änderung
- Komponente Y: Beschreibung der Änderung

### Behoben
- Bug 1: Beschreibung des Bugfixes
- Bug 2: Beschreibung des Bugfixes

## [0.1.0] - 2024-06-01

### Hinzugefügt
- Initiale Version mit Grundfunktionalität
```

## Automatisierung mit GitHub Actions

### Changelog-Generator

Ein GitHub Actions Workflow generiert automatisch Changelogs basierend auf Pull Requests und Commit-Nachrichten:

```yaml
# .github/workflows/changelog.yml
name: Generate Changelog

on:
  push:
    tags:
      - 'v*'

jobs:
  changelog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Generate Changelog
        id: changelog
        uses: metcalfc/changelog-generator@v4.0.1
        with:
          myToken: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Update CHANGELOG.md
        run: |
          VERSION=${GITHUB_REF#refs/tags/v}
          DATE=$(date +"%Y-%m-%d")
          
          # Erzeuge neuen Changelog-Eintrag
          echo -e "## [$VERSION] - $DATE\n\n${{ steps.changelog.outputs.changelog }}\n\n$(cat CHANGELOG.md)" > CHANGELOG.md
          
          # Committe Änderungen zum Branch
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add CHANGELOG.md
          git commit -m "Update CHANGELOG.md for $VERSION"
          git push
```

### Versions-Synchronisierung

Eine weitere GitHub Action stellt sicher, dass alle versionsbezogenen Dateien synchron bleiben:

```yaml
# .github/workflows/version-sync.yml
name: Sync Version Files

on:
  push:
    paths:
      - 'VERSION.txt'

jobs:
  sync-versions:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Update Version in Go File
        run: |
          VERSION=$(cat VERSION.txt)
          MAJOR=$(echo $VERSION | cut -d. -f1)
          MINOR=$(echo $VERSION | cut -d. -f2)
          PATCH=$(echo $VERSION | cut -d. -f3)
          
          # Aktualisiere version.go
          cat > pkg/version/version.go << EOF
          package version

          // Version-Informationen
          const (
              // Semantische Versionskomponenten
              Major      = $MAJOR
              Minor      = $MINOR
              Patch      = $PATCH
              
              // Vollständige Version als String
              Version    = "$VERSION"
              
              // Build-Informationen (werden vom Build-System gesetzt)
              BuildDate  = ""
              CommitHash = ""
          )
          EOF
          
          # Commit und Push
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add pkg/version/version.go
          git commit -m "Sync version.go with VERSION.txt ($VERSION)"
          git push
```

## Versionierungsrichtlinien für Entwickler

### Commit-Nachrichten

Wir verwenden [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Typen:
- `feat`: Neue Funktionen
- `fix`: Bugfixes
- `docs`: Dokumentationsänderungen
- `style`: Formatierungsänderungen
- `refactor`: Code-Refactoring ohne funktionale Änderungen
- `perf`: Performance-Verbesserungen
- `test`: Test-bezogene Änderungen
- `build`: Build-System oder externe Abhängigkeiten
- `ci`: CI-Konfigurationsänderungen

Beispiele:
- `feat(timeline): Zoom-Funktionalität für Event-Timeline hinzufügen`
- `fix(api): Race-Condition in der Websocket-Verbindung beheben`
- `docs(setup): Installationsanleitung für macOS aktualisieren`

### Branching-Strategie

Wir folgen dem GitHub-Flow mit einigen Anpassungen:

1. Feature-/Fix-Branches werden von `main` abgezweigt
2. Nach Fertigstellung wird ein Pull Request erstellt
3. Nach Code-Review und Tests erfolgt der Merge in `main`
4. Für Releases werden Tags auf `main` erstellt

Branch-Namenskonventionen:
- `feature/beschreibung-des-features`
- `fix/beschreibung-des-bugfixes`
- `docs/beschreibung-der-dokumentationsänderung`
- `refactor/beschreibung-des-refactorings`

## Installation und Aktualisierung des Versionierungssystems

### Ersteinrichtung

1. **Initialisierung der Versionsdateien**:
   ```bash
   echo "0.1.0" > VERSION.txt
   mkdir -p pkg/version
   # Erstelle version.go manuell oder per Skript
   ```

2. **Einrichtung des Pre-Commit Hooks**:
   ```bash
   mkdir -p .git/hooks
   # Kopiere den oben beschriebenen Pre-Commit Hook
   chmod +x .git/hooks/pre-commit
   ```

3. **Einrichtung der GitHub Actions**:
   ```bash
   mkdir -p .github/workflows
   # Erstelle die oben beschriebenen Workflow-Dateien
   ```

### Manuelles Erhöhen der Version

Für manuelle Versionsaktualisierungen:

```bash
# Erhöhung der Minor-Version
VERSION=$(cat VERSION.txt)
MAJOR=$(echo $VERSION | cut -d. -f1)
MINOR=$(echo $VERSION | cut -d. -f2)
NEW_MINOR=$((MINOR + 1))
echo "$MAJOR.$NEW_MINOR.0" > VERSION.txt

# Änderungen committen
git add VERSION.txt
git commit -m "Bump version to $(cat VERSION.txt)"
```

## Vorteile dieses Versionierungsansatzes

1. **Konsistenz**: Einheitliche Versionierung über Code und Dokumentation
2. **Automatisierung**: Minimaler manueller Aufwand durch Git-Hooks und CI/CD
3. **Nachvollziehbarkeit**: Klare Verbindung zwischen Änderungen und Versionen
4. **Semantische Klarheit**: Eindeutige Bedeutung von Versionsnummern
5. **Einfache Integration**: Versionsinformationen sind sowohl für Menschen als auch programmatisch zugänglich 