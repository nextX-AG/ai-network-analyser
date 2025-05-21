#!/bin/bash

# Build-Skript für KI-Netzwerk-Analyzer

set -e

# Basisverzeichnis
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

# Build-Verzeichnis erstellen und bereinigen
BUILD_DIR="$PROJECT_DIR/bin"
mkdir -p "$BUILD_DIR"

# Version-Informationen
VERSION=$(cat VERSION.txt)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "KI-Netzwerk-Analyzer Build"
echo "Version: $VERSION"
echo "Commit: $GIT_COMMIT"
echo "Build-Datum: $BUILD_DATE"
echo

# Abhängigkeiten installieren
echo "Installiere Abhängigkeiten..."
go mod tidy

# Build für aktuelle Plattform
echo "Erstelle Build für aktuelle Plattform..."
go build -o "$BUILD_DIR/analyzer" \
  -ldflags "-X github.com/sayedamirkarim/ki-network-analyzer/pkg/version.Version=$VERSION \
  -X github.com/sayedamirkarim/ki-network-analyzer/pkg/version.CommitHash=$GIT_COMMIT \
  -X github.com/sayedamirkarim/ki-network-analyzer/pkg/version.BuildDate=$BUILD_DATE" \
  cmd/server/main.go

echo "Build erfolgreich erstellt: $BUILD_DIR/analyzer"

# Web-Dateien kopieren
echo "Kopiere Web-Dateien..."
mkdir -p "$BUILD_DIR/web"
cp -r web/* "$BUILD_DIR/web/"

# Konfigurationsdatei kopieren
echo "Kopiere Konfigurationsdatei..."
mkdir -p "$BUILD_DIR/configs"
cp configs/config.example.json "$BUILD_DIR/configs/"

echo "Build abgeschlossen."
echo
echo "Verwende: $BUILD_DIR/analyzer --config=configs/config.example.json" 