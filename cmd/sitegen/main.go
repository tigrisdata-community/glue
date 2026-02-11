package main

import (
	"flag"
	"fmt"
	"log"
)

//go:generate go tool templ generate

var (
	configPath = flag.String("config", "sitegen.yaml", "Path to sitegen.yaml configuration file")
	quiet      = flag.Bool("quiet", false, "Suppress progress output")
)

func main() {
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", *configPath, err)
	}

	if !*quiet {
		fmt.Printf("Generating site from %s to %s\n", cfg.ContentDir, cfg.OutputDir)
	}

	if err := Generate(cfg, *quiet); err != nil {
		log.Fatalf("Site generation failed: %v", err)
	}

	if !*quiet {
		fmt.Println("Site generated successfully")
	}
}
