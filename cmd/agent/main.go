package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/sayedamirkarim/ki-network-analyzer/internal/agent"
	"github.com/sayedamirkarim/ki-network-analyzer/internal/config"
)

var (
	configFile = flag.String("config", "", "Path to configuration file")
	listenAddr = flag.String("listen", "0.0.0.0:8090", "Address and port to listen on")
	serverAddr = flag.String("server", "http://localhost:9090", "Address of the main server")
	debug      = flag.Bool("debug", false, "Enable debug mode")
	interface_ = flag.String("interface", "", "Network interface to capture packets from")
	name       = flag.String("name", "", "Agent name (defaults to hostname)")
)

func main() {
	flag.Parse()

	// Set up agent name
	agentName := *name
	if agentName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("Warning: Could not determine hostname: %v", err)
			agentName = "unknown-agent"
		} else {
			agentName = hostname
		}
	}

	// Load configuration if specified
	var cfg *config.Config
	var err error
	if *configFile != "" {
		cfg, err = config.LoadConfig(*configFile)
		if err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}
	} else {
		// Use default configuration
		cfg = config.DefaultConfig()
	}

	// Override config with command line flags
	if *listenAddr != "" {
		cfg.Agent = &config.AgentConfig{
			Listen:    *listenAddr,
			ServerURL: *serverAddr,
			Interface: *interface_,
			Name:      agentName,
		}
	} else if cfg.Agent == nil {
		cfg.Agent = &config.AgentConfig{
			Listen:    "0.0.0.0:8090",
			ServerURL: "http://localhost:9090",
			Interface: *interface_,
			Name:      agentName,
		}
	}

	// Create context for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handler
	setupSignalHandler(cancel)

	// Create and start the agent
	captureAgent := agent.NewCaptureAgent(cfg)
	if err := captureAgent.Init(); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// Register with the main server
	if err := captureAgent.Register(); err != nil {
		log.Printf("Warning: Failed to register with main server: %v", err)
	}

	// Set up the HTTP router
	router := mux.NewRouter()

	// Register API routes
	captureAgent.RegisterRoutes(router)

	// Register Admin UI routes
	captureAgent.RegisterAdminHandlers(router)

	// Set up and start the HTTP server
	server := &http.Server{
		Addr:    cfg.Agent.Listen,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting agent server on %s", cfg.Agent.Listen)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation (signal)
	<-ctx.Done()

	// Graceful shutdown
	log.Println("Shutting down agent...")

	// Unregister from the main server
	if err := captureAgent.Unregister(); err != nil {
		log.Printf("Warning: Failed to unregister from main server: %v", err)
	}

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Agent shutdown complete")
}

// setupSignalHandler sets up signal handling for graceful shutdown
func setupSignalHandler(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received signal: %v", sig)
		cancel()
	}()
}
