<!-- Version: 0.1.0 | Last Updated: 2024-06-19 14:30:00 UTC -->


# Coding-Konventionen für KI-Netzwerk-Analyzer

Dieses Dokument definiert die einheitlichen Coding-Konventionen für das KI-Netzwerk-Analyzer-Projekt. Die Einhaltung dieser Standards gewährleistet eine konsistente, lesbare und wartbare Codebasis für alle Beteiligten.

## Namenskonventionen

### Dateinamen

#### Go (Backend)
- Kleinbuchstaben mit Unterstrichen zwischen Wörtern (snake_case)
- Aussagekräftige Namen, die den Inhalt klar beschreiben
- Spezifische Präfixe für zusammengehörige Dateien
- Test-Dateien mit Suffix `_test.go`

#### TypeScript/React (Frontend)
- CamelCase mit führendem Großbuchstaben für Komponenten (`UserProfile.tsx`)
- Lowercase für Hooks und Hilfsfunktionen (`useNetworkData.ts`)
- `.tsx` für Dateien mit JSX/TSX-Inhalt
- `.ts` für reine TypeScript-Dateien ohne UI-Komponenten

#### Beispiele für gute Dateinamen

| ❌ Schlecht          | ✅ Gut                    | Erklärung                               |
|---------------------|--------------------------|----------------------------------------|
| `api.go`            | `packet_api.go`          | Spezifiziert, welche API implementiert ist |
| `utils.go`          | `packet_utils.go`        | Gibt den Bereich der Hilfsfunktionen an |
| `client.go`         | `ai_client.go`           | Verdeutlicht, welcher Client implementiert ist |
| `db.go`             | `event_db.go`            | Spezifiziert den Datenbankbereich       |
| `App.tsx`           | `NetworkAnalyzer.tsx`    | Spezifischer Komponentenname statt generisch |

### Typnamen und Strukturen

#### Go
- CamelCase mit führendem Großbuchstaben für exportierte Typen
- CamelCase mit führendem Kleinbuchstaben für interne Typen
- Interface-Namen enden oft mit `-er` (z.B. `PacketCapturer`)

#### TypeScript
- CamelCase mit führendem Großbuchstaben für Typen und Interfaces
- `type` für einfache Typen, `interface` für Verträge mit Methoden
- Präfix `I` für Interfaces vermeiden (veraltet)

#### Beispiele für gute Typnamen

| ❌ Schlecht          | ✅ Gut                    | Erklärung                               |
|---------------------|--------------------------|----------------------------------------|
| `ai_model`          | `AIModel`                | Konsistentes CamelCase                 |
| `PacketHandler`     | `NetworkPacketHandler`   | Spezifischer und selbsterklärender     |
| `IUserType`         | `UserType`               | Kein unnötiges Interface-Präfix        |
| `DBManager`         | `EventDatabaseManager`   | Verdeutlicht den Verantwortungsbereich |
| `HTTPUtil`          | `HttpUtil`               | Konsistentes CamelCase für Akronyme    |

### Funktionsnamen

#### Go
- CamelCase mit führendem Großbuchstaben für exportierte Funktionen
- CamelCase mit führendem Kleinbuchstaben für interne Funktionen
- Verben zur Beschreibung von Aktionen

#### TypeScript/React
- camelCase mit führendem Kleinbuchstaben
- Hooks beginnen immer mit `use` (z.B. `usePacketData`)
- Handler beginnen oft mit `handle` (z.B. `handlePacketSelection`)

#### Beispiele für gute Funktionsnamen

| ❌ Schlecht          | ✅ Gut                    | Erklärung                               |
|---------------------|--------------------------|----------------------------------------|
| `GetData`           | `FetchPacketData`        | Spezifischer und aussagekräftiger      |
| `process`           | `processNetworkPacket`   | Verdeutlicht, was verarbeitet wird     |
| `ai`                | `analyzeWithAI`          | Klare Beschreibung der Aktion          |
| `dbsave`            | `saveEventToDatabase`    | Vollständiger, selbsterklärender Name  |
| `useThing`          | `useNetworkTimeline`     | Spezifischer Hook-Name                 |

### Konstanten und Variablen

#### Go
- Konstanten: CamelCase (exportiert) oder camelCase (intern)
- `const` für Go-Konstanten
- Deskriptive Namen, die Zweck und Inhalt verdeutlichen

#### TypeScript
- Konstanten: SCREAMING_SNAKE_CASE für echte Konstanten
- Reguläre Variablen: camelCase
- React-States: camelCase (`useState`)

#### Beispiele für gute Konstanten und Variablen

| ❌ Schlecht          | ✅ Gut                    | Erklärung                               |
|---------------------|--------------------------|----------------------------------------|
| `t`                 | `packetThreshold`        | Selbsterklärender Name                 |
| `aikey`             | `openAIApiKey`           | Spezifischer, selbsterklärender Name   |
| `TIMEOUT`           | `API_TIMEOUT_MS`         | Mit Einheiten für klarere Bedeutung    |
| `arr`               | `packetsList`            | Beschreibt den Inhalt der Datenstruktur|
| `dbConn`            | `eventDatabaseConnection`| Vollständiger, klarer Name             |

## Paketstruktur und Organisation

### Go-Backend-Struktur

Die Go-Anwendung verwendet das Standard-Go-Projektlayout:

```
/cmd                      # Hauptanwendungen
  /server                 # Backend-Server für die API
  /tools                  # CLI-Tools für Verwaltung/Tests
/internal                 # Private Anwendungs- und Bibliothekscode
  /api                    # API-Handler
  /packet                 # Paketerfassung und -verarbeitung
  /storage                # Datenbankintegration
  /ai                     # KI-Integration
  /speech                 # Speech2Text-Integration
  /timeline               # Event- und Timeline-Logik
  /config                 # Konfigurationsmanagement
/pkg                      # Öffentliche Bibliotheken
  /models                 # Gemeinsame Datenmodelle
  /protocol               # Protokoll-Definitionen
  /utils                  # Gemeinsame Hilfsfunktionen
/web                      # Frontend-Anwendung
/docs                     # Dokumentation
/scripts                  # Build- und Deployment-Skripte
```

### React-Frontend-Struktur

```
/src
  /components             # React-Komponenten
    /common               # Wiederverwendbare UI-Komponenten
    /timeline             # Timeline-spezifische Komponenten
    /network              # Netzwerkvisualisierung
    /ai                   # KI-Ergebnis-Komponenten
    /speech               # Sprachannotierende Komponenten
  /hooks                  # Benutzerdefinierte React-Hooks
  /services               # API-Clients und Services
  /utils                  # Hilfsfunktionen
  /types                  # TypeScript-Typdefinitionen
  /context                # React-Kontexte
  /pages                  # Seitenkomponenten
  /assets                 # Statische Assets
```

### Import-Struktur

#### Go
Importe sollten in drei Gruppen organisiert werden:
```go
import (
    // Standardbibliotheken
    "context"
    "time"
    
    // Externe Pakete
    "github.com/google/gopacket"
    "github.com/gorilla/websocket"
    
    // Interne Pakete
    "github.com/yourusername/ki-network-analyzer/internal/packet"
    "github.com/yourusername/ki-network-analyzer/pkg/models"
)
```

#### TypeScript
```typescript
// Externe Bibliotheken
import React, { useState, useEffect } from 'react';
import { useThree, Canvas } from '@react-three/fiber';

// Interne Importe
import { NetworkPacket } from '../types/NetworkTypes';
import { usePacketData } from '../hooks/usePacketData';

// Komponenten
import Timeline from '../components/Timeline';
import PacketDetails from '../components/PacketDetails';
```

## Kommentare und Dokumentation

### Go-Dokumentation

Jede exportierte Funktion, Konstante oder Struktur sollte mit Godoc-konformen Kommentaren versehen sein:

```go
// PacketAnalyzer verarbeitet Netzwerkpakete und extrahiert relevante Informationen.
// Es unterstützt verschiedene Protokolle und kann mit großen PCAP-Dateien umgehen.
type PacketAnalyzer struct {
    // ...
}

// Analyze analysiert ein einzelnes Paket und gibt ein strukturiertes Ergebnis zurück.
// Es extrahiert Protokollinformationen und relevante Metadaten.
// Gibt einen Fehler zurück, wenn das Paket nicht analysiert werden kann.
func (a *PacketAnalyzer) Analyze(packet gopacket.Packet) (*PacketAnalysisResult, error) {
    // ...
}
```

### TypeScript/React-Dokumentation

TypeScript-Komponenten und -Funktionen sollten mit JSDoc-konformen Kommentaren versehen sein:

```typescript
/**
 * TimelineComponent zeigt Netzwerkereignisse in einer interaktiven Timeline an.
 * Unterstützt Zoom, Pan und Ereignisauswahl.
 * 
 * @param {TimelineProps} props - Komponenten-Properties
 * @returns {JSX.Element} Timeline-Komponente
 */
export const TimelineComponent: React.FC<TimelineProps> = ({ events, onSelect }) => {
    // ...
}
```

### Inline-Kommentare

Komplexe Logik, nicht offensichtliche Entscheidungen und Workarounds sollten mit Inline-Kommentaren erklärt werden:

```go
// Wir verwenden hier einen Puffer von 1MB, da größere Pakete
// zu Performance-Problemen führen können und in typischen
// Netzwerken selten vorkommen
const maxBufferSize = 1024 * 1024
```

## Fehlerbehandlung

### Go-Fehlerbehandlung

- Fehler immer zurückgeben, nicht ignorieren
- Fehler mit Kontext anreichern (fmt.Errorf mit %w)
- Strukturierte Fehlertypen für spezifische Fehlerklassen

```go
if err != nil {
    return fmt.Errorf("failed to parse packet %d: %w", packetID, err)
}
```

Für spezifische Fehlerklassen:

```go
type PacketParsingError struct {
    PacketID  uint64
    Cause     error
    RawData   []byte
}

func (e *PacketParsingError) Error() string {
    return fmt.Sprintf("failed to parse packet %d: %v", e.PacketID, e.Cause)
}

func (e *PacketParsingError) Unwrap() error {
    return e.Cause
}
```

### TypeScript/React-Fehlerbehandlung

- Try-Catch für asynchrone Operationen
- Error Boundaries für React-Komponentenfehler
- Strukturierte Fehlermodelle für API-Antworten

```typescript
try {
  const response = await api.getPackets(filter);
  setPackets(response.data);
} catch (error) {
  setError(`Failed to load packets: ${error instanceof Error ? error.message : String(error)}`);
  console.error("Packet loading error:", error);
}
```

## Best Practices

### Go-Best-Practices

1. **Interfaces dort definieren, wo sie verwendet werden**
2. **Kleine, fokussierte Pakete erstellen**
3. **Dependency Injection für bessere Testbarkeit**
4. **Kontext für Abbruch und Trace-Propagation verwenden**
5. **Unveränderliche Datenstrukturen bevorzugen**

```go
// Gutes Beispiel für Interface-Definition am Verwendungsort
type PacketProcessor struct {
    // Definiere das Interface direkt dort, wo es benötigt wird
    capturer PacketCapturer
}

// PacketCapturer definiert die Schnittstelle für Paketerfassungsdienste
type PacketCapturer interface {
    Capture(ctx context.Context) (<-chan gopacket.Packet, error)
    Close() error
}
```

### React/TypeScript-Best-Practices

1. **Funktionale Komponenten mit Hooks verwenden**
2. **Zustand nach unten durchreichen, Ereignisse nach oben**
3. **Vermeiden von `any` und Typen explizit definieren**
4. **Memoization für teure Berechnungen und Renderings**
5. **React.memo, useMemo und useCallback für Performance**

```typescript
// Beispiel für eine gut strukturierte Komponente
import React, { useState, useCallback, useMemo } from 'react';
import { Packet } from '../types/NetworkTypes';

interface PacketListProps {
  packets: Packet[];
  onSelect: (packet: Packet) => void;
}

export const PacketList: React.FC<PacketListProps> = ({ packets, onSelect }) => {
  const [filter, setFilter] = useState('');
  
  const handleFilterChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setFilter(e.target.value);
  }, []);
  
  const filteredPackets = useMemo(() => {
    return packets.filter(packet => 
      packet.protocol.toLowerCase().includes(filter.toLowerCase())
    );
  }, [packets, filter]);
  
  return (
    <div>
      <input 
        type="text" 
        value={filter} 
        onChange={handleFilterChange} 
        placeholder="Filter by protocol..." 
      />
      <ul>
        {filteredPackets.map(packet => (
          <li key={packet.id} onClick={() => onSelect(packet)}>
            {packet.timestamp} - {packet.protocol}: {packet.summary}
          </li>
        ))}
      </ul>
    </div>
  );
};
```

## Testing

### Go-Tests

- Testdateien mit Suffix `_test.go`
- Table-driven Tests für mehrere Testfälle
- Testify für Assertions verwenden (optional)

```go
func TestPacketParser_Parse(t *testing.T) {
    testCases := []struct {
        name    string
        input   []byte
        want    *ParsedPacket
        wantErr bool
    }{
        {
            name:  "valid tcp packet",
            input: []byte{...}, // TCP-Paket-Bytes
            want:  &ParsedPacket{Protocol: "TCP", ...},
            wantErr: false,
        },
        {
            name:  "malformed packet",
            input: []byte{...}, // Fehlerhafte Bytes
            want:  nil,
            wantErr: true,
        },
    }
    
    parser := NewPacketParser()
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := parser.Parse(tc.input)
            
            if tc.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tc.want.Protocol, got.Protocol)
            // Weitere Assertions...
        })
    }
}
```

### React/TypeScript-Tests

- Jest + React Testing Library für Komponententests
- Cypress für End-to-End-Tests
- Mock-Services für API-Aufrufe

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { PacketList } from './PacketList';

describe('PacketList', () => {
  const mockPackets = [
    { id: '1', protocol: 'TCP', timestamp: '2024-01-01 12:00:00', summary: 'Connection established' },
    { id: '2', protocol: 'HTTP', timestamp: '2024-01-01 12:00:01', summary: 'GET /api/data' }
  ];
  
  const mockOnSelect = jest.fn();
  
  it('renders all packets', () => {
    render(<PacketList packets={mockPackets} onSelect={mockOnSelect} />);
    
    expect(screen.getByText(/TCP: Connection established/)).toBeInTheDocument();
    expect(screen.getByText(/HTTP: GET \/api\/data/)).toBeInTheDocument();
  });
  
  it('filters packets based on input', () => {
    render(<PacketList packets={mockPackets} onSelect={mockOnSelect} />);
    
    fireEvent.change(screen.getByPlaceholderText('Filter by protocol...'), {
      target: { value: 'http' }
    });
    
    expect(screen.queryByText(/TCP: Connection established/)).not.toBeInTheDocument();
    expect(screen.getByText(/HTTP: GET \/api\/data/)).toBeInTheDocument();
  });
  
  it('calls onSelect when a packet is clicked', () => {
    render(<PacketList packets={mockPackets} onSelect={mockOnSelect} />);
    
    fireEvent.click(screen.getByText(/TCP: Connection established/));
    
    expect(mockOnSelect).toHaveBeenCalledWith(mockPackets[0]);
  });
});
```

## Performance-Richtlinien

### Go-Performance

1. **Minimale Allokationen in Hot-Paths**
2. **Pufferwiederverwendung bei wiederholten Operationen**
3. **Parallelisierung mit Worker-Pools für CPU-intensive Aufgaben**
4. **Effiziente JSON-Verarbeitung (json.Decoder statt json.Unmarshal für Streams)**
5. **Datenbank-Optimierung (Indizes, Prepared Statements, Pooling)**

### React-Performance

1. **Memo für teure Renderingvorgänge**
2. **Virtualisierung für lange Listen**
3. **Code-Splitting und Lazy-Loading**
4. **Optimierte Bundles (Tree Shaking)**
5. **Service-Worker für Caching**

## Abhängigkeitsmanagement

### Go-Module

- `go.mod` und `go.sum` im Repository
- Explizite Versionierung von Abhängigkeiten
- Regelmäßige Aktualisierung von Abhängigkeiten

### Frontend-Pakete

- `package.json` mit festen Versionen oder engen Versionsbereichen
- Yarn oder npm mit Lockfile im Repository
- Regelmäßiges Dependency-Auditing

## Gängige Codierungsmuster

### Dependency Injection (Go)

```go
// Factory-Funktion mit Abhängigkeiten
func NewPacketAnalyzer(
    capturer packet.Capturer,
    storage storage.EventStorage,
    aiClient ai.Client,
) *PacketAnalyzer {
    return &PacketAnalyzer{
        capturer: capturer,
        storage:  storage,
        aiClient: aiClient,
    }
}
```

### Context-API (React)

```typescript
import React, { createContext, useContext, useState } from 'react';

interface NetworkContextType {
  selectedPacket: Packet | null;
  selectPacket: (packet: Packet) => void;
}

const NetworkContext = createContext<NetworkContextType | undefined>(undefined);

export const NetworkProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  const [selectedPacket, setSelectedPacket] = useState<Packet | null>(null);
  
  const selectPacket = (packet: Packet) => {
    setSelectedPacket(packet);
  };
  
  return (
    <NetworkContext.Provider value={{ selectedPacket, selectPacket }}>
      {children}
    </NetworkContext.Provider>
  );
};

export const useNetworkContext = (): NetworkContextType => {
  const context = useContext(NetworkContext);
  if (context === undefined) {
    throw new Error('useNetworkContext must be used within a NetworkProvider');
  }
  return context;
};
```

## Weiterführende Links

### Go-Ressourcen

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### React/TypeScript-Ressourcen

- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [Clean Code in TypeScript](https://github.com/labs42io/clean-code-typescript)
- [React Performance Optimization](https://reactjs.org/docs/optimizing-performance.html) 