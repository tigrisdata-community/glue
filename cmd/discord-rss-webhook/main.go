package main

import (
	"flag"
	"log"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	flagenv.Parse()
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal("error:", err)
	}
}

func run() error {
	return nil
}
