package main

import (
	"FearNot/internal/orchestrator"
	"log"
)

func main() {
	err := orchestrator.NewOrchestrator()
	if err != nil {
		log.Fatal(err)
	}
}
