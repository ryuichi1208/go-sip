package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/go-sip/config"
	"github.com/user/go-sip/sip"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("config", "config.json", "Path to configuration file")
	generateConfig := flag.Bool("generate-config", false, "Generate default config file and exit")
	overridePort := flag.String("port", "", "Override port setting from config file")
	overrideBindAddr := flag.String("bind", "", "Override bind address setting from config file")
	flag.Parse()

	// Generate default configuration file option
	if *generateConfig {
		cfg := config.DefaultConfig()
		if err := config.SaveConfig(cfg, *configPath); err != nil {
			log.Fatalf("Error generating config file: %v", err)
		}
		fmt.Printf("Default configuration file generated: %s\n", *configPath)
		return
	}

	// Load configuration file
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Warning: %v", err)
		log.Printf("Using default configuration")
		cfg = config.DefaultConfig()
	}

	// Override settings with command line options
	if *overridePort != "" {
		cfg.Server.Port = *overridePort
	}

	if *overrideBindAddr != "" {
		cfg.Server.BindAddr = *overrideBindAddr
	}

	// Create SIP server
	server := sip.NewServer(cfg.Server.Port)
	server.SetBindAddr(cfg.Server.BindAddr)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server startup error: %v", err)
		}
	}()

	fmt.Printf("SIP server started on %s:%s\n", cfg.Server.BindAddr, cfg.Server.Port)
	fmt.Println("Press Ctrl+C to exit...")

	// Wait for signal
	<-sigChan
	fmt.Println("\nShutting down server...")
}
