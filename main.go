package main

import (
	"flag"
	"log"

	"github.com/kroksys/cookie-sync-service/service"
)

var configFile = flag.String("c", "config.yaml", "config file")

func main() {
	// Read configuration file
	s, err := service.ReadConfig(*configFile)
	if err != nil {
		log.Panicf("Error reading config file %s. Err: %v", *configFile, err)
	}

	// Connect to database
	err = service.Connect(s.DSN, s.TablePrefix)
	if err != nil {
		log.Printf("connectionString: %s\n", s.DSN)
		log.Fatalf("Error connecting to database: %s\n", err.Error())
	}

	// Migrate user model to database
	err = service.Migrate()
	if err != nil {
		log.Fatalf("Error migrating user to database: %s\n", err.Error())
	}

	// Start a restful server
	s.Start()
}
